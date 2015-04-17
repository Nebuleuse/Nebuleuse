package core

import (
	"database/sql"
)

type Update struct {
	Version int
	Log     string
	Size    int
	Date    int
	Commit  string
}

func GetUpdateInfos(version int) (Update, error) {
	var up Update

	err := Db.QueryRow("SELECT version, log, size, date, commit FROM neb_updates WHERE version = ?", version).Scan(&up.Version, &up.Log, &up.Size, &up.Date, &up.Commit)

	if err != nil && err == sql.ErrNoRows {
		return up, &NebuleuseError{NebError, "No update found"}
	} else if err != nil {
		return up, err
	}

	return up, nil
}
func GetUpdatesInfos(start int, end int) ([]Update, error) {
	var updates []Update
	rows, err := Db.Query("SELECT version, log, size, date, commit FROM neb_updates WHERE version >= ? AND version <= ?", start, end)
	if err != nil {
		return updates, err
	}
	defer rows.Close()

	for rows.Next() {
		var update Update
		err := rows.Scan(&update.Version, &update.Log, &update.Size, &update.Date, &update.Commit)
		if err != nil {
			Warning.Println("Could not get update infos :", err)
			return updates, err
		}
		updates = append(updates, update)
	}

	err = rows.Err()
	if err != nil {
		Warning.Println("Could not get update infos :", err)
		return updates, err
	}

	return updates, nil
}
func SetActiveUpdate(version int) {
	//Todo
}
func AddUpdate(info Update) {
	Db.Query("INSERT INTO neb_updates VALUES(?,?,?,?,?)", info.Version, info.Log, info.Size, info.Date, info.Commit)
}
func PublishNewUpdate(info Update) {
	//WIP
	if Cfg["updateSystem"] == "GitPatch" {
		GitPreparePatch()
	}
}
func GetUpdateCount() int {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM neb_updates").Scan(&count)
	if err != nil {
		return -1
	}
	return count
}
