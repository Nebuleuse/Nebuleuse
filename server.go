package main

import (
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/Nebuleuse/Nebuleuse/handlers"
	"net/http"
	"time"
)

func CreateServer() {
	handlers.RegisterHandlers()

	go SessionsPurgeTimer()

	core.Info.Fatal(http.ListenAndServe(core.SysCfg["serverAddress"]+":"+core.SysCfg["serverPort"], nil))
}

func SessionsPurgeTimer() {
	core.PurgeSessions()

	timer := time.NewTimer(time.Minute)
	<-timer.C
	go SessionsPurgeTimer()
}
