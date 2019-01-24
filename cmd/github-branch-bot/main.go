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
		Messages: []string{},
	}

	for _, repo := range b.GithubRepo {
		branches := b.GetBranches(b.GithubOrganization, repo, b.HeadBranchPrefixes)

		if len(branches) == 0 {
			log.Printf("No branches matched prefixes %s, check configuration", b.HeadBranchPrefixes)
			continue
		}

		for _, branch := range branches {
			log.Printf("Checking %s branch %s", repo, branch)
			ahead := b.GetAheadBy(b.GithubOrganization, repo, b.BaseBranch, branch)

			if ahead != 0 {
				log.Printf("%s branch %s is ahead of %s", repo, branch, b.BaseBranch)
				sm.Messages = append(sm.Messages, b.GenerateMessage(repo, b.BaseBranch, branch, ahead))
			}
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
