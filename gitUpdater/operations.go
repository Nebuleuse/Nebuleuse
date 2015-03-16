// Copyright 2014 The Gogs Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//Original file from GOGS (https://github.com/gogits/gogs), Edited and extended for our use
package gitUpdater

import (
	"Nebuleuse/gitUpdater/git"
)

var Repo *git.Repository

func Init(path string) error {
	var err error
	Repo, err = git.OpenRepository(path)
	return err
}

func PreparePatch() error {
	//Todo
	return nil
}
