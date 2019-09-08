package service

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/aaron-vaz/github-branch-bot/pkg/config"
	"github.com/aaron-vaz/github-branch-bot/pkg/github"
	"github.com/aaron-vaz/github-branch-bot/pkg/notification"
)

func TestBranchService_GenerateStatusMessage(t *testing.T) {
	tests := []struct {
		name             string
		reposResponse    []byte
		branchesResponse []byte
		compareResponse  []byte
		messageDelivered bool
		messageWant      string
	}{
		{
			name:             "Happy Path Test",
			reposResponse:    readTestResource("repos-happy-path.json"),
			branchesResponse: readTestResource("branches-happy-path.json"),
			compareResponse:  readTestResource("ahead-happy-path.json"),
			messageDelivered: true,
			messageWant:      "*org branch check summary:*\n\n*test*:\nmaster is ahead of develop by 1 commits\n\n",
		},
		{
			name:             "Test Branches inline path",
			reposResponse:    readTestResource("repos-happy-path.json"),
			branchesResponse: readTestResource("branches-happy-path.json"),
			compareResponse:  readTestResource("inline-happy-path.json"),
			messageDelivered: true,
			messageWant:      "*org branch check summary:*\n\n*test*:\nup to date with develop\n\n",
		},
		{
			name:             "Test No matched repos",
			reposResponse:    readTestResource("invalid.json"),
			branchesResponse: readTestResource("branches-happy-path.json"),
			compareResponse:  readTestResource("ahead-happy-path.json"),
			messageDelivered: false,
			messageWant:      "",
		},
		{
			name:             "Test No matched branches",
			reposResponse:    readTestResource("repos-happy-path.json"),
			branchesResponse: readTestResource("invalid.json"),
			compareResponse:  readTestResource("ahead-happy-path.json"),
			messageDelivered: false,
			messageWant:      "",
		},
		{
			name:             "Test No response path",
			reposResponse:    readTestResource("invalid.json"),
			branchesResponse: readTestResource("invalid.json"),
			compareResponse:  readTestResource("invalid.json"),
			messageDelivered: false,
			messageWant:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.Contains(req.RequestURI, "orgs") {
					rw.Write(tt.reposResponse)

				} else if strings.Contains(req.RequestURI, "branches") {
					rw.Write(tt.branchesResponse)

				} else if strings.Contains(req.RequestURI, "compare") {
					rw.Write(tt.compareResponse)
				}
			}))

			params := &config.Params{
				GithubToken:        "token",
				GithubOrganization: "org",
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"master"},
				WebhookURL:         server.URL,
			}

			githubAPI := &github.APIService{BaseURL: server.URL, Client: server.Client()}
			slackAPI := &notification.SlackService{Client: server.Client()}

			bot := &BranchService{
				Params: params,
				API:    githubAPI,
				Msg:    slackAPI,
				Wg:     &sync.WaitGroup{},
			}

			actualMessage := bot.GenerateStatusMessage()

			if actualMessage != tt.messageWant {
				t.Errorf("Unexpected test result for GenerateStatusMessage want = %s, got = %s", tt.messageWant, actualMessage)
			}
		})
	}
}

func readTestResource(path string) []byte {
	content, err := ioutil.ReadFile(filepath.Join("test-resources", path))
	if err != nil {
		panic(err)
	}

	return content
}
