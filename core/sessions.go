package core

import (
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"github.com/nu7hatch/gouuid"
	"time"
)

type UserSession struct {
	LongPolling bool `json:"-"`
	LastAlive   time.Time
	SessionId   string
	UserId      int
	Messages    chan string `json:"-"`
	TimedOut    chan int    `json:"-"`
}

var connectedUsers map[int]UserSession

func initSessions() error {
	connectedUsers = make(map[int]UserSession)
	rows, err := Db.Query("SELECT userid, lastAlive, sessionId FROM neb_sessions")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var userid int
		var lastAlive time.Time
		var sessionId string
		err := rows.Scan(&userid, &lastAlive, &sessionId)
		if err != nil {
			Error.Println("Could not read sessions from db: " + err.Error())
			return err
		}

		var session UserSession
		session.SessionId = sessionId
		session.Messages = make(chan string, GetConfigInt("MaxSessionsChannelBuffer"))
		session.LastAlive = lastAlive
		session.LongPolling = false
		session.UserId = userid
		connectedUsers[userid] = session
	}
	return nil
}
func IsUserLongPolling(userid int) bool {
	return connectedUsers[userid].LongPolling
}
func GetSessionByUserId(userid int) *UserSession {
	session, ok := connectedUsers[userid]
	if ok {
		return &session
	}
	return nil
}
func GetSessionBySessionId(sessionid string) *UserSession {
	for _, session := range connectedUsers {
		if session.SessionId == sessionid {
			return &session
		}
	}
	return nil
}
func DisconnectUser(userid int) {
	delete(connectedUsers, userid)
	stmt, err := Db.Prepare("DELETE FROM neb_sessions WHERE userid = ?")
	if err != nil {
		Error.Println("Could not prepare statement: ", err)
		return
	}
	_, err = stmt.Exec(userid)
	if err != nil {
		Error.Println("Could not delete user session :", err)
	}
}
func CountOnlineUsers() int {
	return len(connectedUsers)
}
func GetOnlineUsersList() []UserSession {
	var list []UserSession

	for _, user := range connectedUsers {
		list = append(list, user)
	}

	return list
}
func (s *UserSession) Heartbeat() {
	s.LastAlive = time.Now()

	stmt, err := Db.Prepare("UPDATE neb_sessions SET lastAlive = NOW() WHERE userid = ?")
	if err != nil {
		Error.Println("Could not prepare statement ", s.UserId, ": ", err)
		return
	}
	_, err = stmt.Exec(s.UserId)
	if err != nil {
		Error.Println("Could not Heartbeat user ", s.UserId, ": ", err)
		return
	}
}

func CreateSession(username string, password string) (string, error) {
	var id int
	var serverPassword, hash string

	err := Db.QueryRow("SELECT id, password, hash FROM neb_users WHERE username = ?", username).Scan(&id, &serverPassword, &hash)

	if err != nil && err == sql.ErrNoRows { //If user are registered on connection
		if Cfg["autoRegister"] == "true" {
			RegisterUser(username, password, 1)
			return CreateSession(username, password)
		} else {
			return "", &NebuleuseError{NebErrorLogin, "Unknown username"}
		}
	} else if err != nil {
		Error.Println("Could not Query DB for user", username, " : ", err)
		return "", err
	}

	if HashPassword(password, hash) != serverPassword {
		return "", &NebuleuseError{NebErrorLogin, "Wrong password"}
	}

	sessionid := GenerateSessionId(username)

	stmt, err := Db.Prepare("REPLACE INTO neb_sessions (userid,lastAlive,sessionId,sessionStart) VALUES (?,NOW(),?,NOW())")
	_, err = stmt.Exec(id, sessionid)
	if err != nil {
		Error.Println("Could not insert session :", err)
		return "", err
	}

	//Create entry in connectedUsers
	var session UserSession
	session.SessionId = sessionid
	session.Messages = make(chan string, GetConfigInt("MaxSessionsChannelBuffer"))
	session.LastAlive = time.Now()
	session.LongPolling = false
	session.UserId = id
	connectedUsers[id] = session

	return sessionid, nil
}
func PurgeSessions() {
	for id, sess := range connectedUsers {
		delta := time.Since(sess.LastAlive).Minutes()
		if delta > GetConfigFloat("sessionTimeout") {
			delete(connectedUsers, id)
		}
	}
	stmt, err := Db.Prepare("DELETE FROM neb_sessions WHERE NOW() > Date_Add( lastAlive, INTERVAL ? SECOND )")
	if err != nil {
		Error.Println("Failed to prepare statement : ", err)
		return
	}
	res, err := stmt.Exec(Cfg["sessionTimeout"])
	if err != nil {
		Error.Println("Failed to purge sessions: ", err)
		return
	}
	af, err := res.RowsAffected()
	if err != nil {
		Error.Println("Failed to get sessions affected rows :", err)
		return
	}
	if af > 0 {
		Info.Println("Purged ", af, " sessions")
	}
}

func HashPassword(password string, hash string) string {
	bhash := []byte(hash)
	bpass := []byte(password)
	b := append(bhash, bpass...)

	hashed := sha512.Sum512(b)
	return base64.URLEncoding.EncodeToString(hashed[:64])
}
func GenerateSessionId(username string) string {
	u4, err := uuid.NewV4()
	if err != nil {
		Error.Println("Failed to generate uuid:", err)
		return ""
	}
	return u4.String()
}
