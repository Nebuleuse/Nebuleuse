package handlers

import (
	"encoding/json"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
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

func getUpdateList(w http.ResponseWriter, r *http.Request) {
	version := context.Get(r, "version").(int)

	list, err := core.GetUpdatesInfos(version)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyDataResponse(w, list)
	}
}

type getUpdateListResponse struct {
	UpdateSystem   string
	Updates        []core.Update
	Commits        []core.Commit
	CurrentVersion int
}

//User connected, must be admin, optional switch : diffs, POST
func getUpdateListComplete(w http.ResponseWriter, r *http.Request) {
	withDiffs := context.Get(r, "diffs").(bool)

	list, err := core.GetUpdatesInfos(0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		var response getUpdateListResponse
		response.UpdateSystem = core.GetUpdateSystem()
		response.Updates = list
		response.CurrentVersion = core.GetCurrentVersion()

		if response.UpdateSystem == "FullGit" || response.UpdateSystem == "GitPatch" {
			commits, err := core.GetGitCommitList()
			if err != nil {
				EasyErrorResponse(w, core.NebError, err)
				return
			}
			response.Commits = commits
			if !withDiffs {
				for i := range response.Commits {
					response.Commits[i].Diff = nil
				}
			}
		}

		EasyDataResponse(w, response)
	}

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
func prepareGitPatch(w http.ResponseWriter, r *http.Request) {
	commit := context.Get(r, "commit").(string)
	res, err := core.PrepareGitPatch(commit)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyDataResponse(w, res)
	}
}
