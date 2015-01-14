package main

import (
	"log"
	"github.com/nu7hatch/gouuid"
)

func CreateSession(username string, password string) (string, error) {
	var id int
	
	err := _db.QueryRow("SELECT id FROM neb_users WHERE username = ? AND password = ?", username, password).Scan(&id)

	if err != nil && err.Error() == "sql: no rows in result set" && _cfg["autoRegister"] == "true" { //If user are registered on connection
		err = RegisterUser(username, password)

		if err != nil{
			return "", err
		}

		return CreateSession(username, password)
	} else if err != nil {
		log.Println("Could not Query DB for user", username, " : ", err)
		return "", err
	}

	sessionid := GenerateSessionId(username)

	stmt, err := _db.Prepare("REPLACE INTO neb_sessions (userid,lastAlive,sessionId,sessionStart) VALUES (?,NOW(),?,NOW())")
	_, err = stmt.Exec(id, sessionid)
	if err != nil {
		log.Println("Could not insert session :", err)
		return "", err
	}

	return sessionid, nil
}
func GenerateSessionId(username string) string{
	u4, err := uuid.NewV4()
	if err != nil {
	    log.Println("Failed to generate uuid:", err)
	    return ""
	}
	return u4.String()
}