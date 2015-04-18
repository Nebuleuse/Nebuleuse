package core

import (
	"database/sql"
)

//Achievements
type Achievement struct {
	Id       int
	Name     string
	Progress uint
	Value    uint
	Icon     string
}

func (a *Achievement) isComplete() bool {
	return a.Progress == a.Value
}

// Stats
type UserStat struct {
	Name  string
	Value int64
}
type KeyValue struct {
	Name  string
	Value string
}
type ComplexStat struct {
	Name   string
	Values []KeyValue
}

// User
type User struct {
	Id           int
	Username     string
	SessionId    string
	Rank         int
	Avatar       string
	Achievements []Achievement
	Stats        []UserStat
}

const (
	UserMaskBase = 1 << iota
	UserMaskOnlyId
	UserMaskAchievements
	UserMaskStats
	UserMaskAll = UserMaskStats | UserMaskAchievements
)

func GetUserBySession(SessionId string, BitMask int) (*User, error) {
	var user User
	user.SessionId = SessionId

	var id int
	err := Db.QueryRow("SELECT userid FROM neb_sessions WHERE sessionid = ?", SessionId).Scan(&id)
	if err != nil && err == sql.ErrNoRows {
		return nil, &NebuleuseError{NebErrorDisconnected, "No session found"}
	} else if err != nil {
		return nil, err
	}
	user.Id = id

	user.FetchUserInfos(BitMask)

	go user.Heartbeat()

	return &user, nil
}

func (u *User) FetchUserInfos(Bitmask int) error {
	if Bitmask&UserMaskOnlyId != 0 {
		return nil
	}

	err := Db.QueryRow("SELECT username, rank, avatars FROM neb_users WHERE id = ?", u.Id).Scan(&u.Username, &u.Rank, &u.Avatar)
	if err != nil && err == sql.ErrNoRows {
		return &NebuleuseError{NebErrorDisconnected, "No user found"}
	} else if err != nil {
		return err
	}
	if u.Avatar == "" {
		u.Avatar = Cfg["defaultAvatar"]
	}

	if Bitmask&UserMaskAchievements != 0 {
		u.PopulateAchievements()
	}
	if Bitmask&UserMaskStats != 0 {
		u.PopulateStats()
	}

	return nil
}

func RegisterUser(username string, password string, hash string) error {
	stmt, err := Db.Prepare("INSERT INTO neb_users (username,password,rank,hash) VALUES (?,?,1,?)")
	_, err = stmt.Exec(username, password, hash)
	if err != nil {
		Warning.Println("Could not register new user :", err)
		return err
	}

	return nil
}

