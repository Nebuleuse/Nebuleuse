package main

import (
	"log"
)

//Achievements
type Achievement struct{
	id int
	Name string
	Progress uint
	Value uint
}
func (a *Achievement) isComplete() bool{
	return a.Progress == a.Value
}

// Stats
type Stat struct{
	Name string
	Value int64
}

// User
type User struct{
	id int
	Username string
	SessionId string
	Rank int
	Achievements[] Achievement
	Stats[] Stat
}
func GetUser(SessionId string) (*User, error){
	var user User
	user.SessionId = SessionId

	var id int
	err := _db.QueryRow("SELECT userid FROM neb_sessions WHERE sessionid = ?", SessionId).Scan(&id)
	if err != nil {
		return nil, err
	}
	user.id = id

	err = _db.QueryRow("SELECT username, rank FROM neb_users WHERE id = ?", id).Scan(&user.Username, &user.Rank)
	if err != nil {
		return nil, err
	}

	user.PopulateAchievements()
	//user.PopulateStats()

	return &user, nil
}

func RegisterUser(username string, password string) error{
	stmt, err := _db.Prepare("INSERT INTO neb_users (username,password,rank) VALUES (?,?,1)")
	_, err = stmt.Exec(username, password)
	if err != nil {
		log.Println("Could not register new user :", err)
		return err
	}

	return nil
}

func (u *User)PopulateAchievements() error{
	rows, err := _db.Query("SELECT achievementid, progress, name, max FROM neb_users_achievements LEFT JOIN neb_achievements using (achievementid) WHERE neb_users_achievements.userid = ?", u.id)
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