package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/aaron-vaz/golang-utils/pkg/ioutils"
)

func TestHandleRequest(t *testing.T) {
	tests := []struct {
		name             string
		reposResponse    []byte
		branchesResponse []byte
		compareResponse  []byte
		messageWant      string
	}{
		{
			name:             "Happy Path Test",
			reposResponse:    readTestResource("repos-happy-path.json"),
			branchesResponse: readTestResource("branches-happy-path.json"),
			compareResponse:  readTestResource("ahead-happy-path.json"),
			messageWant:      `{"text":"*org branch check summary:*\n\n*test*:\nrelease is ahead of develop by 1 commits\n\n"}`,
		},
		{
			name:             "Error has occurred path",
			reposResponse:    readTestResource("invalid.json"),
			branchesResponse: readTestResource("invalid.json"),
			compareResponse:  readTestResource("invalid.json"),
			messageWant:      `{"text":"An error has occurred while performing the branch check"}`,
		},
	}

	for _, tt := range tests {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodGet {
				if strings.Contains(req.RequestURI, "orgs") {
					rw.Write(tt.reposResponse)

				} else if strings.Contains(req.RequestURI, "branches") {
					rw.Write(tt.branchesResponse)

				} else if strings.Contains(req.RequestURI, "compare") {
					rw.Write(tt.compareResponse)
				}

			} else if req.Method == http.MethodPost {
				defer ioutils.Close(req.Body)

				bytes, _ := ioutil.ReadAll(req.Body)
				actualMessage := string(bytes)
				if tt.messageWant != actualMessage {
					t.Errorf("Unexpected test result for HandleRequest want = %s, got = %s", tt.messageWant, actualMessage)
				}
			}
		}))

		os.Setenv("GITHUB_BASE_URL", server.URL)
		os.Setenv("GITHUB_TOKEN", "token")
		os.Setenv("GITHUB_ORGANISATION", "org")
		os.Setenv("BASE_BRANCH", "develop")
		os.Setenv("HEAD_BRANCH_PREFIX", "release")
		os.Setenv("WEBHOOK_URL", server.URL)

		HandleRequest()
	}
}

func readTestResource(path string) []byte {
	content, err := ioutil.ReadFile(filepath.Join("test-resources", path))
	if err != nil {
		panic(err)
	}

	return content
}
