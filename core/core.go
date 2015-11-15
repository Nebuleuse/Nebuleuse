package core

import (
	"bufio"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/config"
	"os"
	"strconv"
	"strings"
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
		initGit()
	}
}
func Die() {
	Db.Close()
	logFile.Close()
}
func Install() {
	initLogging()

	reader := bufio.NewReader(os.Stdin)
	if _, err := os.Stat(".config"); err == nil {
		Info.Println(".config file already exists. Replace it ? Y/N")

		in, err := reader.ReadString('\n')
		in = strings.TrimSpace(in)
		for err != nil || (in != "Y" && in != "N") {
			in, err = reader.ReadString('\n')
			in = strings.TrimSpace(in)
		}
		if in == "N" {
			return
		}
	}
	if _, err := os.Stat("nebuleuse.sql"); os.IsNotExist(err) {
		Error.Println("file nebuleuse.sql not found")
		return
	}

	defaultOptions := map[string]string{"serverAddress": "127.0.0.1",
		"serverPort":               "8080",
		"dbType":                   "mysql",
		"dbAddress":                "127.0.0.1:3306",
		"dbUser":                   "",
		"dbPass":                   "",
		"dbBase":                   "",
		"MaxSessionsChannelBuffer": "10",
		"LongpollingTimeout":       "10",
		"DashboardLocation":        "./admin/"}

	c := config.NewDefault()
	Info.Println("Please enter the following configuration values. Enter empty value to use the default one:")
	for option, val := range defaultOptions {
		outline := option
		canDefault := false
		if val != "" {
			outline += "(Default: " + val + ")"
			canDefault = true
		}
		Info.Print(outline, ":")

		in, _ := reader.ReadString('\n')
		in = strings.TrimSpace(in)
		for in == "" && !canDefault {
			in, _ = reader.ReadString('\n')
			in = strings.TrimSpace(in)
		}
		if canDefault && in == "" {
			in = val
		}

		c.AddOption("default", option, in)
	}

	c.WriteFile(".config", 0644, "")
	Info.Println("Saved config")
	Info.Println("Testing database connectivity")
	initConfig()
	initDb()

	Info.Println("Do you want to create a new admin account ? Y/N")

	in, err := reader.ReadString('\n')
	in = strings.TrimSpace(in)
	for err != nil || (in != "Y" && in != "N") {
		in, err = reader.ReadString('\n')
		in = strings.TrimSpace(in)
	}
	if in == "N" {
		return
	}

	Info.Println("Please enter the user name : ")
	username, _ := reader.ReadString('\n')
	username = strings.TrimSpace(username)
	for username == "" {
		username, _ = reader.ReadString('\n')
		username = strings.TrimSpace(username)
	}

	Info.Println("Please enter the password : ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)
	for password == "" {
		password, _ = reader.ReadString('\n')
		password = strings.TrimSpace(password)
	}

	err = RegisterUser(username, password, 3)
	if err != nil {
		Error.Println("Could not register user : ", err)
		return
	}

	Info.Println("Registered user : ", username)

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
func GetGameSemVer() string {
	update, err := GetUpdateInfos(GetGameVersion())
	if err != nil {
		return ""
	}
	return update.SemVer
}
func GetUpdaterVersion() int {
	n, e := strconv.Atoi(Cfg["updaterVersion"])
	if e != nil {
		return -1
	}
	return n
}
