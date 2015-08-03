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
func getMessages(w http.ResponseWriter, r *http.Request) {
	session := context.Get(r, "session").(*core.UserSession)

	if session == nil {
		EasyResponse(w, core.NebErrorDisconnected, "No session found")
		return
	}

	// If we are already polling messages, signal previous poll that they can stop
	if session.LongPolling {
		session.TimedOut <- 1
	}

	session.LongPolling = true
	session.Heartbeat()
	select {
	case msg := <-session.Messages:
		fmt.Fprint(w, msg)
	case <-session.TimedOut:
		return
	case <-time.After(time.Second * time.Duration(core.GetSysConfigInt("LongpollingTimeout"))):
		EasyResponse(w, core.NebErrorNone, "longpoll timedout")
	}
	session.LongPolling = false
}
