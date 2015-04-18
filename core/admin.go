package core

type AchievementsTable struct {
	Id       int
	Name     string
	Max      int
	FullName string
	FullDesc string
	Icon     string
}

func GetAchievementsData() ([]AchievementsTable, error) {
	var achs []AchievementsTable
	rows, err := Db.Query("SELECT id, name, max, fullName, fullDesc, icon FROM neb_achievements")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ach AchievementsTable
		err := rows.Scan(&ach.Id, &ach.Name, &ach.Max, &ach.FullName, &ach.FullDesc, &ach.Icon)
		if err != nil {
			Warning.Println("Could not get achievement table infos :", err)
			return nil, err
		}
		achs = append(achs, ach)
	}

	err = rows.Err()
	if err != nil {
		Warning.Println("Could not get achievements tables infos :", err)
		return nil, err
	}
	return achs, nil
}

func SetAchievementData(id int, data AchievementsTable) error {
	stmt, err := Db.Prepare("UPDATE neb_achievements SET name=?, max=?, fullName=?, fullDesc=?, icon=? WHERE id=?")
	if err != nil {
		Warning.Println("Could not prepare achievement data update :", err)
		return err
	}

	_, err = stmt.Exec(data.Name, data.Max, data.FullName, data.FullDesc, data.Icon, data.Id)
	if err != nil {
		Warning.Println("Could not execute achievement data update : ", err, data)
		return err
	}
	return nil
}

func AddAchievementData(data AchievementsTable) (AchievementsTable, error) {
	stmt, err := Db.Prepare("INSERT INTO neb_achievements (name, max, fullName, fullDesc, icon) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		Warning.Println("Could not prepare achievement data update :", err)
		return data, err
	}

	res, err := stmt.Exec(data.Name, data.Max, data.FullName, data.FullDesc, data.Icon)
	if err != nil {
		Warning.Println("Could not execute achievement data update : ", err, data)
		return data, err
	}
	id, err := res.LastInsertId()
	data.Id = int(id)

	return data, nil
}

func DeleteAchievementData(id int) error {
	stmt, err := Db.Prepare("DELETE FROM neb_achievements WHERE id=?")
	if err != nil {
		Warning.Println("Could not prepare achievement data deletion :", err)
		return err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		Warning.Println("Could not execute achievement data deletion : ", err)
		return err
	}

	return nil
}
