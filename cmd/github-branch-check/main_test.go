package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleRequest(t *testing.T) {
	var reposResponse []byte
	var branchesResponse []byte
	var compareResponse []byte

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if strings.Contains(req.RequestURI, "orgs") {
			rw.Write(reposResponse)

		} else if strings.Contains(req.RequestURI, "branches") {
			rw.Write(branchesResponse)

		} else if strings.Contains(req.RequestURI, "compare") {
			rw.Write(compareResponse)
		}
	}))

	tests := []struct {
		name         string
		args         Event
		responseFunc func()
		wantErr      bool
	}{
		{
			name:    "No validation token",
			args:    Event{},
			wantErr: true,
		},
		{
			name: "No response url",
			args: Event{
				Query: map[string]string{
					"token": "token",
				},
			},
			wantErr: true,
		},
		{
			name: "Error Path",
			args: Event{
				Query: map[string]string{
					"token":        "token",
					"response_url": server.URL,
				},
			},
			wantErr: true,
		},
		{
			name: "Happy Path",
			args: Event{
				Query: map[string]string{
					"token":        "token",
					"response_url": server.URL,
				},
			},
			responseFunc: func() {
				reposResponse = readTestResource("repos-happy-path.json")
				branchesResponse = readTestResource("branches-happy-path.json")
				compareResponse = readTestResource("ahead-happy-path.json")
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("GITHUB_BASE_URL", server.URL)
			os.Setenv("GITHUB_TOKEN", "token")
			os.Setenv("GITHUB_ORGANISATION", "org")
			os.Setenv("GITHUB_REPO", "repo")
			os.Setenv("BASE_BRANCH", "develop")
			os.Setenv("HEAD_BRANCH_PREFIX", "release")
			os.Setenv("SLACK_COMMAND_TOKEN", "token")

			if tt.responseFunc != nil {
				tt.responseFunc()
			}

			err := HandleRequest(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("HandleRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
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
