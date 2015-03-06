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
	NebErrorPartialFail  = iota
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

func getComplexStatsTableInfos(table string) (ComplexStatTableInfo, error) {
	var values string
	var info ComplexStatTableInfo
	err := _db.QueryRow("SELECT fields, autoCount FROM neb_stats_tables WHERE tableName = ?", table).Scan(&values, &info.AutoCount)
	if err != nil {
		return info, err
	}
	info.Fields = strings.Split(values, ",")
	return info, nil
}

func getComplexStatsTablesInfos() ([]ComplexStatTableInfo, error) {
	var ret = make([]ComplexStatTableInfo, 0)

	rows, err := _db.Query("SELECT tableName, fields, autoCount FROM neb_stats_tables")
	defer rows.Close()

	if err != nil {
		return ret, err
	}

	for rows.Next() {
		var info ComplexStatTableInfo
		var fields string
		err = rows.Scan(&info.Name, fields, &info.AutoCount)
		if err != nil {
			return ret, err
		}
		info.Fields = strings.Split(fields, ",")
		ret = append(ret, info)
	}

	err = rows.Err()

	if err != nil {
		return ret, err
	}
	return ret, nil
}
