package core

import (
	"encoding/json"
	"errors"
	"time"
)

type Update struct {
	Build        *Build `json:"-"`
	BuildId      int
	Branch       string
	Size         int64
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
	Updates     map[string]*Update `json:"-"`
}
type Branch struct {
	Name        string
	AccessRank  int
	ActiveBuild int
	Head        *Update `json:"-"`
}

var updateBuilds []*Build
var updateBranches []*Branch

func initUpdateSystem() error {
	updateBuilds = make([]*Build, 0)
	updateBranches = make([]*Branch, 0)

	//	Load builds
	buildRows, err := Db.Query("SELECT id, commit, date, changelist, obselete, log FROM neb_updates_builds ORDER BY id")
	if err != nil {
		return err
	}
	defer buildRows.Close()
	for buildRows.Next() {
		var build Build
		err := buildRows.Scan(&build.Id, &build.Commit, &build.Date, &build.FileChanged, &build.Obselete, &build.Log)
		if err != nil {
			Warning.Println("Could not scan build from DB:", err)
			return err
		}
		build.Updates = make(map[string]*Update)
		updateBuilds = append(updateBuilds, &build)
	}

	//Load branches
	branchRows, err := Db.Query("SELECT name, rank, activeBuild FROM neb_updates_branches")
	if err != nil {
		return err
	}
	defer branchRows.Close()
	for branchRows.Next() {
		var branch Branch
		err := branchRows.Scan(&branch.Name, &branch.AccessRank, &branch.ActiveBuild)
		if err != nil {
			Warning.Println("Could not scan branch from DB:", err)
			return err
		}
		updateBranches = append(updateBranches, &branch)
	}
	//Load updates
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
		branch, err := GetBranch(update.Branch)
		if err != nil {
			Warning.Println("Skipped update #" + string(update.BuildId) + " because branch " + update.Branch + " does not exist")
			continue
		}
		build, err := GetBuild(update.BuildId)
		if err != nil {
			Warning.Println("Skipped update on branch " + update.Branch + " because buildid #" + string(update.BuildId) + " is incorrect")
			continue
		}
		update.Build = build
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
	if isGitUpdateSystem() {
		err := initGit()
		if err != nil {
			return err
		}
		err = gitUpdateCommitCache()
		if err != nil {
			return err
		}
	}

	return nil
}

