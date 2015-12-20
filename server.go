package main

import (
	"github.com/Nebuleuse/Nebuleuse/core"
	"github.com/Nebuleuse/Nebuleuse/handlers"
	"net/http"
	"time"
)

func createServer() {
	handlers.RegisterHandlers()

	go sessionsPurgeTimer()

	core.Info.Fatal(http.ListenAndServe(core.SysCfg["serverAddress"]+":"+core.SysCfg["serverPort"], nil))
}

func sessionsPurgeTimer() {
	timer := time.NewTimer(time.Minute)
	<-timer.C

	core.PurgeSessions()
	go sessionsPurgeTimer()
}