func (u *User) PopulateAchievements() error {
	rows, err := Db.Query("SELECT id, name, max, icon, progress FROM neb_achievements AS ach LEFT JOIN neb_users_achievements AS usr ON ach.id = usr.achievementid AND usr.userid = ?", u.Id)
	if err != nil {
		Warning.Println("Could not get user achievements :", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ach Achievement
		var progress sql.NullInt64
		progress.Int64 = 0
		err := rows.Scan(&ach.Id, &ach.Progress, &ach.Name, &ach.Icon, &progress)
		if err != nil {
			Warning.Println("Could not get user achievements :", err)
			return err
		}
		ach.Value = uint(progress.Int64)
		u.Achievements = append(u.Achievements, ach)
	}

	err = rows.Err()
	if err != nil {
		Warning.Println("Could not get user achievements :", err)
		return err
	}
	return nil
}
func (u *User) PopulateStats() error {
	rows, err := Db.Query("SELECT name, value FROM neb_users_stats WHERE userid = ?", u.Id)
	if err != nil {
		Warning.Println("Could not get user stats :", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var st UserStat
		err := rows.Scan(&st.Name, &st.Value)
		if err != nil {
			Warning.Println("Could not get user Stats :", err)
			return err
		}
		u.Stats = append(u.Stats, st)
	}

	err = rows.Err()
	if err != nil {
		Warning.Println("Could not get user Stats :", err)
		return err
	}
	return nil
}
func (u *User) Heartbeat() {
	if sess, ok := connectedUsers[u.Id]; ok {
		sess.Heartbeat()
	}
}
func (u *User) Disconnect() {
	stmt, err := Db.Prepare("DELETE FROM neb_sessions WHERE userid = ?")
	_, err = stmt.Exec(u.Id)
	if err != nil {
		Warning.Println("Could not delete user session :", err)
	}
}
func (u *User) SetAchievementProgress(aid int, value int) error {
	stmt, err := Db.Prepare("UPDATE neb_users_achievements SET progress= ? WHERE userid = ? AND achievementid = ? LIMIT 1")
	if err != nil {
		Warning.Println("Could not create statement : ", err)
		return err
	}

	res, err := stmt.Exec(value, u.Id, aid)
	if err != nil {
		Warning.Println("Could not update achievement :", err)
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		Warning.Println("Could not get update rowcount :", err)
		return err
	}
	if rowCnt == 0 {
		Warning.Println("Tried to update achievementid : ", aid, " but no rows affected")
		return &NebuleuseError{NebError, "No rows affected by operation"}
	}

	return nil
}
func (u *User) SetStats(stats []UserStat) error {
	stmt, err := Db.Prepare("UPDATE neb_users_stats SET value = ? WHERE userid = ? AND name = ? LIMIT 1")
	if err != nil {
		Warning.Println("Could not create statement : ", err)
		return err
	}

	for _, stat := range stats {
		if stat.Name == "userid" {
			Warning.Println("Could not update user stats, userid present in stat list")
			return &NebuleuseError{NebErrorLogin, "Could not update user stats, userid present in stat list"}
		}
		_, err := stmt.Exec(stat.Value, u.Id, stat.Name)
		if err != nil {
			Warning.Println("Could not update user stats : ", err)
			return err
		}
	}
	return nil
}
func (u *User) SetComplexStats(stats []ComplexStat) error {
	var count = 0
	for _, stat := range stats {
		tableInfo, err := GetComplexStatsTableInfos(stat.Name)
		if err != nil {
			Warning.Println("Could not get fields for table : ", stat.Name, err)
			continue
		}
		//Prepare SQL request
		cmd := "INSERT INTO neb_users_stats_"
		cmd += stat.Name
		cmd += " VALUES ("
		for i := 0; i < len(tableInfo.Fields); i++ {
			if i == len(tableInfo.Fields)-1 {
				cmd += "?"
			} else {
				cmd += "?,"
			}
		}
		cmd += ")"

		stmt, err := Db.Prepare(cmd)
		if err != nil {
			Warning.Println("Could not prepare statement : ", err)
			continue
		}
		//Sort values so they match the table definition
		var sortedValues []interface{}
		for _, field := range tableInfo.Fields {
			if field == "userid" {
				sortedValues = append(sortedValues, u.Id)
				continue
			}
			for _, value := range stat.Values {
				if value.Name != field {
					continue
				}
				sortedValues = append(sortedValues, value.Value)
			}
		}

		if len(sortedValues) == 0 {
			Warning.Println("No correct values to insert into stat table: ", stat.Name)
			continue
		}
		_, err = stmt.Exec(sortedValues...)
		if err != nil {
			Warning.Println("Could not insert data into stat table : ", err)
			continue
		}
		count++

		if tableInfo.AutoCount { // Do we need to update the player stat associated ?
			var stt []UserStat
			var st UserStat
			for _, s := range u.Stats {
				if s.Name == stat.Name {
					st = s
					break
				}
			}
			st.Value = st.Value + 1
			stt = append(stt, st)
			u.SetStats(stt)
		}
	}
	if count < len(stats) {
		return &NebuleuseError{NebErrorPartialFail, "Inserted " + string(count) + " stats out of " + string(len(stats))}
	}
	return nil
}
func GetUserCount() int {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM neb_users").Scan(&count)
	if err != nil {
		return -1
	}
	return count
}

func GetUsersInfos(start, count, mask int) ([]*User, error) {
	var Users []*User
	rows, err := Db.Query("SELECT id FROM neb_users LIMIT ?, ?", start, count)
	if err != nil {
		Warning.Println("Could not fetch users infos : ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var usr User
		err := rows.Scan(&usr.Id)
		if err != nil {
			Warning.Println("Could not scan id users infos : ", err)
			return nil, err
		}
		usr.FetchUserInfos(mask)
		Users = append(Users, &usr)
	}

	err = rows.Err()
	if err != nil {
		Warning.Println("Could not get Users infos (after loop) :", err)
		return nil, err
	}

	return Users, nil
}
