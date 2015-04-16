package handlers

import (
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"time"
)

func subscribeTo(w http.ResponseWriter, r *http.Request) {
	channel := context.Get(r, "channel").(string)
	user := context.Get(r, "user").(*core.User)

	core.Listen(channel, user.Id)
	EasyResponse(w, core.NebErrorNone, "subscribed to "+channel)
}
func unSubscribeTo(w http.ResponseWriter, r *http.Request) {
	channel := context.Get(r, "channel").(string)
	user := context.Get(r, "user").(*core.User)

	core.StopListen(channel, user.Id)
	EasyResponse(w, core.NebErrorNone, "unSubscribed from "+channel)
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	channel := context.Get(r, "channel").(string)
	message := context.Get(r, "message").(string)

	core.Dispatch(channel, message)
	EasyResponse(w, core.NebErrorNone, "Sent message ("+channel+")"+message)
}
func fetchMessage(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm["sessionid"] == nil {
		EasyResponse(w, core.NebError, "Missing sessionid")
		return
	}

	session := context.Get(r, "session").(*core.UserSession)

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
