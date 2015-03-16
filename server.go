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
	PurgeSessions()

	timer := time.NewTimer(time.Minute)
	<-timer.C
	go SessionsPurgeTimer()
}

func PurgeSessions() {
	stmt, err := core.Db.Prepare("DELETE FROM neb_sessions WHERE NOW() > Date_Add( lastAlive, INTERVAL ? SECOND )")
	if err != nil {
		core.Warning.Println("Failed to prepare statement : ", err)
		return
	}
	res, err := stmt.Exec(core.Cfg["sessionTimeout"])
	if err != nil {
		core.Warning.Println("Failed to purge sessions: ", err)
		return
	}
	af, err := res.RowsAffected()
	if err != nil {
		core.Warning.Println("Failed to get sessions affected rows :", err)
		return
	}
	if af > 0 {
		core.Info.Println("Purged ", af, " sessions")
	}
}
