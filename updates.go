package main

import (
	"Nebuleuse/gitUpdater"
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

	err := _db.QueryRow("SELECT version, log, size, date, commit FROM neb_updates WHERE version = ?", version).Scan(&up.version, &up.log, &up.size, &up.date, &up.commit)

	if err != nil && err == sql.ErrNoRows {
		return up, &NebuleuseError{NebError, "No update found"}
	} else if err != nil {
		return up, err
	}

	return up, nil
}
func GetUpdatesInfos(start int, end int) ([]Update, error) {
	var updates []Update
	rows, err := _db.Query("SELECT version, log, size, date FROM neb_updates WHERE version >= ? AND version <= ?", start, end)
	if err != nil {
		return updates, err
	}
	defer rows.Close()

	for rows.Next() {
		var update Update
		err := rows.Scan(&update.version, &update.log, &update.size, &update.date)
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
	_db.Query("INSERT INTO neb_updates VALUES(?,?,?,?,?)", info.Version, info.Log, info.Size, info.Date, info.Commit)
}
func PublishNewUpdate(info Update) {
	//WIP
	var info Update
	info.Commit = commit
	if _cfg["updateSystem"] == "GitPatch" {
		gitUpdater.PreparePatch(commit)
	}
}
