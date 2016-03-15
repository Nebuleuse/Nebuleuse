package core

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
)

//Achievements
type Achievement struct {
	Id       int
	Name     string
	Progress uint
	Max      uint
	Icon     string
}

func (a *Achievement) isComplete() bool {
	return a.Progress == a.Max
}

// User
type User struct {
	Id           int
	Username     string
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

const (
    UserRankBanned = 1 << iota
    UserRankNormal
    UserRankDev
    UserRankAdmin
)

func GetUserBySession(SessionId string, BitMask int) (*User, error) {
	var user User

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

func RegisterUser(username string, password string, rank int) error {
	c := sha512.Size
	bhash := make([]byte, c)
	_, err := rand.Read(bhash)
	if err != nil {
		Error.Println("Error generating crytpo hash:", err)
		return err
	}
	hash := base64.URLEncoding.EncodeToString(bhash)
	hashedPassword := HashPassword(password, string(hash))

	stmt, err := Db.Prepare("INSERT INTO neb_users (username,password,rank,hash) VALUES (?,?,?,?)")
	_, err = stmt.Exec(username, hashedPassword, rank, string(hash))
	if err != nil {
		Error.Println("Could not register new user :", err)
		return err
	}

	return nil
}

func (u *User) PopulateAchievements() error {
	rows, err := Db.Query("SELECT id, name, max, icon, progress FROM neb_achievements AS ach LEFT JOIN neb_users_achievements AS usr ON ach.id = usr.achievementid AND usr.userid = ?", u.Id)
	if err != nil {
		Error.Println("Could not get user achievements :", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var ach Achievement
		var progress sql.NullInt64
		progress.Int64 = 0
		err := rows.Scan(&ach.Id, &ach.Name, &ach.Max, &ach.Icon, &progress)
		if err != nil {
			Error.Println("Could not get user achievements :", err)
			return err
		}
		ach.Progress = uint(progress.Int64)
		u.Achievements = append(u.Achievements, ach)
	}

	err = rows.Err()
	if err != nil {
		Error.Println("Could not get user achievements :", err)
		return err
	}
	return nil
}

// We load complex stats informations stored in neb_stats_tables
// to make a list of stats the user has. Entry for users is the users' fields
// and entries with AutoCount true are additional fields updated when a complex stat is added
func (u *User) PopulateStats() error {
	Fields, err := GetUserStatsFields()
	if err != nil {
		Error.Println("Could not get users stats fields :", err)
		return err
	}
	StatFields := make(map[string]int64)
	for _, field := range Fields {
		StatFields[field] = 0
	}

	rows, err := Db.Query("SELECT name, value FROM neb_users_stats WHERE userid = ?", u.Id)
	if err != nil {
		Error.Println("Could not get user stats :", err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var value int64
		err := rows.Scan(&name, &value)
		if err != nil {
			Error.Println("Could not get user Stats :", err)
			return err
		}
		StatFields[name] = value
	}

	err = rows.Err()
	if err != nil {
		Error.Println("Could not get user Stats :", err)
		return err
	}

	for field, value := range StatFields {
		var st UserStat
		st.Name = field
		st.Value = value
		u.Stats = append(u.Stats, st)
	}
	return nil
}
func (u *User) Heartbeat() {
	if sess, ok := connectedUsers[u.Id]; ok {
		sess.Heartbeat()
	}
}
func (u *User) Disconnect() {
	DisconnectUser(u.Id)
}
func (u *User) SetAchievementProgress(aid int, value int) error {
	stmt, err := Db.Prepare("UPDATE neb_users_achievements SET progress= ? WHERE userid = ? AND achievementid = ? LIMIT 1")
	if err != nil {
		Error.Println("Could not create statement : ", err)
		return err
	}

	res, err := stmt.Exec(value, u.Id, aid)
	if err != nil {
		Error.Println("Could not update achievement :", err)
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		Error.Println("Could not get update rowcount :", err)
		return err
	}
	if rowCnt == 0 {
		res, err = Db.Exec("INSERT INTO neb_users_achievements (userid, achievementid, progress) VALUES (?, ?, ?)", u.Id, aid, value)
		if err != nil {
			Error.Println("Could not insert achievement for user:", err)
			return err
		}
	}

	return nil
}
func (u *User) SetStats(stats []UserStat) error {
	stmt, err := Db.Prepare("UPDATE neb_users_stats SET value = ? WHERE userid = ? AND name = ? LIMIT 1")
	if err != nil {
		Error.Println("Could not create statement : ", err)
		return err
	}

	for _, stat := range stats {
		if stat.Name != "userid" {
			res, err := stmt.Exec(stat.Value, u.Id, stat.Name)
			if err != nil {
				Error.Println("Could not update user stats : ", err)
				return err
			}
			rows, err := res.RowsAffected()
			if err != nil {
				Error.Println("Error getting rows affected : ", err)
				return err
			}
			if rows == 0 {
				_, err = Db.Exec("INSERT INTO neb_users_stats (userid, name, value) VALUES (?,?,?)", u.Id, stat.Name, stat.Value)
				if err != nil {
					Error.Println("Could not update user stats : ", err)
					return err
				}
			}
		}
	}
	return nil
}

func (u *User) SetComplexStats(stats []ComplexStat) error {
	var count = 0
	for _, stat := range stats {
		tableInfo, err := GetComplexStatsTableInfos(stat.Name)
		if err != nil {
			Error.Println("Could not get fields for table : ", stat.Name, err)
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
			Error.Println("Could not prepare statement : ", err)
			continue
		}
		//Sort values so they match the table definition
		var sortedValues []interface{}
		for _, field := range tableInfo.Fields {
			if field.Name == "userid" {
				sortedValues = append(sortedValues, u.Id)
				continue
			}
			for _, value := range stat.Values {
				if value.Name != field.Name {
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
			Error.Println("Could not insert data into stat table : ", err)
			continue
		}
		count++

		if tableInfo.AutoCount { // Do we need to update the player stat associated ?
			if len(u.Stats) == 0 {
				u.PopulateStats()
			}

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
		Error.Println("Could not fetch users infos : ", err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var usr User
		err := rows.Scan(&usr.Id)
		if err != nil {
			Error.Println("Could not scan id users infos : ", err)
			return nil, err
		}
		usr.FetchUserInfos(mask)
		Users = append(Users, &usr)
	}

	err = rows.Err()
	if err != nil {
		Error.Println("Could not get Users infos (after loop) :", err)
		return nil, err
	}

	return Users, nil
}
