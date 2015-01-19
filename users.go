package main

import (
	"database/sql"
	"log"
)

//Achievements
type Achievement struct {
	id       int
	Name     string
	Progress uint
	Value    uint
}

func (a *Achievement) isComplete() bool {
	return a.Progress == a.Value
}

// Stats
type Stat struct {
	Name  string
	Value int64
}

// Kills
type Kill struct {
	Weapon  string
	X, Y, Z int
}

// User
type User struct {
	id           int
	Username     string
	SessionId    string
	Rank         int
	Avatar       string
	Achievements []Achievement
	Stats        []Stat
}

func GetUser(SessionId string) (*User, error) {
	var user User
	user.SessionId = SessionId

	var id int
	err := _db.QueryRow("SELECT userid FROM neb_sessions WHERE sessionid = ?", SessionId).Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return nil, &NebuleuseError{NebErrorDisconnected, "No session found"}
	} else if err != nil {
		return nil, err
	}
	user.id = id

	err = _db.QueryRow("SELECT username, rank, avatars FROM neb_users WHERE id = ?", id).Scan(&user.Username, &user.Rank, &user.Avatar)
	if err != nil && err == sql.ErrNoRows {
		return nil, &NebuleuseError{NebErrorDisconnected, "No user found"}
	} else if err != nil {
		return nil, err
	}
	if user.Avatar == "" {
		user.Avatar = _cfg["defaultAvatar"]
	}

	user.PopulateAchievements()
	user.PopulateStats()

	return &user, nil
}

func RegisterUser(username string, password string, hash string) error {
	stmt, err := _db.Prepare("INSERT INTO neb_users (username,password,rank,hash) VALUES (?,?,1,?)")
	_, err = stmt.Exec(username, password, hash)
	if err != nil {
		log.Println("Could not register new user :", err)
		return err
	}

	return nil
}

func (u *User) PopulateAchievements() error {
	rows, err := _db.Query("SELECT achievementid, progress, name, max FROM neb_users_achievements LEFT JOIN neb_achievements ON (neb_achievements.id = neb_users_achievements.achievementid) WHERE neb_users_achievements.userid = ?", u.id)
	if err != nil {
		log.Println("Could not get user achievements :", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ach Achievement
		err := rows.Scan(&ach.id, &ach.Progress, &ach.Name, &ach.Value)
		if err != nil {
			log.Println("Could not get user achievements :", err)
			return err
		}
		u.Achievements = append(u.Achievements, ach)
	}

	err = rows.Err()
	if err != nil {
		log.Println("Could not get user achievements :", err)
		return err
	}
	return nil
}
func (u *User) PopulateStats() error {
	rows, err := _db.Query("SELECT name, value FROM neb_users_stats WHERE userid = ?", u.id)
	if err != nil {
		log.Println("Could not get user stats :", err)
		return err
	}
	defer rows.Close()

	if err != nil {
		log.Println("Could not get columns:", err)
		return err
	}

	for rows.Next() {
		var st Stat
		err := rows.Scan(&st.Name, &st.Name)
		if err != nil {
			log.Println("Could not get user Stats :", err)
			return err
		}
		u.Stats = append(u.Stats, st)
	}

	err = rows.Err()
	if err != nil {
		log.Println("Could not get user Stats :", err)
		return err
	}
	return nil
}
func (u *User) Heartbeat() {
	stmt, err := _db.Prepare("UPDATE neb_sessions SET lastAlive = NOW() WHERE userid = ?")
	_, err = stmt.Exec(u.id)
	if err != nil {
		log.Println("Could not Heartbeat :", err)
	}
}
func (u *User) UpdateAchievementProgress(aid int, value int) error {
	stmt, err := _db.Prepare("UPDATE neb_users_achievements SET progress= ? WHERE userid = ? AND achievementid = ? LIMIT 1")
	if err != nil {
		log.Println("Could not create statement : ", err)
		return err
	}

	res, err := stmt.Exec(value, u.id, aid)
	if err != nil {
		log.Println("Could not update achievement :", err)
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Println("Could not get update rowcount :", err)
		return err
	}
	if rowCnt == 0 {
		log.Println("Tried to update achievementid : ", aid, " but no rows affected")
		return &NebuleuseError{NebError, "No rows affected by operation"}
	}

	return nil
}
func (u *User) UpdateStats(stats []Stat) error {
	stmt, err := _db.Prepare("UPDATE neb_users_stats SET value = ? WHERE userid = ? AND name = ? LIMIT 1")
	if err != nil {
		log.Println("Could not create statement : ", err)
		return err
	}

	for _, stat := range stats {
		if stat.Name == "userid" {
			log.Println("Could not update user stats, userid present in stat list")
			return &NebuleuseError{NebErrorLogin, "Could not update user stats, userid present in stat list"}
		}
		_, err := stmt.Exec(stat.Value, u.id, stat.Name)
		if err != nil {
			log.Println("Could not update user stats : ", err)
			return err
		}
	}
	return nil
}
func (u *User) InsertKills(kills []Kill, mapName string) error {
	stmt, err := _db.Prepare("INSERT INTO neb_users_stats_kills VALUES (?,?,?,?,?,?)")
	if err != nil {
		log.Println("Could not create statement : ", err)
		return err
	}
	for _, kill := range kills {
		_, err := stmt.Exec(u.id, kill.X, kill.Y, kill.Z, kill.Weapon, mapName)
		if err != nil {
			log.Println("Could not insert kill : ", err)
			return err
		}
	}
	return nil
}
