package main

import (
	"Nebuleuse/gitUpdater"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
)

const nebuleuseVersion = 1

var _cfg map[string]string
var _db *sql.DB

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func main() {
	InitLogging(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	initDb()
	defer _db.Close()
	gitUpdater.InitGit(".")

	readConfig()

	createServer()
}

func InitLogging(traceHandle io.Writer, infoHandle io.Writer, warningHandle io.Writer, errorHandle io.Writer) {
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
	_db = db
}
func readConfig() {
	var (
		name  string
		value string
	)
	_cfg = make(map[string]string)

	rows, err := _db.Query("select name, value from neb_config")
	if err != nil {
		Error.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			Error.Fatal(err)
		}
		_cfg[name] = value
	}

	err = rows.Err()
	if err != nil {
		Error.Fatal(err)
	} else {
		Info.Println("Successfully read configuration")
	}
}

func getGameVersion() int {
	n, e := strconv.Atoi(_cfg["gameVersion"])
	if e != nil {
		return -1
	}
	return n
}
func getUpdaterVersion() int {
	n, e := strconv.Atoi(_cfg["updaterVersion"])
	if e != nil {
		return -1
	}
	return n
}
