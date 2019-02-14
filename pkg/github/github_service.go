package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/aaron-vaz/golang-utils/pkg/errorutil"
	"github.com/aaron-vaz/golang-utils/pkg/ioutils"
)

const (
	getReposInOrgPath   = "/orgs/%s/repos?per_page=100"
	getBranchesPath     = "/repos/%s/%s/branches?per_page=100"
	compareBranchesPath = "/repos/%s/%s/compare/%s...%s"
)

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

// GetReposInOrg returns a list projects that contain the configured base branch as their default branch
func (service *APIService) GetReposInOrg(org, baseBranch string) []string {
	url := service.BaseURL + fmt.Sprintf(getReposInOrgPath, org)
	responses := []Response{}

	service.getGithubResponse(url, &responses)

	var repos []string
	for _, response := range responses {
		if response.DefaultBranch == baseBranch {
			repos = append(repos, response.Name)
		}
	}

	return repos
}

// GetBranches return all the branches matching the supplied prefix
// if no prefix is supplied it returns all the branches from the repo
func (service *APIService) GetBranches(owner, repo string, prefix []string) []string {
	url := service.BaseURL + fmt.Sprintf(getBranchesPath, owner, repo)
	responses := []Response{}

	service.getGithubResponse(url, &responses)

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

// GetAheadBy returns how many commits the supplied branch is ahead of the supplied base branch
func (service *APIService) GetAheadBy(owner, repo, base string, heads []string) map[string]int {
	results := make(map[string]int)
	for _, head := range heads {
		log.Printf("Checking %s branch %s", repo, head)
		url := service.BaseURL + fmt.Sprintf(compareBranchesPath, owner, repo, base, head)
		response := &CompareBranches{}

		service.getGithubResponse(url, response)
		results[head] = response.Ahead
	}

	return results
}

func (service *APIService) getGithubResponse(url string, model interface{}) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	errorutil.ErrCheck(err, false)

	req.Header.Add("Authorization", "token "+service.Token)

	res, err := service.Do(req)
	errorutil.ErrCheck(err, false)

	if res != nil {
		defer ioutils.Close(res.Body)
	}

	content, err := ioutil.ReadAll(res.Body)
	errorutil.ErrCheck(err, false)

	errorutil.ErrCheck(json.Unmarshal(content, &model), false)
}
