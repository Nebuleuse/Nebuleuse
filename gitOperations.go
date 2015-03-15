// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//Original file from GOGS (https://github.com/gogits/gogs), improved for our uses
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Unknwon/com"
	"github.com/gogits/gogs/modules/git"
	"github.com/gogits/gogs/modules/process"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Diff line types.
const (
	DIFF_LINE_PLAIN = iota + 1
	DIFF_LINE_ADD
	DIFF_LINE_DEL
	DIFF_LINE_SECTION
)
const (
	DIFF_FILE_ADD = iota + 1
	DIFF_FILE_CHANGE
	DIFF_FILE_DEL
)

type DiffFile struct {
	Name               string
	Index              int
	Addition, Deletion int
	Type               int
	IsCreated          bool
	IsDeleted          bool
	IsBin              bool
}
type Diff struct {
	TotalAddition, TotalDeletion int
	Files                        []*DiffFile
}

func (diff *Diff) NumFiles() int {
	return len(diff.Files)
}

const DIFF_HEAD = "diff --git "

func ParsePatch(pid int64, cmd *exec.Cmd, reader io.Reader) (*Diff, error) {
	scanner := bufio.NewScanner(reader)
	var (
		curFile             *DiffFile
		leftLine, rightLine int
		// FIXME: use first 30 lines to detect file encoding. Should use cache in the future.
		buf bytes.Buffer
	)
	diff := &Diff{Files: make([]*DiffFile, 0)}
	var i int
	for scanner.Scan() {
		line := scanner.Text()
		// fmt.Println(i, line)
		if strings.HasPrefix(line, "+++ ") || strings.HasPrefix(line, "--- ") {
			continue
		}
		if line == "" {
			continue
		}
		i = i + 1
		// FIXME: use first 30 lines to detect file encoding.
		if i <= 30 {
			buf.WriteString(line)
		}

		switch {
		case line[0] == ' ':
			leftLine++
			rightLine++
			continue
		case line[0] == '@':
			ss := strings.Split(line, "@@")
			// Parse line number.
			ranges := strings.Split(ss[len(ss)-2][1:], " ")
			leftLine, _ = com.StrTo(strings.Split(ranges[0], ",")[0][1:]).Int()
			rightLine, _ = com.StrTo(strings.Split(ranges[1], ",")[0]).Int()
			continue
		case line[0] == '+':
			curFile.Addition++
			diff.TotalAddition++
			rightLine++
			continue
		case line[0] == '-':
			curFile.Deletion++
			diff.TotalDeletion++
			if leftLine > 0 {
				leftLine++
			}
		case strings.HasPrefix(line, "Binary"):
			curFile.IsBin = true
			continue
		}
		// Get new file.
		if strings.HasPrefix(line, DIFF_HEAD) {
			fs := strings.Split(line[len(DIFF_HEAD):], " ")
			a := fs[0]
			curFile = &DiffFile{
				Name:  a[strings.Index(a, "/")+1:],
				Index: len(diff.Files) + 1,
				Type:  DIFF_FILE_CHANGE,
			}
			diff.Files = append(diff.Files, curFile)
			// Check file diff type.
			for scanner.Scan() {
				switch {
				case strings.HasPrefix(scanner.Text(), "new file"):
					curFile.Type = DIFF_FILE_ADD
					curFile.IsDeleted = false
					curFile.IsCreated = true
				case strings.HasPrefix(scanner.Text(), "deleted"):
					curFile.Type = DIFF_FILE_DEL
					curFile.IsCreated = false
					curFile.IsDeleted = true
				case strings.HasPrefix(scanner.Text(), "index"):
					curFile.Type = DIFF_FILE_CHANGE
					curFile.IsCreated = false
					curFile.IsDeleted = false
				}
				if curFile.Type > 0 {
					break
				}
			}
		}
	}

	return diff, nil
}
func (r *git.Repository) GetDiffRange(beforeCommitId string, afterCommitId string) (*Diff, error) {
	repo, err := git.OpenRepository(r.Path)
	if err != nil {
		return nil, err
	}
	commit, err := repo.GetCommit(afterCommitId)
	if err != nil {
		return nil, err
	}
	rd, wr := io.Pipe()
	var cmd *exec.Cmd
	// if "after" commit given
	if beforeCommitId == "" {
		// First commit of repository.
		if commit.ParentCount() == 0 {
			cmd = exec.Command("git", "show", afterCommitId)
		} else {
			c, _ := commit.Parent(0)
			cmd = exec.Command("git", "diff", c.Id.String(), afterCommitId)
		}
	} else {
		cmd = exec.Command("git", "diff", beforeCommitId, afterCommitId)
	}
	cmd.Dir = r.Path
	cmd.Stdout = wr
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	done := make(chan error)
	go func() {
		cmd.Start()
		done <- cmd.Wait()
		wr.Close()
	}()
	defer rd.Close()
	desc := fmt.Sprintf("GetDiffRange(%s)", r.Path)
	pid := process.Add(desc, cmd)
	go func() {
		// In case process became zombie.
		select {
		case <-time.After(5 * time.Minute):
			if errKill := process.Kill(pid); errKill != nil {
				Warning.Println(4, "git_diff.ParsePatch(Kill): %v", err)
			}
			<-done
		// return "", ErrExecTimeout.Error(), ErrExecTimeout
		case err = <-done:
			process.Remove(pid)
		}
	}()
	return ParsePatch(pid, cmd, rd)
}
func (r *git.Repository) GetDiffCommit(commitId string) (*Diff, error) {
	return r.GetDiffRange("", commitId)
}

func (r *git.Repository) UpdateGitRepo() {
	cmd := exec.Command("git", "pull", "origin", _cfg["productionBranch"])
	cmd.Dir = r.Path
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()
}
