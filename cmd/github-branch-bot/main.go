package main

import (
	"log"
	"net/http"

	"github.com/aaron-vaz/github-branch-bot/pkg/config"
	"github.com/aaron-vaz/github-branch-bot/pkg/github"
	"github.com/aaron-vaz/github-branch-bot/pkg/notification"
	"github.com/aws/aws-lambda-go/lambda"
)

// Bot is the main struct, it is used to start the application
type Bot struct {
	*config.Params
	*github.APIService
	*notification.SlackService
}

// Start is used to start the application
func (b *Bot) Start() {
	sm := &notification.SlackMessage{
		Org:      b.GithubOrganization,
		Messages: make(map[string][]string),
	}

	for _, repo := range b.GithubRepo {
		branches := b.GetBranches(b.GithubOrganization, repo, b.HeadBranchPrefixes)

		if len(branches) == 0 {
			log.Printf("No branches of %s matched prefixes %s, check configuration", repo, b.HeadBranchPrefixes)
			continue
		}

		aheadBranches := b.GetAheadBy(b.GithubOrganization, repo, b.BaseBranch, branches)

		var branchMessages []string
		for branch, aheadBy := range aheadBranches {
			if message := b.GenerateMessage(repo, b.BaseBranch, branch, aheadBy); message != "" {
				branchMessages = append(branchMessages, message)
			}
		}

		if len(branchMessages) > 0 {
			sm.Messages[repo] = branchMessages

		} else {
			sm.Messages[repo] = []string{"up to date"}
		}

	}

	b.Notify(sm.String())
}

// HandleRequest is the main entry point to the application, it will be executed by the AWS
func HandleRequest() {
	params := config.ParseParams()
	githubAPI := &github.APIService{BaseURL: params.GithubBaseURL, Token: params.GithubToken, Client: http.DefaultClient}
	slackAPI := &notification.SlackService{URL: params.WebhookURL, Client: http.DefaultClient}

	bot := &Bot{params, githubAPI, slackAPI}
	bot.Start()
}

func main() {
	lambda.Start(HandleRequest)
}
