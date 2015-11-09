package core

import (
	"database/sql"
	"time"
)

type Update struct {
	Version int
	SemVer  string
	Log     string
	Size    int
	Url     string
	Date    time.Time
	Commit  string
}

func GetUpdateInfos(version int) (Update, error) {
	var up Update

	err := Db.QueryRow("SELECT version, semVer, log, size, url, date, commit FROM neb_updates WHERE version = ?", version).Scan(&up.Version, &up.SemVer, &up.Log, &up.Size, &up.Url, &up.Date, &up.Commit)

	if err != nil && err == sql.ErrNoRows {
		return up, &NebuleuseError{NebError, "No update found"}
	} else if err != nil {
		return up, err
	}

	return up, nil
}

func GetUpdatesInfos(start int) ([]Update, error) {
	var updates []Update
	rows, err := Db.Query("SELECT version, semVer, log, size, url, date, commit FROM neb_updates WHERE version >= ?", start)
	if err != nil {
		return updates, err
	}
	defer rows.Close()

	for rows.Next() {
		var update Update
		err := rows.Scan(&update.Version, &update.SemVer, &update.Log, &update.Size, &update.Url, &update.Date, &update.Commit)
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

func UpdateGitCommitCache() error {
	return gitUpdateCommitCache()
}

func GetCurrentCommit() string {
	return Cfg["currentCommit"]
}

func GetGitCommitList() ([]Commit, error) {
	return gitGetLatestCommitsCached(Cfg["currentCommit"], 5)
}

func AddUpdate(info Update) error {
	if Cfg["updateSystem"] == "GitPatch" {
		return createGitPatch(info)
	} else if Cfg["updateSystem"] == "FullGit" {
		return addFullGitPatch(info)
	} else if Cfg["updateSystem"] == "Manual" {

	}
	return nil
}

//Assumes update selected to create patch from is forward in history tree
func createGitPatch(info Update) error {
	gitCreatePatch(info.Commit)
	return nil
}

func addFullGitPatch(info Update) error {
	_, err := Db.Exec("INSERT INTO neb_updates VALUES(?,?,?,?,?,NOW(),?)", info.Version, info.SemVer, info.Log, info.Size, info.Url, info.Commit)
	if err != nil {
		Error.Println("Failed to insert update : ", info)
		return err
	}

	err = Cfg.SetConfig("currentCommit", info.Commit)
	if err != nil {
		return err
	}
	SignalGameUpdated()
	return nil
}

func SignalGameUpdated() {
	Dispatch("system", "game update released")
}

func GetUpdateCount() int {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM neb_updates").Scan(&count)
	if err != nil {
		return -1
	}
	return count
}
