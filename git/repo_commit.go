// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bufio"
	"bytes"
	"container/list"
	"errors"
	"fmt"
	"github.com/gogits/gogs/modules/process"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/Unknwon/com"
)

func (repo *Repository) getCommitIdOfRef(refpath string) (string, error) {
	stdout, stderr, err := com.ExecCmdDir(repo.Path, "git", "show-ref", "--verify", refpath)
	if err != nil {
		return "", errors.New(stderr)
	}
	return strings.Split(stdout, " ")[0], nil
}

func (repo *Repository) GetCommitIdOfBranch(branchName string) (string, error) {
	return repo.getCommitIdOfRef("refs/heads/" + branchName)
}

// get branch's last commit or a special commit by id string
func (repo *Repository) GetCommitOfBranch(branchName string) (*Commit, error) {
	commitId, err := repo.GetCommitIdOfBranch(branchName)
	if err != nil {
		return nil, err
	}
	return repo.GetCommit(commitId)
}

func (repo *Repository) GetCommitIdOfTag(tagName string) (string, error) {
	return repo.getCommitIdOfRef("refs/tags/" + tagName)
}

func (repo *Repository) GetCommitOfTag(tagName string) (*Commit, error) {
	tag, err := repo.GetTag(tagName)
	if err != nil {
		return nil, err
	}
	return tag.Commit()
}

// Parse commit information from the (uncompressed) raw
// data from the commit object.
// \n\n separate headers from message
func parseCommitData(data []byte) (*Commit, error) {
	commit := new(Commit)
	commit.parents = make([]sha1, 0, 1)
	// we now have the contents of the commit object. Let's investigate...
	nextline := 0
l:
	for {
		eol := bytes.IndexByte(data[nextline:], '\n')
		switch {
		case eol > 0:
			line := data[nextline : nextline+eol]
			spacepos := bytes.IndexByte(line, ' ')
			reftype := line[:spacepos]
			switch string(reftype) {
			case "tree":
				id, err := NewIdFromString(string(line[spacepos+1:]))
				if err != nil {
					return nil, err
				}
				commit.Tree.Id = id
			case "parent":
				// A commit can have one or more parents
				oid, err := NewIdFromString(string(line[spacepos+1:]))
				if err != nil {
					return nil, err
				}
				commit.parents = append(commit.parents, oid)
			case "author":
				sig, err := newSignatureFromCommitline(line[spacepos+1:])
				if err != nil {
					return nil, err
				}
				commit.Author = sig
			case "committer":
				sig, err := newSignatureFromCommitline(line[spacepos+1:])
				if err != nil {
					return nil, err
				}
				commit.Committer = sig
			}
			nextline += eol + 1
		case eol == 0:
			commit.CommitMessage = string(data[nextline+1:])
			break l
		default:
			break l
		}
	}
	return commit, nil
}

func (repo *Repository) getCommit(id sha1) (*Commit, error) {
	if repo.commitCache != nil {
		if c, ok := repo.commitCache[id]; ok {
			return c, nil
		}
	} else {
		repo.commitCache = make(map[sha1]*Commit, 10)
	}

	data, bytErr, err := com.ExecCmdDirBytes(repo.Path, "git", "cat-file", "-p", id.String())
	if err != nil {
		return nil, errors.New(err.Error() + ": " + string(bytErr))
	}

	commit, err := parseCommitData(data)
	if err != nil {
		return nil, err
	}
	commit.repo = repo
	commit.Id = id

	repo.commitCache[id] = commit
	return commit, nil
}

// Find the commit object in the repository.
func (repo *Repository) GetCommit(commitId string) (*Commit, error) {
	id, err := NewIdFromString(commitId)
	if err != nil {
		return nil, err
	}

	return repo.getCommit(id)
}

func (repo *Repository) commitsCount(id sha1) (int, error) {
	if gitVer.LessThan(MustParseVersion("1.8.0")) {
		stdout, stderr, err := com.ExecCmdDirBytes(repo.Path, "git", "log", "--pretty=format:''", id.String())
		if err != nil {
			return 0, errors.New(string(stderr))
		}
		return len(bytes.Split(stdout, []byte("\n"))), nil
	}

	stdout, stderr, err := com.ExecCmdDir(repo.Path, "git", "rev-list", "--count", id.String())
	if err != nil {
		return 0, errors.New(stderr)
	}
	return com.StrTo(strings.TrimSpace(stdout)).Int()
}

