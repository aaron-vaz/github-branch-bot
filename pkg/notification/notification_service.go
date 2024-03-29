package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/aaron-vaz/golang-utils/pkg/errorutil"
	"github.com/aaron-vaz/golang-utils/pkg/ioutils"
)

// SlackService provides operations that allow you to post notifications to slack
type SlackService struct {
	Client *http.Client
}

// SlackMessage is used to build the message we will be posting to slack
type SlackMessage struct {
	Org      string
	Messages map[string][]string
}

func (sm *SlackMessage) String() string {
	if sm.Org == "" || len(sm.Messages) == 0 {
		return ""
	}

	ret := fmt.Sprintf("*%s branch check summary:*\n", sm.Org)
	ret += "\n"

	var repos []string
	for repo := range sm.Messages {
		repos = append(repos, repo)
	}

	sort.Strings(repos)

	for _, repo := range repos {
		ret += fmt.Sprintf("*%s*:\n", repo)
		for _, message := range sm.Messages[repo] {
			ret += message
			ret += "\n"
		}
	}

	return ret
}

// GenerateMessage build a mesage that will be posted to the slack channel
func (service *SlackService) GenerateMessage(repo, base, head string, aheadBy int) string {
	var message string

	if aheadBy > 0 {
		log.Printf("%s branch %s is ahead of %s", repo, head, base)
		message += fmt.Sprintf("%s is ahead of %s by %d commits\n", head, base, aheadBy)
	}

	return message
}

// Notify sends slack message in the form of a json payload to the URL provided
func (service *SlackService) Notify(url, message string) {
	if message == "" {
		log.Println("No message received, notification will not be performed")
		return
	}

	payload, err := json.Marshal(map[string]string{"text": message})
	errorutil.ErrCheck(err, false)

	res, err := service.Client.Post(url, "application/json", bytes.NewReader(payload))
	errorutil.ErrCheck(err, false)

	if res != nil {
		defer ioutils.Close(res.Body)
	}

	body, err := ioutil.ReadAll(res.Body)
	errorutil.ErrCheck(err, false)

	log.Printf("Slack response: %s", body)
}
