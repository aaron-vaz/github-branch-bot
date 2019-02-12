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
	getBranchesPath     = "/repos/%s/%s/branches?per_page=100"
	compareBranchesPath = "/repos/%s/%s/compare/%s...%s"
)

// BranchModel is the struct that represents the github branches response
type BranchModel struct {
	Name string `json:"name"`
}

// CompareBranchesModel is the struct that represents the github compare branches response
type CompareBranchesModel struct {
	Ahead int `json:"ahead_by"`
}

// APIService is a service that provides operations allowing you to interact with github api
type APIService struct {
	BaseURL string
	Token   string
	*http.Client
}

// GetBranches return all the branches matching the supplied prefix
// if no prefix is supplied it returns all the branches from the repo
func (service *APIService) GetBranches(owner, repo string, prefix []string) []string {
	url := service.BaseURL + fmt.Sprintf(getBranchesPath, owner, repo)
	responseModel := []BranchModel{}

	service.getGithubResponse(url, &responseModel)

	var branches []string
	for _, value := range responseModel {
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
		responseModel := &CompareBranchesModel{}

		service.getGithubResponse(url, responseModel)
		results[head] = responseModel.Ahead
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