// used only for single tree, (]
func (repo *Repository) CommitsBetween(last *Commit, before *Commit) (*list.List, error) {
	l := list.New()
	if last == nil || last.ParentCount() == 0 {
		return l, nil
	}

	var err error
	cur := last
	for {
		if cur.Id.Equal(before.Id) {
			break
		}
		l.PushBack(cur)
		if cur.ParentCount() == 0 {
			break
		}
		cur, err = cur.Parent(0)
		if err != nil {
			return nil, err
		}
	}
	return l, nil
}

func (repo *Repository) commitsBefore(lock *sync.Mutex, l *list.List, parent *list.Element, id sha1, limit int) error {
	commit, err := repo.getCommit(id)
	if err != nil {
		return err
	}

	var e *list.Element
	if parent == nil {
		e = l.PushBack(commit)
	} else {
		var in = parent
		for {
			if in == nil {
				break
			} else if in.Value.(*Commit).Id.Equal(commit.Id) {
				return nil
			} else {
				if in.Next() == nil {
					break
				}
				if in.Value.(*Commit).Committer.When.Equal(commit.Committer.When) {
					break
				}

				if in.Value.(*Commit).Committer.When.After(commit.Committer.When) &&
					in.Next().Value.(*Commit).Committer.When.Before(commit.Committer.When) {
					break
				}
			}
			in = in.Next()
		}

		e = l.InsertAfter(commit, in)
	}

	var pr = parent
	if commit.ParentCount() > 1 {
		pr = e
	}

	for i := 0; i < commit.ParentCount(); i++ {
		id, err := commit.ParentId(i)
		if err != nil {
			return err
		}
		err = repo.commitsBefore(lock, l, pr, id, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repo *Repository) CommitsCount(commitId string) (int, error) {
	id, err := NewIdFromString(commitId)
	if err != nil {
		return 0, err
	}
	return repo.commitsCount(id)
}

func (repo *Repository) FileCommitsCount(branch, file string) (int, error) {
	stdout, stderr, err := com.ExecCmdDir(repo.Path, "git", "rev-list", "--count",
		branch, "--", file)
	if err != nil {
		return 0, errors.New(stderr)
	}
	return com.StrTo(strings.TrimSpace(stdout)).Int()
}

func (repo *Repository) CommitsByFileAndRange(branch, file string, page int) (*list.List, error) {
	stdout, stderr, err := com.ExecCmdDirBytes(repo.Path, "git", "log", branch,
		"--skip="+com.ToStr((page-1)*50), "--max-count=50", prettyLogFormat, "--", file)
	if err != nil {
		return nil, errors.New(string(stderr))
	}
	return parsePrettyFormatLog(repo, stdout)
}

func (repo *Repository) getCommitsBefore(id sha1) (*list.List, error) {
	l := list.New()
	lock := new(sync.Mutex)
	err := repo.commitsBefore(lock, l, nil, id, 0)
	return l, err
}

func (repo *Repository) searchCommits(id sha1, keyword string) (*list.List, error) {
	stdout, stderr, err := com.ExecCmdDirBytes(repo.Path, "git", "log", id.String(), "-100",
		"-i", "--grep="+keyword, prettyLogFormat)
	if err != nil {
		return nil, err
	} else if len(stderr) > 0 {
		return nil, errors.New(string(stderr))
	}
	return parsePrettyFormatLog(repo, stdout)
}

func (repo *Repository) commitsByRange(id sha1, page int) (*list.List, error) {
	stdout, stderr, err := com.ExecCmdDirBytes(repo.Path, "git", "log", id.String(),
		"--skip="+com.ToStr((page-1)*50), "--max-count=50", prettyLogFormat)
	if err != nil {
		return nil, errors.New(string(stderr))
	}
	return parsePrettyFormatLog(repo, stdout)
}

func (repo *Repository) getCommitOfRelPath(id sha1, relPath string) (*Commit, error) {
	stdout, _, err := com.ExecCmdDir(repo.Path, "git", "log", "-1", prettyLogFormat, id.String(), "--", relPath)
	if err != nil {
		return nil, err
	}

	id, err = NewIdFromString(string(stdout))
	if err != nil {
		return nil, err
	}

	return repo.getCommit(id)
}

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

func parsePatch(pid int64, cmd *exec.Cmd, reader io.Reader) (*Diff, error) {
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
func (r *Repository) GetDiffRange(beforeCommitId string, afterCommitId string) (*Diff, error) {
	commit, err := r.GetCommit(afterCommitId)
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
				log.Println(4, "git_diff.ParsePatch(Kill): %v", err)
			}
			<-done
		// return "", ErrExecTimeout.Error(), ErrExecTimeout
		case err = <-done:
			process.Remove(pid)
		}
	}()
	return parsePatch(pid, cmd, rd)
}
func (r *Repository) GetDiffCommit(commitId string) (*Diff, error) {
	return r.GetDiffRange("", commitId)
}
