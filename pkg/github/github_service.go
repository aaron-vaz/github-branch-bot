package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/aaron-vaz/golang-utils/pkg/errorutil"
	"github.com/aaron-vaz/golang-utils/pkg/ioutils"
)

const (
	getRepositoriesInOrgPath = "/orgs/%s/repos"
	getBranchesPath          = "/repos/%s/%s/branches"
	compareBranchesPath      = "/repos/%s/%s/compare/%s...%s"

	authorizationHeader = "Authorization"
	contentTypeHeader   = "Content-Type"
	linkHeader          = "Link"

	jsonMediaType = "application/json; charset=utf-8"

	tokenHeaderPrefix = "token "
)

var linkHeaderRegex = regexp.MustCompile("<([^>]+)>;\\srel=\"next\"+")

// Response is the struct that represents the github branches response
type Response struct {
	Name          string `json:"name"`
	DefaultBranch string `json:"default_branch"`
}

// CompareBranches is the struct that represents the github compare branches response
type CompareBranches struct {
	Ahead int `json:"ahead_by"`
}

// APIService is a service that provides operations allowing you to interact with github api
type APIService struct {
	BaseURL string
	Token   string
	*http.Client
}

// GetRepositoriesInOrg returns a list projects that contain the configured base branch as their default branch
func (s *APIService) GetRepositoriesInOrg(org, baseBranch string) []string {
	url := s.BaseURL + fmt.Sprintf(getRepositoriesInOrgPath, org)
	responses := s.executePaginatedGithubRequest(url)

	var repositories []string
	for _, response := range responses {
		if response.DefaultBranch == baseBranch {
			repositories = append(repositories, response.Name)
		}
	}

	return repositories
}

// GetBranches return all the branches matching the supplied prefix
// if no prefix is supplied it returns all the branches from the repo
func (s *APIService) GetBranches(owner, repo string, prefix []string) []string {
	url := s.BaseURL + fmt.Sprintf(getBranchesPath, owner, repo)
	responses := s.executePaginatedGithubRequest(url)

	var branches []string
	for _, value := range responses {
		for _, branch := range prefix {
			if branch == "" || strings.HasPrefix(value.Name, branch) {
				branches = append(branches, value.Name)
			}
		}
	}

	return branches
}

func (s *APIService) executePaginatedGithubRequest(url string) []Response {
	body, nextURL := s.executeGithubRequest(url)

	responses := []Response{}
	errorutil.ErrCheck(json.Unmarshal(body, &responses), false)

	if nextURL == "" {
		return responses

	} else {
		if nextPage := s.executePaginatedGithubRequest(nextURL); len(nextPage) > 0 {
			responses = append(responses, nextPage...)
		}
	}

	return responses
}

// GetAheadBy returns how many commits the supplied branch is ahead of the supplied base branch
func (s *APIService) GetAheadBy(owner, repo, base string, heads []string) map[string]int {
	results := make(map[string]int)
	for _, head := range heads {
		log.Printf("Checking %s branch %s", repo, head)
		url := s.BaseURL + fmt.Sprintf(compareBranchesPath, owner, repo, base, head)
		response := &CompareBranches{}

		body, _ := s.executeGithubRequest(url)
		errorutil.ErrCheck(json.Unmarshal(body, response), false)

		results[head] = response.Ahead
	}

	return results
}

func (s *APIService) executeGithubRequest(url string) ([]byte, string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	errorutil.ErrCheck(err, false)

	req.Header.Add(authorizationHeader, tokenHeaderPrefix+s.Token)
	req.Header.Add(contentTypeHeader, jsonMediaType)

	res, err := s.Do(req)
	errorutil.ErrCheck(err, false)

	if res != nil {
		defer ioutils.Close(res.Body)
	}

	body, err := ioutil.ReadAll(res.Body)
	errorutil.ErrCheck(err, false)

	return body, s.getNextLink(res.Header.Get(linkHeader))
}

func (s *APIService) getNextLink(header string) string {
	if matches := linkHeaderRegex.FindStringSubmatch(header); len(matches) > 1 {
		if value := strings.TrimSpace(matches[1]); value != "" {
			return value
		}
	}

	return ""
}
