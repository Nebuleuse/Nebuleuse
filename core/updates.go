package core

import (
	"errors"
	"os"
	"time"
)

type Update struct {
	Build        *Build
	Branch       string
	Size         int
	RollBack     bool
	NextInBranch *Update
}
type Build struct {
	Id          int
	SemVer      string
	Commit      string
	Log         string
	Date        time.Time
	FileChanged string
	Obselete    bool
	Updates     map[string]*Update
}
type Branch struct {
	Name       string
	AccessRank int
	Head       *Update
}

var updateBuilds map[int]*Build
var updateBranches map[string]Branch

func initUpdateSystem() error {
	updateBuilds = make(map[int]*Build)
	updateBranches = make(map[string]Branch)

	buildRows, err := Db.Query("SELECT id, semver, commit, log, date, changelist, obselete FROM neb_updates_builds ORDER BY id")
	if err != nil {
		return err
	}
	defer buildRows.Close()
	for buildRows.Next() {
		var build Build
		err := buildRows.Scan(&build.Id, &build.SemVer, &build.Commit, &build.Log, &build.Date, &build.FileChanged, &build.Obselete)
		if err != nil {
			Warning.Println("Could not scan build from DB:", err)
			return err
		}
		build.Updates = make(map[string]*Update)
		updateBuilds[build.Id] = &build
	}

	branchRows, err := Db.Query("SELECT name, rank FROM neb_updates_branches")
	if err != nil {
		return err
	}
	defer branchRows.Close()
	for branchRows.Next() {
		var branch Branch
		err := branchRows.Scan(&branch.Name, &branch.AccessRank)
		if err != nil {
			Warning.Println("Could not scan branch from DB:", err)
			return err
		}
		updateBranches[branch.Name] = branch
	}

	updateRows, err := Db.Query("SELECT build, branch, size, rollback FROM neb_updates ORDER BY build")
	if err != nil {
		return err
	}
	defer updateRows.Close()
	for updateRows.Next() {
		var update Update
		var buildid int

		err := updateRows.Scan(&buildid, update.Branch, update.Size, update.RollBack)
		if err != nil {
			return err
		}

		update.Build = updateBuilds[buildid]
		update.NextInBranch = updateBranches[update.Branch].Head
		branch := updateBranches[update.Branch]
		branch.Head = &update
	}

	return nil
}
func GetUpdateInfos(version int) (Update, error) {
	var up Update

	/*err := Db.QueryRow("SELECT version, semVer, log, size, url, date, commit FROM neb_updates WHERE version = ?", version).Scan(&up.Version, &up.SemVer, &up.Log, &up.Size, &up.Url, &up.Date, &up.Commit)

	if err != nil && err == sql.ErrNoRows {
		return up, &NebuleuseError{NebError, "No update found"}
	} else if err != nil {
		return up, err
	}
	*/
	return up, nil
}

func GetUpdatesInfos(start int) ([]Update, error) {
	var updates []Update
	/*rows, err := Db.Query("SELECT version, semVer, log, size, url, date, commit FROM neb_updates WHERE version >= ?", start)
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
	*/
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
func GetCurrentVersion() int {
	return Cfg.GetConfigInt("gameVersion")
}

func GetLatestBuildCommit() (string, error) {
	if len(updateBuilds) == 0 {
		return "", errors.New("No build recorded, cannot acces latest build commit")
	}
	return updateBuilds[len(updateBuilds)-1].Commit, nil
}

func GetUpdateSystem() string {
	return Cfg.GetConfig("updateSystem")
}

func GetProductionBranch() string {
	return Cfg.GetConfig("productionBranch")
}

func GetGitCommitList() ([]Commit, error) {
	comm, err := GetLatestBuildCommit()
	if err != nil {
		return nil, err
	}

	return gitGetLatestCommitsCached(comm, 0)
}

func AddUpdate(info Update) error {
	/*if Cfg.GetConfig("updateSystem") == "GitPatch" {
		return createGitPatch(info)
	} else if Cfg.GetConfig("updateSystem") == "FullGit" {
		return addFullGitPatch(info)
	} else if Cfg.GetConfig("updateSystem") == "Manual" {

	}*/
	return nil
}

//Assumes update selected to create patch from is forward in history tree
/*func createGitPatch(info Update) error {
	gitCreatePatch(info.Commit)
	return nil
}
*/
type gitPatchPrepInfos struct {
	Diffs     []Diff
	TotalSize int64
}

func PrepareGitPatch(commit string) (gitPatchPrepInfos, error) {
	var res gitPatchPrepInfos
	comm, err := GetLatestBuildCommit()
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
		file, err := os.Open(Cfg.GetConfig("gitRepositoryPath") + c.Name)
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

func SignalGameUpdated() {
	Dispatch("system", "game", "game update released")
}

// Todo : no duplicate
func GetUpdateCount() int {
	var count int
	err := Db.QueryRow("SELECT COUNT(*) FROM neb_updates").Scan(&count)
	if err != nil {
		return -1
	}
	return count
}
