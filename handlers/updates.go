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
	Updates []core.Update
}

func getUpdateGraphList(w http.ResponseWriter, r *http.Request) {
	list, err := core.GetUpdatesInfos(0)
	if err != nil {
		EasyErrorResponse(w, core.NebError, err)
	} else {
		var response getupdateGraphListResponse
		response.Updates = list
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
