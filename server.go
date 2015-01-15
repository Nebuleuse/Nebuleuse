package main

import (
	"time"
	"log"
	"net/http"
)

const (  
	NebErrorNone = iota
	NebError = iota
	NebErrorDisconnected = iota
	NebErrorLogin = iota
)

func createServer(){
	registerHandlers()

	go SessionsPurgeTimer()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SessionsPurgeTimer(){
	PurgeSessions()

	timer := time.NewTimer(time.Minute)
    <- timer.C
    go SessionsPurgeTimer()
}
func PurgeSessions() {
	stmt, err := _db.Prepare("DELETE FROM neb_sessions WHERE NOW() > Date_Add( lastAlive, INTERVAL ? SECOND )")
	if err != nil {
		log.Println("Failed to prepare statement : ", err)
		return
	}
	res, err := stmt.Exec(_cfg["sessionTimeout"])
	if err != nil {
		log.Println("Failed to purge sessions: ", err)
		return	
	}
	af, err := res.RowsAffected()
	if err != nil {
		log.Println("Failed to get sessions affected rows :", err)
		return
	}
	if af > 0 {
		log.Println("Purged ", af, " sessions")
	}
}