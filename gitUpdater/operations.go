// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//Original file from GOGS (https://github.com/gogits/gogs), Edited and extended for our use
package gitUpdater

import (
	"Nebuleuse/gitUpdater/git"
)

type Repository git.Repository

var _repo *Repository

func InitGit(path string) error {
	_repo, err := git.OpenRepository(path)
	return err
}

func PreparePatch() error {
	//Todo
	return nil
}
