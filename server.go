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

	core.Info.Fatal(http.ListenAndServe(core.Cfg.GetSysConfig("serverAddress")+":"+core.Cfg.GetSysConfig("serverPort"), nil))
}
func createInstallServer() {
	handlers.RegisterInstallHandlers()
	core.Info.Fatal(http.ListenAndServe(core.Cfg.GetSysConfig("serverAddress")+":"+core.Cfg.GetSysConfig("serverPort"), nil))
}
func sessionsPurgeTimer() {
	timer := time.NewTimer(time.Minute)
	<-timer.C

	core.PurgeSessions()
	go sessionsPurgeTimer()
}
