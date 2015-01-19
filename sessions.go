package main

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"github.com/nu7hatch/gouuid"
	"log"
)

func CreateSession(username string, password string) (string, error) {
	var id int
	var serverPassword, hash string

	err := _db.QueryRow("SELECT id, password, hash FROM neb_users WHERE username = ?", username).Scan(&id, &serverPassword, &hash)

	if err != nil && err == sql.ErrNoRows && _cfg["autoRegister"] == "true" { //If user are registered on connection
		c := sha512.Size
		bhash := make([]byte, c)
		_, err := rand.Read(bhash)
		if err != nil {
			log.Println("Error generating crytpo hash:", err)
			return "", err
		}
		hash := base64.URLEncoding.EncodeToString(bhash)
		hashedPassword := HashPassword(password, string(hash))
		err = RegisterUser(username, hashedPassword, string(hash))

		if err != nil {
			return "", err
		}

		return CreateSession(username, password)
	} else if err != nil && err == sql.ErrNoRows {
		return "", &NebuleuseError{NebErrorLogin, "Unknown username"}
	} else if err != nil {
		log.Println("Could not Query DB for user", username, " : ", err)
		return "", err
	}

	if HashPassword(password, hash) != serverPassword {
		return "", &NebuleuseError{NebErrorLogin, "Wrong password"}
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
func HashPassword(password string, hash string) string {
	bhash := []byte(hash)
	bpass := []byte(password)
	b := append(bhash, bpass...)

	hashed := sha512.Sum512(b)
	return base64.URLEncoding.EncodeToString(hashed[:64])
	//return string(hashed[:64])
}
func GenerateSessionId(username string) string {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Println("Failed to generate uuid:", err)
		return ""
	}
	return u4.String()
}
