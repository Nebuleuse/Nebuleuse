package core

import (
	"container/list"
	"github.com/Nebuleuse/Nebuleuse/git"
	"os"
	"os/exec"
	"strconv"
	"sync"
)

var gitRepo *git.Repository
var commitCache []Commit
var gitRepoLock sync.Mutex

func initGit() error {
	var err error
	gitRepo, err = git.OpenRepository(Cfg.GetConfig("gitRepositoryPath"))
	if err != nil {
		Error.Println("Failed to open repository")
		return err
	}

	return nil
}

func gitUpdateCommitCache() error {
	gitUpdateRepo()
	Info.Println("Updating commit cache")

	if len(commitCache) == 0 {
		cache, err := gitGetCommits("")
		if err != nil {
			Warning.Println("Could not get complete commits list")
			return err
		}
		Info.Printf("Recorded %d commits from HEAD", len(cache))
		commitCache = cache
	} else {
		newCommits, err := gitGetCommits(commitCache[0].Id)
		if err != nil {
			Warning.Println("Could not recent commits list")
			return err
		}
		if len(newCommits) != 0 {
			commitCache = append(newCommits, commitCache...)
			Info.Printf("Recorded %d new commits", len(newCommits))
		}
	}
	return nil
}

func gitUpdateRepo() {
	gitRepo.UpdateGitRepo(Cfg.GetConfig("productionBranch"))
	Info.Println("Updated git repository")
}
func gitLockRepo() {
	gitRepoLock.Lock()
}
func gitUnlockRepo() {
	gitRepoLock.Unlock()
}

type Commit struct {
	Id        string
	Message   string
	Committer string
	Date      string

	TotalAddition, TotalDeletion int
	Diff                         []Diff
}
type Diff struct {
	Name               string
	Index              int
	Addition, Deletion int
	Type               int
	IsCreated          bool
	IsDeleted          bool
	IsBin              bool
}

func gitGetDiffs(commits []Commit) []Diff {
	var diffs []Diff
	found := make(map[string]int) //-1 if removed ,0 not present, else indicate position+1 in array

	for _, commit := range commits {
		for _, diff := range commit.Diff {
			if found[diff.Name] == 0 {
				diffs = append(diffs, diff)
				found[diff.Name] = len(diffs)
			} else if found[diff.Name] == -1 {
				continue
			} else {
				storedDiff := diffs[found[diff.Name]-1]
				if diff.IsCreated && !storedDiff.IsCreated {
					if !storedDiff.IsDeleted {
						diffs[found[diff.Name]-1].IsCreated = true
					} else { // It was created and deleted in between
						//Remove file from list
						pos := found[diff.Name] - 1
						diffs = append(diffs[:pos], diffs[pos+1:]...)
						found[diff.Name] = -1
					}
				} else if diff.IsDeleted {
					diffs[found[diff.Name]-1].IsCreated = false
					diffs[found[diff.Name]-1].IsDeleted = true
				}
			}
		}
	}
	return diffs
}

func gitParseCommitList(list *list.List) []Commit {
	var Commits []Commit
	i := 0
	for e := list.Front(); e != nil; e = e.Next() {
		var commit Commit
		gitCommit := e.Value.(*git.Commit)
		commit.Id = gitCommit.Id.String()
		commit.Message = gitCommit.Message()
		commit.Committer = gitCommit.Committer.Name + " " + gitCommit.Committer.Email
		commit.Date = gitCommit.Committer.When.String()

		Diffs, _ := gitRepo.GetDiffCommit(commit.Id)
		commit.TotalAddition = Diffs.TotalAddition
		commit.TotalDeletion = Diffs.TotalDeletion

		for _, file := range Diffs.Files {
			var diff Diff
			diff.Addition = file.Addition
			diff.Deletion = file.Deletion
			diff.Index = file.Index
			diff.IsBin = file.IsBin
			diff.IsCreated = file.IsCreated
			diff.IsDeleted = file.IsDeleted
			diff.Name = file.Name
			diff.Type = file.Type
			commit.Diff = append(commit.Diff, diff)
		}
		Commits = append(Commits, commit)
		i++
	}
	return Commits
}

