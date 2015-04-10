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
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or channel"))
		return
	}
	channel := r.FormValue("channel")
	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
	}

	core.Listen(channel, user.Id)
	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "subscribed to "+channel))
}
func unSubscribeTo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil && r.PostForm["channel"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or channel"))
		return
	}
	channel := r.FormValue("channel")
	user, err := core.GetUserBySession(r.PostForm["sessionid"][0], core.UserMaskOnlyId)
	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
	}
	core.StopListen(channel, user.Id)
	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "unSubscribed from "+channel))
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil && r.PostForm["channel"] == nil && r.PostForm["message"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid or channel or message"))
		return
	}

	_, err := core.GetUserBySession(r.FormValue("sessionid"), core.UserMaskOnlyId)

	if err != nil {
		fmt.Fprint(w, EasyErrorResponse(core.NebErrorDisconnected, err))
		return
	}

	channel := r.FormValue("channel")
	message := r.FormValue("message")
	core.Dispatch(channel, message)
	fmt.Fprint(w, EasyResponse(core.NebErrorNone, "Sent message ("+channel+")"+message))
}
func longPollRequest(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		fmt.Fprint(w, EasyResponse(core.NebError, "Missing sessionid"))
		return
	}

	session := core.GetSessionBySessionId(r.FormValue("sessionid"))
	if session == nil {
		fmt.Fprint(w, EasyResponse(core.NebErrorDisconnected, "Could not get session data using session Id: "+r.FormValue("sessionid")))
		return
	}

	session.LongPolling = true
	session.Heartbeat()
	select {
	case msg := <-session.Messages:
		fmt.Fprint(w, msg)

	case <-time.After(time.Second * time.Duration(core.GetConfigInt("LongpollingTimeout"))):
		fmt.Fprint(w, EasyResponse(core.NebErrorNone, "longpoll timedout"))
		return
	}
	session.LongPolling = false
}
