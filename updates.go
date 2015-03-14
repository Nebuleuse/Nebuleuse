package main

type Update struct {
	version int
	log     string
	size    int
	date    int
}

func GetUpdateInfos(version int) (Update, error) {
	var up Update
	err := _db.QueryRow("SELECT version, log, size, date FROM neb_updates WHERE version = ?", version).Scan(&up.version, &up.log, &up.size, &up.date)
	if err != nil && err == sql.ErrNoRows {
		return nil, &NebuleuseError{NebError, "No update found"}
	} else if err != nil {
		return nil, err
	}
	return up, nil
}
func GetUpdatesInfos(start int, end int) ([]Update, error) {
	var updates []Update
	err := _db.Query("SELECT version, log, size, date FROM neb_updates WHERE version >= ? AND version <= ?", start, end)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var update Update
		err := rows.Scan(&update.version, &update.log, &update.size, &update.date)
		if err != nil {
			Warning.Println("Could not get update infos :", err)
			return nil, err
		}
		updates = append(updates, update)
	}

	err = rows.Err()
	if err != nil {
		Warning.Println("Could not get update infos :", err)
		return err
	}

	return updates, nil
}
