// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"bufio"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Repository represents a Git repository.
type Repository struct {
	Path string

	commitCache map[sha1]*Commit
	tagCache    map[sha1]*Tag
}

// OpenRepository opens the repository at the given path.
func OpenRepository(repoPath string) (*Repository, error) {
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return nil, err
	}

	return &Repository{Path: repoPath}, nil
}
func parseUpdateRepo(rd io.Reader) error {
	scanner := bufio.NewScanner(rd)
	failed := false
	var err string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "fatal:") {
			if !failed {
				failed = true
				err = "Git pull failed: "
			}
			err = err + line + "\n"
		}
	}
	if failed {
		return errors.New(err)
	}
	return nil
}
func (r *Repository) UpdateGitRepo(branch string) error {
	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Dir = r.Path
	rd, wr := io.Pipe()

	//Todo: Add parser for error output from git
	//Bug: setting cmd.Stderr to wr makes cmd.start freeze
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()
	wr.Close()
	defer rd.Close()
	return parseUpdateRepo(rd)
}
func (r *Repository) Checkout(commit string) {
	cmd := exec.Command("git", "checkout", commit)
	cmd.Dir = r.Path
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()
}

func (r *Repository) GetCommitsSinceLastUpdate(branch, latestCommit string) ([]*Commit, error) {
	var ret []*Commit

	actual, err := r.GetCommitOfBranch(branch)
	last, err := r.GetCommit(latestCommit)
	list, err := r.CommitsBetween(actual, last)

	for e := list.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Commit))
	}

	return ret, err
}

func (r *Repository) GetRecentCommits(branch string) ([]*Commit, error) {
	var ret []*Commit

	latest, err := r.GetCommitOfBranch(branch)
	list, err := latest.CommitsByRange(1)

	for e := list.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(*Commit))
	}

	return ret, err
}

func (r *Repository) GetFilesChangedSinceUpdate(branch, latestCommit string) (*Diff, error) {
	latest, err := r.GetCommitOfBranch(branch)

	res, err := r.GetDiffRange(latestCommit, latest.Id.String())

	return res, err
}
func (r *Repository) GetFilesChangedSinceUpdateRange(start, end string) (*Diff, error) {
	res, err := r.GetDiffRange(start, end)

	return res, err
}
