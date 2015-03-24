package core

import (
	"github.com/Nebuleuse/Nebuleuse/git"
)

var Repo *git.Repository

func InitGitUpdater(path string) error {
	var err error
	Repo, err = git.OpenRepository(path)
	return err
}

func GitPreparePatch() error {
	//Todo
	return nil
}
