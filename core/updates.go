package core

import (
	"database/sql"
	"os"
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

func SetActiveUpdate(version int) error {
	err := Cfg.SetConfig("gameVersion", string(version))
	if err != nil {
		return err
	}
	SignalGameUpdated()
	return nil
}

func UpdateGitCommitCache() error {
	return gitUpdateCommitCache()
}

func GetCurrentCommit() (string, error) {
	info, err := GetUpdateInfos(GetConfigInt("gameVersion"))
	if err != nil {
		return "", err
	}
	return info.Commit, nil
}

func GetUpdateSystem() string {
	return Cfg["updateSystem"]
}

func GetProductionBranch() string {
	return Cfg["productionBranch"]
}

func GetGitCommitList() ([]Commit, error) {
	comm, err := GetCurrentCommit()
	if err != nil {
		return nil, err
	}

	return gitGetLatestCommitsCached(comm, 0)
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

type gitPatchPrepInfos struct {
	Diffs     []Diff
	TotalSize int64
}

func PrepareGitPatch(commit string) (gitPatchPrepInfos, error) {
	var res gitPatchPrepInfos
	comm, err := GetCurrentCommit()
	if err != nil {
		return res, err
	}
	com, err := gitGetCommitsBetween(commit, comm)
	if err != nil {
		return res, err
	}
	diffs := gitGetDiffs(com)
	var total int64
	for _, c := range diffs {
		file, err := os.Open(Cfg["gitRepositoryPath"] + c.Name)
		if err != nil {
			return res, err
		}
		stat, err := file.Stat()
		if err != nil {
			return res, err
		}
		total = total + stat.Size()
		file.Close()
	}

	res.Diffs = diffs
	res.TotalSize = total

	return res, err
}
func addFullGitPatch(info Update) error {
	_, err := Db.Exec("INSERT INTO neb_updates VALUES(?,?,?,?,?,NOW(),?)", info.Version, info.SemVer, info.Log, info.Size, info.Url, info.Commit)
	if err != nil {
		Error.Println("Failed to insert update : ", info)
		return err
	}

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
