package handlers

import (
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"net/http"
	"time"
)

func subscribeTo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil && r.PostForm["channel"] == nil {
		EasyResponse(w, core.NebError, "Missing sessionid or channel")
		return
	}
	channel := r.FormValue("channel")
	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)
	if err != nil {
		EasyErrorResponse(w, core.NebErrorDisconnected, err)
	}

	core.Listen(channel, user.Id)
	EasyResponse(w, core.NebErrorNone, "subscribed to "+channel)
}
func unSubscribeTo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil && r.PostForm["channel"] == nil {
		EasyResponse(w, core.NebError, "Missing sessionid or channel")
		return
	}
	channel := r.FormValue("channel")
	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)
	if err != nil {
		EasyErrorResponse(w, core.NebErrorDisconnected, err)
	}
	core.StopListen(channel, user.Id)
	EasyResponse(w, core.NebErrorNone, "unSubscribed from "+channel)
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil && r.PostForm["channel"] == nil && r.PostForm["message"] == nil {
		EasyResponse(w, core.NebError, "Missing sessionid or channel or message")
		return
	}

	_, err := core.GetUserBySession(r.FormValue("sessionid"), core.UserMaskOnlyId)

	if err != nil {
		EasyErrorResponse(w, core.NebErrorDisconnected, err)
		return
	}

	channel := r.FormValue("channel")
	message := r.FormValue("message")
	core.Dispatch(channel, message)
	EasyResponse(w, core.NebErrorNone, "Sent message ("+channel+")"+message)
}
func longPollRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		EasyResponse(w, core.NebError, "Missing sessionid")
		return
	}

	session := core.GetSessionBySessionId(r.FormValue("sessionid"))
	if session == nil {
		EasyResponse(w, core.NebErrorDisconnected, "Could not get session data using session Id: "+r.FormValue("sessionid"))
		return
	}

	session.LongPolling = true
	session.Heartbeat()
	select {
	case msg := <-session.Messages:
		fmt.Fprint(w, msg)

	case <-time.After(time.Second * time.Duration(core.GetConfigInt("LongpollingTimeout"))):
		EasyResponse(w, core.NebErrorNone, "longpoll timedout")
		return
	}
	session.LongPolling = false
}
