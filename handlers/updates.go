package handlers

import (
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
func addUpdate(w http.ResponseWriter, r *http.Request) {
	
}