package handlers

import (
	"encoding/json"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
)

func getUpdateList(w http.ResponseWriter, r *http.Request) {
	version := context.Get(r, "version").(int)

	list, err := core.GetUpdatesInfos(version)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyDataResponse(w, list)
	}
}

type getupdateGraphListResponse struct {
	Updates       []core.Update
	Commits       []core.Commit
	CurrentCommit string
}

func getUpdateListWithGit(w http.ResponseWriter, r *http.Request) {
	withDiffs := context.Get(r, "diffs").(bool)

	list, err := core.GetUpdatesInfos(0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		commits, err := core.GetGitCommitList()
		if err != nil {
			EasyErrorResponse(w, core.NebError, err)
			return
		}
		var response getupdateGraphListResponse
		response.Updates = list
		response.CurrentCommit, err = core.GetCurrentCommit()
		if err != nil {
			EasyErrorResponse(w, core.NebError, err)
		}
		response.Commits = commits
		if !withDiffs {
			for i, _ := range response.Commits {
				core.Info.Println(response.Commits[i].Message)
				response.Commits[i].Diff = nil
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

func updateGitCommitCacheList(w http.ResponseWriter, r *http.Request) {
	err := core.UpdateGitCommitCache()
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyResponse(w, core.NebErrorNone, "updated cache list")
	}
}

func prepareGitPatch(w http.ResponseWriter, r *http.Request) {
	commit := context.Get(r, "commit").(string)
	res, err := core.PrepareGitPatch(commit)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		EasyDataResponse(w, res)
	}
}
