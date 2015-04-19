package core

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

const NebuleuseVersion = 1

const (
	NebErrorNone = iota
	// Generic error
	NebError = iota
	// Session is dead, you were disconnected
	NebErrorDisconnected = iota
	// Login failed
	NebErrorLogin = iota
	// There were errors during multiple operations
	NebErrorPartialFail = iota
	// User is not authorized to do that
	NebErrorAuthFail = iota
)

type NebuleuseError struct {
	Code int
	Msg  string
}

func (e NebuleuseError) Error() string {
	return e.Msg
}

type ConfigMgr map[string]string

var Cfg ConfigMgr
var SysCfg map[string]string
var Db *sql.DB

func Init() {
	initLogging()
	initConfig()
	initDb()
	loadConfig()
	initSessions()
	initMessaging()

	//Todo: if update system is Git
	if Cfg["updateSystem"] == "GitPatch" || Cfg["updateSystem"] == "FullGit" {
		InitGitUpdater(SysCfg["gitPath"])
	}
}
func Die() {
	Db.Close()
	logFile.Close()
}

func initDb() {
	con := SysCfg["dbUser"] + ":" + SysCfg["dbPass"] + "@tcp(" + SysCfg["dbAddress"] + ")/" + SysCfg["dbBase"] + "?parseTime=true"
	db, err := sql.Open(SysCfg["dbType"], con)

	if err != nil {
		Error.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		Error.Fatal(err)
	} else {
		Info.Println("Successfully connected to db")
	}
	Db = db
}
func GetGameVersion() int {
	n, e := strconv.Atoi(Cfg["gameVersion"])
	if e != nil {
		return -1
	}
	return n
}
func GetUpdaterVersion() int {
	n, e := strconv.Atoi(Cfg["updaterVersion"])
	if e != nil {
		return -1
	}
	return n
}
