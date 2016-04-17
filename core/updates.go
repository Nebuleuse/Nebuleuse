package core

import (
	"errors"
	"os"
	"time"
)

type Update struct {
	Build        *Build `json:"-"`
	BuildId      int
	Branch       string
	Size         int
	RollBack     bool
	SemVer       string
	Log          string
	Date         time.Time
	NextInBranch *Update `json:"-"`
	PrevInBranch *Update `json:"-"`
}
type Build struct {
	Id          int
	Commit      string
	Log         string
	Date        time.Time
	FileChanged string
	Obselete    bool
	Updates     map[string]*Update
}
type Branch struct {
	Name        string
	AccessRank  int
	ActiveBuild int
	Head        *Update
}

var updateBuilds map[int]*Build
var updateBranches map[string]Branch

func initUpdateSystem() error {
	updateBuilds = make(map[int]*Build)
	updateBranches = make(map[string]Branch)

	buildRows, err := Db.Query("SELECT id, commit, date, changelist, obselete FROM neb_updates_builds ORDER BY id")
	if err != nil {
		return err
	}
	defer buildRows.Close()
	for buildRows.Next() {
		var build Build
		err := buildRows.Scan(&build.Id, &build.Commit, &build.Date, &build.FileChanged, &build.Obselete)
		if err != nil {
			Warning.Println("Could not scan build from DB:", err)
			return err
		}
		build.Updates = make(map[string]*Update)
		updateBuilds[build.Id] = &build
	}

	branchRows, err := Db.Query("SELECT name, rank, activeBuild FROM neb_updates_branches")
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

	updateRows, err := Db.Query("SELECT build, branch, size, rollback, semver, log, date FROM neb_updates ORDER BY build")
	if err != nil {
		return err
	}
	defer updateRows.Close()
	for updateRows.Next() {
		var update Update

		err := updateRows.Scan(&update.BuildId, &update.Branch, &update.Size, &update.RollBack, &update.SemVer, &update.Log, &update.Date)
		if err != nil {
			return err
		}
		branch, ok := updateBranches[update.Branch]
		if !ok {
			Warning.Println("Skipped update #" + string(update.BuildId) + " because branch " + update.Branch + " does not exist")
			continue
		}
		if updateBuilds[update.BuildId] == nil {
			Warning.Println("Skipped update on branch " + update.Branch + " because buildid #" + string(update.BuildId) + " is incorrect")
			continue
		}
		update.Build = updateBuilds[update.BuildId]
		update.Build.Updates[update.Branch] = &update

		update.NextInBranch = branch.Head
		//If an update is rolled back, you can only go backward in update history
		if !update.RollBack {
			if branch.Head != nil {
				branch.Head.PrevInBranch = &update
			}
			branch.Head = &update
		}
	}

	//Setup Git if needed
	if Cfg.GetConfig("updateSystem") == "GitPatch" || Cfg.GetConfig("updateSystem") == "FullGit" {
		initGit()
	}

	return nil
}
func GetBranchList(rank int) []string {
	ret := []string{}
	for _, branch := range updateBranches {
		if branch.AccessRank&rank != 0 {
			ret = append(ret, branch.Name)
		}
	}
	return ret
}
func GetBranchHead(name string) (*Update, error) {
	branch, ok := updateBranches[name]
	if !ok || branch.Head == nil {
		return nil, errors.New("Branch not found or head missing:" + name)
	}
	return branch.Head, nil
}

//If branch doesn't exist, returns false
func CanUserAccessBranch(name string, rank int) bool {
	branch, ok := updateBranches[name]
	if !ok || branch.AccessRank&rank == 0 {
		return false
	}
	return true
}
func GetBranchUpdates(name string) ([]Update, error) {
	branch, ok := updateBranches[name]
	if !ok {
		return nil, errors.New("Branch not Found: " + name)
	}
	cur := branch.Head
	var ret []Update
	for cur != nil {
		ret = append(ret, *cur)
		cur = cur.NextInBranch
	}
	return ret, nil
}

func GetUpdateInfos(branchName string, buildId int) (*Update, error) {
	build := updateBuilds[buildId]
	if build == nil {
		return nil, errors.New("Build not found: " + string(buildId))
	}

	update, ok := build.Updates[branchName]
	if !ok {
		return nil, errors.New("No update info for build #" + string(buildId) + " on branch: " + branchName)
	}
	return update, nil
}

func SetActiveUpdate(branchName string, buildId int) error {
	branch, ok := updateBranches[branchName]
	if !ok {
		return errors.New("Branch not found: " + branchName)
	}
	build := updateBuilds[buildId]
	if build == nil {
		return errors.New("Build not found: " + string(buildId))
	}
	if build.Updates[branchName] == nil {
		return errors.New("Udate not found for build #" + string(buildId) + " on branch " + branchName + "")
	}

	branch.ActiveBuild = build.Id
	Db.Exec("UPDATE neb_updates_branches SET activeBuild = ? WHERE name = ?", build.Id, branch.Name)

	SignalGameUpdated(branch, *build.Updates[branchName])
	return nil
}

func UpdateGitCommitCache() error {
	return gitUpdateCommitCache()
}

func GetBranchActiveBuild(branchName string) (int, error) {
	branch, ok := updateBranches[branchName]
	if !ok {
		return 0, errors.New("Could not find branch: " + branchName)
	}
	return branch.ActiveBuild, nil
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

func SignalGameUpdated(branch Branch, update Update) {
	DispatchRank("system", "gameUpdate", update, branch.AccessRank)
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
