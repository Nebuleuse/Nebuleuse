// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package git

import (
	"os"
	"os/exec"
	"path/filepath"
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

func (r *Repository) UpdateGitRepo(branch string) {
	cmd := exec.Command("git", "pull", "origin", branch)
	cmd.Dir = r.Path
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Start()
	cmd.Wait()
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
