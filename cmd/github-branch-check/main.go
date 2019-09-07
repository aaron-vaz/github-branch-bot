package main

import (
	"errors"
	"net/http"
	"sync"

	"github.com/aaron-vaz/github-branch-bot/pkg/config"
	"github.com/aaron-vaz/github-branch-bot/pkg/github"
	"github.com/aaron-vaz/github-branch-bot/pkg/notification"
	"github.com/aaron-vaz/github-branch-bot/pkg/service"
	"github.com/aws/aws-lambda-go/lambda"
)

// Event is the struct for the lambda event
type Event struct {
	Query map[string]string `json:"query"`
}

// HandleRequest is the main entry point to the application, it will be executed by the AWS
func HandleRequest(request Event) error {
	params := config.ParseParams()
	githubAPI := &github.APIService{BaseURL: params.GithubBaseURL, Token: params.GithubToken, Client: http.DefaultClient}
	slackAPI := &notification.SlackService{Client: http.DefaultClient}

	branchService := &service.BranchService{
		Params: params,
		API:    githubAPI,
		Msg:    slackAPI,
		Wg:     &sync.WaitGroup{},
	}

	// first check validation token
	token := request.Query["token"]
	if token != params.SlackCommandToken {
		return errors.New("Incorrect validation token")
	}

	// then check for respose webhook url
	responseURL := request.Query["response_url"]
	if responseURL == "" {
		return errors.New("No response_url provided")
	}

	// provide response to stop slack command from timing out
	slackAPI.Notify(responseURL, "Processing request...")

	// do branch check
	if message := branchService.GenerateStatusMessage(); message != "" {
		slackAPI.Notify(responseURL, message)
		return nil
	}

	errorMsg := "Error occurred while processing request"

	slackAPI.Notify(responseURL, errorMsg)
	return errors.New(errorMsg)
}

func main() {
	lambda.Start(HandleRequest)
}
