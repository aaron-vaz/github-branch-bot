package service

import (
	"fmt"
	"log"
	"sync"

	"github.com/aaron-vaz/github-branch-bot/pkg/config"
	"github.com/aaron-vaz/github-branch-bot/pkg/github"
	"github.com/aaron-vaz/github-branch-bot/pkg/notification"
)

const projectUpToDateText = "up to date with %s\n"

// BranchService is the main struct, it is used to start the application
type BranchService struct {
	Params *config.Params
	API    *github.APIService
	Msg    *notification.SlackService
	Wg     *sync.WaitGroup
}

// GenerateStatusMessage is used to start the application
func (b *BranchService) GenerateStatusMessage() string {
	sm := &notification.SlackMessage{
		Org:      b.Params.GithubOrganization,
		Messages: make(map[string][]string),
	}

	repos := b.API.GetReposInOrg(b.Params.GithubOrganization, b.Params.BaseBranch)

	if len(repos) == 0 {
		log.Printf("No branches in %s contain a default branch %s", b.Params.GithubOrganization, b.Params.BaseBranch)
		return ""
	}

	b.Wg.Add(len(repos))

	for _, repo := range repos {
		go b.processRepo(repo, sm)
	}

	b.Wg.Wait()

	return sm.String()
}

func (b *BranchService) processRepo(repo string, sm *notification.SlackMessage) {
	defer b.Wg.Done()

	branches := b.API.GetBranches(b.Params.GithubOrganization, repo, b.Params.HeadBranchPrefixes)

	if len(branches) == 0 {
		log.Printf("No branches of %s matched prefixes %s, check configuration", repo, b.Params.HeadBranchPrefixes)
		return
	}

	aheadBranches := b.API.GetAheadBy(b.Params.GithubOrganization, repo, b.Params.BaseBranch, branches)

	var branchMessages []string
	for branch, aheadBy := range aheadBranches {
		if message := b.Msg.GenerateMessage(repo, b.Params.BaseBranch, branch, aheadBy); message != "" {
			branchMessages = append(branchMessages, message)
		}
	}

	if len(branchMessages) > 0 {
		sm.Messages[repo] = branchMessages

	} else {
		sm.Messages[repo] = []string{fmt.Sprintf(projectUpToDateText, b.Params.BaseBranch)}
	}
}
