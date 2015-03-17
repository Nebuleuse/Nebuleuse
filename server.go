package main

import (
	"Nebuleuse/core"
	"Nebuleuse/handlers"
	"net/http"
	"time"
)

func CreateServer() {
	handlers.RegisterHandlers()

	go SessionsPurgeTimer()

	core.Info.Fatal(http.ListenAndServe(":8080", nil))
}

func SessionsPurgeTimer() {
	core.PurgeSessions()

	timer := time.NewTimer(time.Minute)
	<-timer.C
	go SessionsPurgeTimer()
}
