package main

import (
	"strconv"
	"log"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

const nebuleuseVersion = 1
var _cfg map[string]string
var _db *sql.DB

func main() {
	initDb()
	defer _db.Close()

	readConfig()

	createServer()
}

func initDb(){
	db, err := sql.Open("mysql", "nebuleuse:abc@tcp(127.0.0.1:3306)/nebuleuse")

	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully connected to db")
	}
	_db = db
}

func readConfig(){
	var (name string
		value string)
	_cfg = make(map[string]string)

	rows, err := _db.Query("select name, value from neb_config")
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&name, &value)
		if err != nil {
			log.Fatal(err)
		}
		_cfg[name] = value
	}

	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully read configuration")
	}
}

func getGameVersion() int{
	n, e := strconv.Atoi(_cfg["gameVersion"])
	if(e != nil){
		return -1
	}
	return n
}
func getUpdaterVersion() int{
	n, e := strconv.Atoi(_cfg["updaterVersion"])
	if(e != nil){
		return -1
	}
	return n
}