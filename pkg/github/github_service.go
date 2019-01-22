package github

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aaron-vaz/golang-utils/pkg/util"
)

const getBranchesPath = "/repos/%s/%s/branches"
const compareBranchesPath = "/repos/%s/%s/compare/%s...%s"

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
	*http.Client
	baseURL string
}

// GetBranches return all the branches matching the supplied prefix
// if no prefix is supplied it returns all the branches from the repo
func (service *APIService) GetBranches(owner string, repo string, prefix string) []string {
	url := service.baseURL + fmt.Sprintf(getBranchesPath, owner, repo)
	responseModel := []BranchModel{}

	service.getGithubResponse(url, &responseModel)

	var branches []string
	for _, value := range responseModel {
		if prefix == "" || strings.HasPrefix(value.Name, prefix) {
			branches = append(branches, value.Name)
		}
	}

	return branches
}

// GetAheadBy returns how many commits the supplied branch is ahead of the supplied base branch
func (service *APIService) GetAheadBy(owner string, repo string, base string, head string) int {
	url := service.baseURL + fmt.Sprintf(compareBranchesPath, owner, repo, base, head)
	responseModel := &CompareBranchesModel{}

	service.getGithubResponse(url, responseModel)

	return responseModel.Ahead
}

func (service *APIService) getGithubResponse(url string, model interface{}) {
	res, err := service.Get(url)
	util.ErrCheck(err, false)

	if res != nil {
		defer util.Close(res.Body)
	}

	content, err := ioutil.ReadAll(res.Body)
	util.ErrCheck(err, false)

	util.ErrCheck(json.Unmarshal(content, &model), false)
}
