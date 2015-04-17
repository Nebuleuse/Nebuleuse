package core

type AchievementsTable struct {
	Id       int
	Name     string
	Max      int
	FullName string
	FullDesc string
}

func GetAchievements() ([]AchievementsTable, error) {
	var achs []AchievementsTable
	rows, err := Db.Query("SELECT id, name, max, fullName, fullDesc FROM neb_achievements")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ach AchievementsTable
		err := rows.Scan(&ach.Id, &ach.Name, &ach.Max, &ach.FullName, &ach.FullDesc)
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
