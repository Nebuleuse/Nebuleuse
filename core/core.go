package core

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const NebuleuseVersion = 1

const (
	NebErrorNone         = iota
	NebError             = iota
	NebErrorDisconnected = iota
	NebErrorLogin        = iota
	NebErrorPartialFail  = iota
)

type NebuleuseError struct {
	Code int
	Msg  string
}

func (e NebuleuseError) Error() string {
	return e.Msg
}

var Cfg map[string]string
var Db *sql.DB

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init() {
	initLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	initDb()
	readConfig()

	//Todo: if update system is Git
	InitGitUpdater(".")
}
func Die() {
	Db.Close()
}

func initLogging(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func initDb() {
	db, err := sql.Open("mysql", "nebuleuse:abc@tcp(127.0.0.1:3306)/nebuleuse")

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
func readConfig() {
	var (
		name  string
		value string
	)
	Cfg = make(map[string]string)

	rows, err := Db.Query("select name, value from neb_config")
	if err != nil {
		Error.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			Error.Fatal(err)
		}
		Cfg[name] = value
	}

	err = rows.Err()
	if err != nil {
		Error.Fatal(err)
	} else {
		Info.Println("Successfully read configuration")
	}
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
