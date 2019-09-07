package main

import (
	"net/http"
	"sync"

	"github.com/aaron-vaz/github-branch-bot/pkg/config"
	"github.com/aaron-vaz/github-branch-bot/pkg/github"
	"github.com/aaron-vaz/github-branch-bot/pkg/notification"
	"github.com/aaron-vaz/github-branch-bot/pkg/service"
	"github.com/aws/aws-lambda-go/lambda"
)

// HandleRequest is the main entry point to the application, it will be executed by the AWS
func HandleRequest() {
	params := config.ParseParams()
	githubAPI := &github.APIService{BaseURL: params.GithubBaseURL, Token: params.GithubToken, Client: http.DefaultClient}
	slackAPI := &notification.SlackService{Client: http.DefaultClient}

	branchService := &service.BranchService{
		Params: params,
		API:    githubAPI,
		Msg:    slackAPI,
		Wg:     &sync.WaitGroup{},
	}

	if message := branchService.GenerateStatusMessage(); message != "" {
		slackAPI.Notify(params.WebhookURL, message)

	} else {
		slackAPI.Notify(params.WebhookURL, "An error has occurred while performing the branch check")
	}
}

func main() {
	lambda.Start(HandleRequest)
}
