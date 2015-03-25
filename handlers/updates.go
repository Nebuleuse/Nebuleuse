package handlers

import (
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"net/http"
)

type getUpdateListRequest struct {
	Version int
}

func getUpdateList(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil || r.PostForm["data"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or data"))
		return
	}
	//user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)
}
