package handlers

import (
	"encoding/json"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"strconv"
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

func addUpdate(w http.ResponseWriter, r *http.Request) {
	data := context.Get(r, "data").([]byte)

	var request core.Update
	err := json.Unmarshal([]byte(data), &request)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	}
	err = core.AddUpdate(request)
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
