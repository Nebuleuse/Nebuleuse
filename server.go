package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	NebErrorNone         = iota
	NebError             = iota
	NebErrorDisconnected = iota
	NebErrorLogin        = iota
)

type NebuleuseError struct {
	code int
	msg  string
}

func (e NebuleuseError) Error() string {
	return e.msg
}

func createServer() {
	registerHandlers()

	go SessionsPurgeTimer()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func SessionsPurgeTimer() {
	PurgeSessions()

	timer := time.NewTimer(time.Minute)
	<-timer.C
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
func getFieldListForStatTable(table string) ([]string, error) {
	var values string
	err := _db.QueryRow("SELECT fields FROM neb_stats_tables WHERE tableName = ?", table).Scan(&values)
	if err != nil {
		return nil, err
	}
	return strings.Split(values, ","), nil
}