//Inserts update in the structure and add to db
func insertUpdate(update *Update, build *Build, branch *Branch) error {
	stmt, err := Db.Prepare("INSERT INTO neb_updates(build, branch, size, rollback, semver, log) VALUES (?, ?, ?, 0, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(update.BuildId, update.Branch, update.Size, update.SemVer, update.Log)
	if err != nil {
		return err
	}
	update.NextInBranch = branch.Head
	if branch.Head != nil {
		branch.Head.PrevInBranch = update
	}
	branch.Head = update
	build.Updates[branch.Name] = update
	return nil
}

func GetBuild(id int) (*Build, error) {
	for _, build := range updateBuilds {
		if build.Id == id {
			return build, nil
		}
	}
	return nil, errors.New("Could not find build id:" + string(id))
}
func GetBranch(name string) (*Branch, error) {
	for _, branch := range updateBranches {
		if branch.Name == name {
			return branch, nil
		}
	}
	return nil, errors.New("Could not find branch named:" + name)
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
	branch, err := GetBranch(name)
	if err != nil || branch.Head == nil {
		return nil, err
	}
	return branch.Head, nil
}

//If branch doesn't exist, returns false
func CanUserAccessBranch(name string, rank int) bool {
	branch, err := GetBranch(name)
	if err != nil || branch.AccessRank&rank == 0 {
		return false
	}
	return true
}

//Only return updates that are active or are older than the active one
func GetBranchUpdates(name string) ([]Update, error) {
	branch, err := GetBranch(name)
	if err != nil {
		return nil, err
	}
	cur := branch.Head
	active := false
	var ret []Update
	for cur != nil {
		if branch.ActiveBuild == cur.BuildId {
			active = true
		}
		if active {
			ret = append(ret, *cur)
			cur = cur.NextInBranch
		}
	}
	return ret, nil
}

type branchUpdatesData struct {
	Name        string
	AccessRank  int
	ActiveBuild int
	Updates     []Update
}
type completeBranchUpdatesData struct {
	Branches []branchUpdatesData
	Builds   []Build
	Commits  []Commit
}

func GetCompleteUpdatesInfos() completeBranchUpdatesData {
	var res completeBranchUpdatesData
	for _, branch := range updateBranches {
		var branchData branchUpdatesData
		branchData.Name = branch.Name
		branchData.ActiveBuild = branch.ActiveBuild
		branchData.AccessRank = branch.AccessRank
		cur := branch.Head
		for cur != nil {
			branchData.Updates = append(branchData.Updates, *cur)
			cur = cur.NextInBranch
		}
		res.Branches = append(res.Branches, branchData)
	}
	for _, build := range updateBuilds {
		res.Builds = append(res.Builds, *build)
	}
	if isGitUpdateSystem() {
		commits, err := GetGitCommitList()
		if err != nil {
			Warning.Println("Could not get commit list:" + err.Error())
		} else {
			res.Commits = commits
		}
	}
	return res
}

func GetUpdateInfos(branchName string, buildId int) (*Update, error) {
	build, err := GetBuild(buildId)
	if err != nil {
		return nil, errors.New("Build not found: " + string(buildId))
	}

	update, ok := build.Updates[branchName]
	if !ok {
		return nil, errors.New("No update info for build #" + string(buildId) + " on branch: " + branchName)
	}
	return update, nil
}

func SetActiveUpdate(branchName string, buildId int) error {
	branch, err := GetBranch(branchName)
	if err != nil {
		return errors.New("Branch not found: " + branchName)
	}
	build, err := GetBuild(buildId)
	if err != nil {
		return errors.New("Build not found: " + string(buildId))
	}
	if build.Updates[branchName] == nil {
		return errors.New("Udate not found for build #" + string(buildId) + " on branch " + branchName + "")
	}

	_, err = Db.Exec("UPDATE neb_updates_branches SET activeBuild = ? WHERE name = ?", build.Id, branch.Name)
	if err != nil {
		return err
	}
	branch.ActiveBuild = build.Id

	SignalGameUpdated(*branch, *build.Updates[branchName])
	return nil
}

func UpdateGitCommitCache() error {
	return gitUpdateCommitCache()
}

func GetBranchActiveBuild(branchName string) (int, error) {
	branch, err := GetBranch(branchName)
	if err != nil {
		return 0, errors.New("Could not find branch: " + branchName)
	}
	return branch.ActiveBuild, nil
}

func GetLatestBuildCommit() (string, error) {
	if len(updateBuilds) == 0 {
		return "", errors.New("No build recorded, cannot acces latest build commit")
	}
	build := updateBuilds[len(updateBuilds)-1]
	return build.Commit, nil
}

func GetUpdateSystem() string {
	return Cfg.GetConfig("updateSystem")
}

func GetProductionBranch() string {
	return Cfg.GetConfig("productionBranch")
}
func GetUpdatesLocation() string {
	return Cfg.GetSysConfig("UpdatesLocation")
}
func isGitUpdateSystem() bool {
	if GetUpdateSystem() == "GitPatch" || GetUpdateSystem() == "FullGit" {
		return true
	}
	return false
}
func GetGitCommitList() ([]Commit, error) {
	comm, err := GetLatestBuildCommit()
	if err != nil { // No builds
		return gitGetAllCommitsCached(), nil
	}

	return gitGetLatestCommitsCached(comm, 0)
}

type gitBuildPrepInfos struct {
	Diffs     []Diff
	TotalSize int64
}

func PrepareGitBuild(commit string) (gitBuildPrepInfos, error) {
	var res gitBuildPrepInfos
	baseCommit, err := GetLatestBuildCommit()
	if err != nil { //No build recorded
		baseCommit, err = gitGetFirstCommit()

		if err != nil {
			return res, err
		}
	}
	commitsBetween, err := gitGetCommitsBetweenCached(commit, baseCommit)
	if err != nil {
		return res, err
	}

	Warning.Println("gitGetCommitsBetween c " + commit + " com " + baseCommit)

	diffs := gitGetDiffs(commitsBetween)
	var total int64
	gitLockRepo()
	gitCheckoutCommit(commit)
	for _, c := range diffs {
		if c.IsDeleted {
			continue
		}
		size, err := getFileSize(Cfg.GetConfig("gitRepositoryPath") + c.Name)
		if err != nil {
			Warning.Println("Could not read file size: " + err.Error())
		}
		total = total + size
	}
	gitCheckoutCommit("master")
	gitUnlockRepo()

	res.Diffs = diffs
	res.TotalSize = total

	return res, err
}
func CreateGitBuild(commit string, log string) error {
	prerInfos, err := PrepareGitBuild(commit)
	if err != nil {
		return err
	}
	msh, err := json.Marshal(prerInfos.Diffs)
	if err != nil {
		return err
	}
	changelist := string(msh)

	stmt, err := Db.Prepare("INSERT INTO neb_updates_builds(commit, log, changelist) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(commit, log, changelist)
	if err != nil {
		return err
	}

	var build Build
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	build.Id = int(id)
	build.Commit = commit
	build.Date = time.Now()
	build.Log = log
	build.FileChanged = changelist
	build.Updates = make(map[string]*Update)
	updateBuilds = append(updateBuilds, &build)

	return nil
}
func checkObseleteBuilds() error {
	foundFiles := make(map[string]int)
	for _, build := range updateBuilds {
		var diffs []Diff
		obselete := true
		err := json.Unmarshal([]byte(build.FileChanged), &diffs)
		if err != nil {
			Warning.Println("Could not unmarshal build changed files: " + err.Error())
			continue
		}
		for _, file := range diffs {
			_, ok := foundFiles[file.Name]
			if !ok {
				obselete = false
				foundFiles[file.Name] = build.Id
			}
		}
		if obselete {
			build.Obselete = true
			stmt, _ := Db.Prepare("UPDATE neb_updates_builds SET obselete=1 WHERE id = ?")
			stmt.Exec(build.Id)
		}
	}
	return nil
}
func CreateUpdate(build int, branch, semver, log string) error {
	var update Update
	buildObj, err := GetBuild(build)
	if err != nil {
		return err
	}
	branchObj, err := GetBranch(branch)
	if err != nil {
		return err
	}
	baseCommit := ""
	baseId := 0
	head := branchObj.Head
	if head == nil {
		//return errors.New("Branch " + branch + " has no head")
		baseCommit, err = gitGetFirstCommit()
		if err != nil {
			return err
		}
	} else {
		baseCommit = head.Build.Commit
		baseId = head.Build.Id
	}

	update.Branch = branch
	update.BuildId = build
	update.Build = buildObj
	update.Date = time.Now()
	update.Log = log
	update.SemVer = semver
	update.RollBack = false

	size, err := gitCreatePatch(buildObj.Commit, baseCommit, build, baseId)
	if err != nil {
		return err
	}
	update.Size = size

	err = insertUpdate(&update, buildObj, branchObj)

	return err
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
