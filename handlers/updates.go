package handlers

import (
	"net/http"
	"strconv"

	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
)

//User connected
func getBranchList(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*core.UserSession)
	list := core.GetBranchList(session.UserRank)
	EasyDataResponse(w, list)
}

//User connected, form values: branch
func getBranchUpdates(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*core.UserSession)
	branch := context.Get(r, "branch").(string)
	if !core.CanUserAccessBranch(branch, session.UserRank) {
		EasyResponse(w, core.NebErrorAuthFail, "Cannot access branch: unauthorized or branch does not exist")
		return
	}

	list, err := core.GetBranchUpdates(branch)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	}
	EasyDataResponse(w, list)
}

//User connected, must be Admin
func getCompleteBranchUpdates(w http.ResponseWriter, r *http.Request) {
	data := core.GetCompleteUpdatesInfos()
	EasyDataResponse(w, data)
}

//User connected, must be admin
func updateGitCommitCacheList(w http.ResponseWriter, r *http.Request) {
	err := core.UpdateGitCommitCache()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "updated cache list")
	}
}

//User connected, must be admin, form value: commit
func prepareGitBuild(w http.ResponseWriter, r *http.Request) {
	commit := context.Get(r, "commit").(string)
	res, err := core.PrepareGitBuild(commit)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyDataResponse(w, res)
	}
}

//User connected, must be admin, form value: commit, log
func createGitBuild(w http.ResponseWriter, r *http.Request) {
	commit := context.Get(r, "commit").(string)
	log := context.Get(r, "log").(string)
	err := core.CreateGitBuild(commit, log)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "created build")
	}
}

//User connected, must be admin, form values: build,branch, semver, log
func createUpdate(w http.ResponseWriter, r *http.Request) {
	ibuild := context.Get(r, "build").(string)
	build, err := strconv.Atoi(ibuild)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	branch := context.Get(r, "branch").(string)
	semver := context.Get(r, "semver").(string)
	log := context.Get(r, "log").(string)
	err = core.CreateUpdate(build, branch, semver, log)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "created update")
	}
}

//User connected, must be admin, form values: build,branch
func setActiveUpdate(w http.ResponseWriter, r *http.Request) {
	ibuild := context.Get(r, "build").(string)
	build, err := strconv.Atoi(ibuild)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	branch := context.Get(r, "branch").(string)

	err = core.SetActiveUpdate(branch, build)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "activated update")
	}
}

func addBranchFromBuild(w http.ResponseWriter, r *http.Request) {
	name := context.Get(r, "name").(string)
	log := context.Get(r, "log").(string)
	semver := context.Get(r, "semver").(string)
	iAccessRank := context.Get(r, "accessrank").(string)
	ibuild := context.Get(r, "build").(string)

	build, err := strconv.Atoi(ibuild)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	accessRank, err := strconv.Atoi(iAccessRank)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	err = core.AddBranchFromBuild(build, accessRank, name, log, semver)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "build created")
	}
}

func addEmptyBranch(w http.ResponseWriter, r *http.Request) {
	name := context.Get(r, "name").(string)
	iAccessRank := context.Get(r, "accessrank").(string)
	accessRank, err := strconv.Atoi(iAccessRank)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
		return
	}
	err = core.AddEmptyBranch(name, accessRank)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "build created")
	}
}
