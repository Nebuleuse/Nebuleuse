package core

import (
	"container/list"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"errors"

	"path/filepath"

	"encoding/json"

	"github.com/Nebuleuse/Nebuleuse/git"
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
	err := gitRepo.UpdateGitRepo(Cfg.GetConfig("productionBranch"))
	if err != nil {
		Warning.Println(err.Error())
	} else {
		Info.Println("Updated git repository")
	}
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
	var outDiffs []Diff
	found := make(map[string]int) //-1 if removed ,0 not present, else indicate position+1 in array

	for _, commit := range commits {
		for _, diff := range commit.Diff {
			if found[diff.Name] == 0 {
				outDiffs = append(outDiffs, diff)
				found[diff.Name] = len(outDiffs)
			} else if found[diff.Name] == -1 {
				continue
			} else {
				storedDiff := outDiffs[found[diff.Name]-1]
				if diff.IsCreated && !storedDiff.IsCreated {
					if !storedDiff.IsDeleted {
						outDiffs[found[diff.Name]-1].IsCreated = true
					} else { // It was created and deleted in between, no need to know about this file in our diffs
						found[diff.Name] = -1
					}
				} else if diff.IsDeleted {
					outDiffs[found[diff.Name]-1].IsCreated = false
					outDiffs[found[diff.Name]-1].IsDeleted = true
				}
			}
		}
	}
	var returnDiffs []Diff
	for _, diff := range outDiffs {
		if found[diff.Name] != -1 {
			returnDiffs = append(returnDiffs, diff)
		}
	}
	return returnDiffs
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

func gitGetCommitsBetweenCached(last, before string) ([]Commit, error) {
	startPos, endPos := -1, -1
	for i, commit := range commitCache {
		if commit.Id == last {
			startPos = i
		}
		if commit.Id == before {
			endPos = i
			break
		}
	}
	if startPos == -1 {
		return nil, errors.New("No commit found " + last)
	}
	if endPos == -1 {
		endPos = len(commitCache) - 1
	}
	return commitCache[startPos:endPos], nil
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

func gitGetAllCommitsCached() []Commit {
	ret := make([]Commit, len(commitCache))
	copy(ret, commitCache)
	return ret
}

func gitGetFirstCommit() (string, error) {
	if len(commitCache) == 0 {
		return "", errors.New("No commit recorder")
	}
	return commitCache[len(commitCache)-1].Id, nil
}

func gitCreatePatch(start, end string, buildTo, buildFrom int) (int64, error) {
	gitUpdateRepo()

	gitRepoLock.Lock()
	defer gitRepoLock.Unlock()

	commitsBetween, err := gitGetCommitsBetweenCached(start, end)
	if err != nil {
		return 0, err
	}
	diffs := gitGetDiffs(commitsBetween)

	size, err := _createPatch(start, strconv.Itoa(buildFrom)+"to"+strconv.Itoa(buildTo), diffs, true)
	_createPatch(end, strconv.Itoa(buildTo)+"to"+strconv.Itoa(buildFrom), diffs, false)
	gitRepo.Checkout("master")

	return size, err
}
func _createPatch(commit, filename string, diff []Diff, skipDeleted bool) (int64, error) {
	var deletedFiles []string
	updatesLocation := Cfg.GetSysConfig("UpdatesLocation")
	path := updatesLocation + "tmp/" + filename
	repoPath := Cfg.GetConfig("gitRepositoryPath")
	gitRepo.Checkout(commit)

	os.MkdirAll(path, 0764)
	for _, file := range diff {
		if skipDeleted && file.IsDeleted {
			deletedFiles = append(deletedFiles, file.Name)
			continue
		} else if !skipDeleted && file.IsCreated {
			deletedFiles = append(deletedFiles, file.Name)
			continue
		}

		outPath := filepath.Dir(file.Name)
		if outPath != "." {
			os.MkdirAll(path+"/"+outPath, 0764)
		}
		err := os.Link(repoPath+file.Name, path+"/"+file.Name)
		if err != nil {
			Error.Println("Could not Link file for patch building: " + err.Error())
		}
	}

	deletedManifest, _ := os.Create(path + "/" + "_deleted.json")

	if len(deletedFiles) > 0 {
		jsonbytes, _ := json.Marshal(deletedFiles)
		deletedManifest.Write(jsonbytes[:])
	} else {
		deletedManifest.WriteString("{}")
	}
	deletedManifest.Close()

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

func gitCheckoutCommit(commit string) {
	gitRepo.Checkout(commit)
}
