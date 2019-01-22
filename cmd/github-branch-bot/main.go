package main

import (
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
	for _, repo := range b.GithubRepo {
		branches := b.GetBranches(b.GithubOrganization, repo, b.HeadBranchPrefixes)

		for _, branch := range branches {
			ahead := b.GetAheadBy(b.GithubOrganization, repo, b.BaseBranch, branch)

			if ahead != 0 {
				message := b.GenerateMessage(repo, b.BaseBranch, branch, ahead)
				b.Notify(message)
			}
		}
	}
}

func HandleRequest() {
	params := config.ParseParams()

	githubAPI := &github.APIService{
		BaseURL: params.GithubBaseURL,
		Client:  http.DefaultClient,
	}

	slackAPI := &notification.SlackService{
		URL:    params.WebhookURL,
		Client: http.DefaultClient,
	}

	bot := &Bot{params, githubAPI, slackAPI}
	bot.Start()
}

func main() {
	lambda.Start(HandleRequest)
}
