package handlers

import (
	"fmt"
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/gorilla/context"
	"net/http"
	"time"
)

//User connected, form value: pipe, channel
func subscribeTo(w http.ResponseWriter, r *http.Request) {
	pipe := context.Get(r, "pipe").(string)
	channel := context.Get(r, "channel").(string)
	session := context.Get(r, "session").(*core.UserSession)

		core.Listen(pipe, channel, session)
	EasyResponse(w, core.NebErrorNone, "subscribed to "+channel)
}


//User connected, form value: pipe, channel
func unSubscribeTo(w http.ResponseWriter, r *http.Request) {
	pipe := context.Get(r, "pipe").(string)
	channel := context.Get(r, "channel").(string)
	session := context.Get(r, "session").(*core.UserSession)

	core.StopListen(pipe, channel, session)
	EasyResponse(w, core.NebErrorNone, "unSubscribed from "+channel)
}

//User connected, form value: pipe, channel, message
func sendMessage(w http.ResponseWriter, r *http.Request) {
	pipe := context.Get(r, "pipe").(string)
	channel := context.Get(r, "channel").(string)
	message := context.Get(r, "message").(string)

	core.Dispatch(pipe, channel, message)
	EasyResponse(w, core.NebErrorNone, "Sent message ("+channel+")"+message)
}

//User connected
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
		core.Info.Println(msg)
		fmt.Fprint(w, msg)
	case <-session.TimedOut:
		return
	case <-time.After(time.Second * time.Duration(core.Cfg.GetSysConfigInt("LongpollingTimeout"))):
		EasyResponse(w, core.NebErrorNone, "longpoll timedout")
	}
	session.LongPolling = false
}
