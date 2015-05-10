package core

import (
	"strings"
)

type AchievementData struct {
	Id       int
	Name     string
	Max      int
	FullName string
	FullDesc string
	Icon     string
}

func GetAchievementsData() ([]AchievementData, error) {
	var achs []AchievementData
	rows, err := Db.Query("SELECT id, name, max, fullName, fullDesc, icon FROM neb_achievements")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ach AchievementData
		err := rows.Scan(&ach.Id, &ach.Name, &ach.Max, &ach.FullName, &ach.FullDesc, &ach.Icon)
		if err != nil {
			Error.Println("Could not get achievement table infos :", err)
			return nil, err
		}
		achs = append(achs, ach)
	}

	err = rows.Err()
	if err != nil {
		Error.Println("Could not get achievements tables infos :", err)
		return nil, err
	}
	return achs, nil
}
func GetAchievementData(id string) (AchievementData, error) {
	var ach AchievementData
	err := Db.QueryRow("SELECT id, name, max, fullName, fullDesc, icon FROM neb_achievements WHERE id=?", id).Scan(&ach.Id, &ach.Name, &ach.Max, &ach.FullName, &ach.FullDesc, &ach.Icon)
	if err != nil {
		return ach, err
	}
	return ach, nil
}

func SetAchievementData(id int, data AchievementData) error {
	stmt, err := Db.Prepare("UPDATE neb_achievements SET name=?, max=?, fullName=?, fullDesc=?, icon=? WHERE id=?")
	if err != nil {
		Error.Println("Could not prepare achievement data update :", err)
		return err
	}

	_, err = stmt.Exec(data.Name, data.Max, data.FullName, data.FullDesc, data.Icon, data.Id)
	if err != nil {
		Error.Println("Could not execute achievement data update : ", err, data)
		return err
	}
	return nil
}

func AddAchievementData(data AchievementData) (AchievementData, error) {
	stmt, err := Db.Prepare("INSERT INTO neb_achievements (name, max, fullName, fullDesc, icon) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		Error.Println("Could not prepare achievement data update :", err)
		return data, err
	}

	res, err := stmt.Exec(data.Name, data.Max, data.FullName, data.FullDesc, data.Icon)
	if err != nil {
		Error.Println("Could not execute achievement data update : ", err, data)
		return data, err
	}
	id, err := res.LastInsertId()
	data.Id = int(id)

	return data, nil
}

func DeleteAchievementData(id int) error {
	stmt, err := Db.Prepare("DELETE FROM neb_achievements WHERE id=?")
	if err != nil {
		Error.Println("Could not prepare achievement data deletion :", err)
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		Error.Println("Could not execute achievement data deletion : ", err)
		return err
	}

	return nil
}

func AddUserStat(u UserStat) error {
	var fields string
	err := Db.QueryRow("SELECT fields FROM neb_stats_table WHERE tableName = users").Scan(&fields)
	if err != nil {
		Error.Println("Could not get fields from neb_stats_table for users")
		return err
	}

	tFields := strings.Split(fields, ",")
	for _, field := range tFields {
		if field == u.Name {
			return &NebuleuseError{Code: NebError, Msg: "Field already exists in users"}
		}
	}
	err = Db.QueryRow("SELECT fields FROM neb_stats_table WHERE tableName = ?", u.Name).Scan(&fields)
	if err != nil {
		return &NebuleuseError{Code: NebError, Msg: "Table already exists"}
	}
	//TODO!
	return nil
}