func gitGetCommits(commit string) ([]Commit, error) {
	//Get commits between HEAD and the provided commit
	headCommit, err := gitRepo.GetCommitOfBranch(Cfg.GetConfig("productionBranch"))
	if err != nil {
		Error.Println("Could not get Head commit of production branch")
		return nil, err
	}
	var list *list.List
	if commit == "" { // Get list of all commits
		list, err = headCommit.CommitsBefore()
		if err != nil {
			Error.Println("Could not get commits of production branch")
			return nil, err
		}
	} else { // Get list of commits between HEAD and specified commit
		list, err = headCommit.CommitsBeforeUntil(commit)
		if err != nil {
			Error.Println("Could not get commits of production branch")
			return nil, err
		}
	}

	return gitParseCommitList(list), err
}

func gitGetCommitsBetween(last, before string) ([]Commit, error) {
	lastCom, err := gitRepo.GetCommit(last)
	if err != nil {
		return nil, err
	}
	beforeCom, err := gitRepo.GetCommit(before)
	if err != nil {
		return nil, err
	}
	list, err := gitRepo.CommitsBetween(lastCom, beforeCom)
	if err != nil {
		return nil, err
	}
	return gitParseCommitList(list), err
}

//Get latest commits with no duplicates
func gitGetLatestCommitsCached(commit string, after int) ([]Commit, error) {
	if len(commitCache) == 0 {
		return commitCache, nil
	}
	found := false
	endPos := 0
	for _, com := range commitCache {
		if found && after == 0 {
			break
		} else if found {
			after--
		} else if com.Id == commit {
			found = true
		}
		endPos++
	}
	if !found {
		Warning.Println("Could not find commit during lookup:" + commit)
		//Means the commit isn't in this branch. Get full list instead
		endPos++
	}
	ret := make([]Commit, endPos-1)
	copy(ret, commitCache[0:endPos-1])
	return ret, nil
}

func gitCreatePatch(start, end string, buildTo, buildFrom int) (int64, error) {
	gitUpdateRepo()

	gitRepoLock.Lock()
	defer gitRepoLock.Unlock()

	diff, _ := gitRepo.GetFilesChangedSinceUpdateRange(end, start)
	size, err := _createPatch(start, strconv.Itoa(buildFrom)+"to"+strconv.Itoa(buildTo), diff, true)
	_createPatch(end, strconv.Itoa(buildTo)+"to"+strconv.Itoa(buildFrom), diff, false)
	gitRepo.Checkout("master")

	return size, err
}
func _createPatch(commit, filename string, diff *git.Diff, skipDeleted bool) (int64, error) {
	updatesLocation := Cfg.GetSysConfig("UpdatesLocation")
	path := updatesLocation + "tmp/" + filename
	repoPath := Cfg.GetConfig("gitRepositoryPath")
	gitRepo.Checkout(commit)

	os.MkdirAll(path, 0764)
	for _, file := range diff.Files {
		if skipDeleted && file.IsDeleted {
			continue
		} else if !skipDeleted && file.IsCreated {
			continue
		}
		err := os.Link(repoPath+file.Name, path+"/"+file.Name)
		if err != nil {
			Error.Println("Could not Link file for patch building: " + err.Error())
		}
	}

	cmdTar := exec.Command("tar", "-c", filename, "-f", filename+".tar")
	cmdTar.Stdout = os.Stdout
	cmdTar.Stderr = os.Stderr
	cmdTar.Dir = "./updates/tmp/"
	cmdTar.Run()
	cmdXz := exec.Command("xz", "-z", filename+".tar")
	cmdXz.Stdout = os.Stdout
	cmdXz.Stderr = os.Stderr
	cmdXz.Dir = "./updates/tmp/"
	cmdXz.Run()
	os.RemoveAll(path)
	os.Rename(path+".tar.xz", updatesLocation+filename+".tar.xz")

	return getFileSize(updatesLocation + filename + ".tar.xz")
}
