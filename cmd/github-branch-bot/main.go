package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/aaron-vaz/github-branch-bot/pkg/config"
	"github.com/aaron-vaz/github-branch-bot/pkg/github"
	"github.com/aaron-vaz/github-branch-bot/pkg/notification"
	"github.com/aws/aws-lambda-go/lambda"
)

const projectUpToDateText = "up to date with %s\n"

// Bot is the main struct, it is used to start the application
type Bot struct {
	params *config.Params
	api    *github.APIService
	msg    *notification.SlackService
	wg     *sync.WaitGroup
}

// Start is used to start the application
func (b *Bot) Start() {
	sm := &notification.SlackMessage{
		Org:      b.params.GithubOrganization,
		Messages: make(map[string][]string),
	}

	repos := b.api.GetReposInOrg(b.params.GithubOrganization, b.params.BaseBranch)

	if len(repos) == 0 {
		log.Printf("No branches in %s contain a default branch %s", b.params.GithubOrganization, b.params.BaseBranch)
		return
	}

	b.wg.Add(len(repos))

	for _, repo := range repos {
		go b.processRepo(repo, sm)
	}

	b.wg.Wait()

	b.msg.Notify(sm.String())
}

func (b *Bot) processRepo(repo string, sm *notification.SlackMessage) {
	defer b.wg.Done()

	branches := b.api.GetBranches(b.params.GithubOrganization, repo, b.params.HeadBranchPrefixes)

	if len(branches) == 0 {
		log.Printf("No branches of %s matched prefixes %s, check configuration", repo, b.params.HeadBranchPrefixes)
		return
	}

	aheadBranches := b.api.GetAheadBy(b.params.GithubOrganization, repo, b.params.BaseBranch, branches)

	var branchMessages []string
	for branch, aheadBy := range aheadBranches {
		if message := b.msg.GenerateMessage(repo, b.params.BaseBranch, branch, aheadBy); message != "" {
			branchMessages = append(branchMessages, message)
		}
	}

	if len(branchMessages) > 0 {
		sm.Messages[repo] = branchMessages

	} else {
		sm.Messages[repo] = []string{fmt.Sprintf(projectUpToDateText, b.params.BaseBranch)}
	}
}

// HandleRequest is the main entry point to the application, it will be executed by the AWS
func HandleRequest() {
	params := config.ParseParams()
	githubAPI := &github.APIService{BaseURL: params.GithubBaseURL, Token: params.GithubToken, Client: http.DefaultClient}
	slackAPI := &notification.SlackService{URL: params.WebhookURL, Client: http.DefaultClient}

	bot := &Bot{
		params: params,
		api:    githubAPI,
		msg:    slackAPI,
		wg:     &sync.WaitGroup{},
	}

	bot.Start()
}

func main() {
	lambda.Start(HandleRequest)
}
