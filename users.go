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
	user.PopulateStats()

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
func (u *User)PopulateStats() error{
	rows, err := _db.Query("SELECT * FROM neb_users_stats WHERE userid = ?", u.id)
	if err != nil {
		log.Println("Could not get user stats :", err)
		return err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Println("Could not get columns:", err)
		return err
	}
	vals := make([]interface{}, len(cols))
	for i, _ := range cols {
		vals[i] = new(int64)
	}

	for rows.Next() {
		var st Stat
		err := rows.Scan(vals...)
		if err != nil {
			log.Println("Could not get user Stats :", err)
			return err
		}
		for i, _ := range cols {
			if i == 0 { //First column is userid
				continue
			}
			st.Name = cols[i]
			st.Value = *vals[i].(*int64)
			u.Stats = append(u.Stats, st)
		}
	}

	err = rows.Err()
	if err != nil {
		log.Println("Could not get user Stats :", err)
		return err
	}
	return nil
}