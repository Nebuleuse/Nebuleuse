package handlers

import (
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"net/http"
)

func sendMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil && r.PostForm["channel"] == nil && r.PostForm["message"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid"))
		return
	}

	_, err := core.GetUserBySession(r.FormValue("sessionid"), core.UserMaskOnlyId)

	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	core.Dispatch(r.FormValue("channel"), r.FormValue("message"))
}
func longPollRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid"))
		return
	}

	user, err := core.GetUserBySession(r.FormValue("sessionid"), core.UserMaskOnlyId)

	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	fmt.Fprint(w, <-core.GetMessages(user.Id))
	//io.WriteString(w, <-messages)
}
