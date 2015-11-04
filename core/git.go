package core

import (
	"github.com/Nebuleuse/Nebuleuse/git"
)

var gitRepo *git.Repository

func initGit() error {
	var err error
	gitRepo, err = git.OpenRepository(Cfg["gitRepositoryPath"])
	if err != nil {
		Error.Println("Failed to open repository")
		return err
	}
	gitUpdateRepo()

	return nil
}

func gitUpdateRepo() {
	gitRepo.UpdateGitRepo(Cfg["productionBranch"])
}

func gitCreatePatch(commit string) {
	gitUpdateRepo()
	diff, _ := gitRepo.GetFilesChangedSinceUpdateRange(Cfg["productionBranch"], Cfg["currentCommit"], commit)

	Info.Println(diff)
}

func gitGetCommits(commit string) {
	gitUpdateRepo()
	list := gitRepo.CommitsBetween(gitRepo.GetCommit(""), gitRepo.GetCommit(""))
}
