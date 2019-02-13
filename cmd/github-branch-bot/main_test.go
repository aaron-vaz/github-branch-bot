package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aaron-vaz/github-branch-bot/pkg/config"
	"github.com/aaron-vaz/github-branch-bot/pkg/github"
	"github.com/aaron-vaz/github-branch-bot/pkg/notification"
)

func TestBot_Start(t *testing.T) {
	tests := []struct {
		name             string
		branchesResponse []byte
		compareResponse  []byte
		messageWant      string
	}{
		{
			name:             "Happy Path Test",
			branchesResponse: readTestResource("branches-happy-path.json"),
			compareResponse:  readTestResource("ahead-happy-path.json"),
			messageWant:      `{"text":"*org branch check summary:*\n\n*repo*:\nmaster is ahead of develop by 1 commits\n\n"}`,
		},
		{
			name:             "Test Branches inline path",
			branchesResponse: readTestResource("branches-happy-path.json"),
			compareResponse:  readTestResource("inline-happy-path.json"),
			messageWant:      `{"text":"*org branch check summary:*\n\n*repo*:\nup to date\n"}`,
		},
		{
			name:             "Test No response path",
			branchesResponse: readTestResource("invalid.json"),
			compareResponse:  readTestResource("invalid.json"),
			messageWant:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if strings.Contains(req.RequestURI, "branches") {
					rw.Write(tt.branchesResponse)

				} else if strings.Contains(req.RequestURI, "compare") {
					rw.Write(tt.compareResponse)
				} else {
					body, _ := ioutil.ReadAll(req.Body)

					if content := string(body); content != tt.messageWant {
						t.Errorf("Slack message not correct want = %s, got = %s", tt.messageWant, content)
					}

				}
			}))

			params := &config.Params{
				GithubToken:        "token",
				GithubOrganization: "org",
				GithubRepo:         []string{"repo"},
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"master"},
			}

			githubAPI := &github.APIService{BaseURL: server.URL, Client: server.Client()}
			slackAPI := &notification.SlackService{URL: server.URL, Client: server.Client()}

			bot := &Bot{params, githubAPI, slackAPI}
			bot.Start()
		})
	}
}

func TestHandleRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
	}))

	os.Setenv("GITHUB_BASE_URL", server.URL)
	os.Setenv("GITHUB_TOKEN", "token")
	os.Setenv("GITHUB_ORGANISATION", "org")
	os.Setenv("GITHUB_REPO", "repo")
	os.Setenv("BASE_BRANCH", "develop")
	os.Setenv("HEAD_BRANCH_PREFIX", "release")
	os.Setenv("WEBHOOK_URL", server.URL)

	HandleRequest()
}

func readTestResource(path string) []byte {
	content, err := ioutil.ReadFile(filepath.Join("test-resources", path))
	if err != nil {
		panic(err)
	}

	return content
}
