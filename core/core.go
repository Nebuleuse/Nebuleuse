package core

import (
	"bufio"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/robfig/config"
	"os"
	"strings"
	"sync"
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

type configMgr struct {
	Cfg        map[string]string
	SysCfg     map[string]string
	configLock sync.RWMutex
}

var Cfg configMgr
var Db *sql.DB

func Init() {
	initLogging()
	Cfg.InitConfig()
	initDb()
	Cfg.LoadConfig()
	initSessions()
	initMessaging()
	err := initUpdateSystem()
	if err != nil {
		Warning.Println("Update system not setup: " + err.Error())
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
		in = strings.ToUpper(strings.TrimSpace(in))
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
			Info.Println("No default, please enter a value: ")
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
	Cfg.InitConfig()
	initDb()

	Info.Println("Do you want to create a new admin account ? Y/N")

	in, err := reader.ReadString('\n')
	in = strings.ToUpper(strings.TrimSpace(in))
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

	err = RegisterUser(username, password, UserRankAdmin)
	if err != nil {
		Error.Println("Could not register user : ", err)
		return
	}

	Info.Println("Registered user : ", username)

}

func initDb() {
	con := Cfg.GetSysConfig("dbUser") + ":" + Cfg.GetSysConfig("dbPass") + "@tcp(" + Cfg.GetSysConfig("dbAddress") + ")/" + Cfg.GetSysConfig("dbBase") + "?parseTime=true"
	db, err := sql.Open(Cfg.GetSysConfig("dbType"), con)

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
func GetUpdaterVersion() int {
	return Cfg.GetConfigInt("updaterVersion")
}

func getFileSize(path string) (int64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return 0, err
	}
	return stat.Size(), nil

}
