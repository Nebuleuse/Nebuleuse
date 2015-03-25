package handlers

import (
	"github.com/Nebuleuse/Nebuleuse/core"
	"io"
	"net/http"
)

var messages chan string

func InitLongPoll() {
	messages = make(chan string, 100)
}
func sendMessage(w http.ResponseWriter, r *http.Request) {
	ms := r.FormValue("msg")
	messages <- ms
	io.WriteString(w, "done")
}
func longPollRequest(w http.ResponseWriter, r *http.Request) {
	core.Info.Println("poll")
	io.WriteString(w, <-messages)
}
