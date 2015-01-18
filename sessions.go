package main

import (
	"errors"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"log"
	"github.com/nu7hatch/gouuid"
)

func CreateSession(username string, password string) (string, error) {
	var id int
	var serverPassword, hash string

	err := _db.QueryRow("SELECT id, password, hash FROM neb_users WHERE username = ?", username, password).Scan(&id, &serverPassword, &hash)

	if err != nil && err == sql.ErrNoRows && _cfg["autoRegister"] == "true" { //If user are registered on connection
		c := sha512.Size
		hash := make([]byte, c)
		_, err := rand.Read(hash)
		if err != nil {
			log.Println("Error generating crytpo hash:", err)
			return "", err
		}
		hashedPassword := HashPassword(password, string(hash))
		err = RegisterUser(username, hashedPassword)

		if err != nil{
			return "", err
		}

		return CreateSession(username, password)
	} else if err != nil && err == sql.ErrNoRows {
		return "", errors.New("Unknown username")
	} else if err != nil {
		log.Println("Could not Query DB for user", username, " : ", err)
		return "", err
	}

	if HashPassword(password, hash) != serverPassword {
		return "", errors.New("Wrong password")
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
func HashPassword(password string, hash string) string{
	bhash := []byte(hash)
	bpass := []byte(password)
	b := append(bhash, bpass...)

	hashed := sha512.Sum512(b)
	return string(hashed[:64])
}
func GenerateSessionId(username string) string{
	u4, err := uuid.NewV4()
	if err != nil {
	    log.Println("Failed to generate uuid:", err)
	    return ""
	}
	return u4.String()
}