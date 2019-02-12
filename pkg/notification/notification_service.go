package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aaron-vaz/golang-utils/pkg/errorutil"
	"github.com/aaron-vaz/golang-utils/pkg/ioutils"
)

// SlackService provides operations that allow you to post notifications to slack
type SlackService struct {
	URL string
	*http.Client
}

// SlackMessage is used to build the message we will be posting to slack
type SlackMessage struct {
	Org      string
	Messages []string
}

func (sm *SlackMessage) String() string {
	if sm.Org == "" || len(sm.Messages) == 0 {
		return ""
	}

	ret := fmt.Sprintf("*%s branch check summary:*\n", sm.Org)
	ret += "\n"

	for _, message := range sm.Messages {
		ret += message
		ret += "\n"
	}

	return ret
}

// GenerateMessage build a mesage that will be posted to the slack channel
func (service *SlackService) GenerateMessage(repo, base, head string, aheadBy int) string {
	message := "*%s*:\n"
	message += "%s is ahead of %s by %d commits\n"
	return fmt.Sprintf(message, repo, head, base, aheadBy)
}

// Notify sends slack message in the form of a json payload to the URL provided
func (service *SlackService) Notify(message string) {
	if message == "" {
		log.Println("No message received, notification will not be performed")
		return
	}

	payload, err := json.Marshal(map[string]string{"text": message})
	errorutil.ErrCheck(err, false)

	res, err := service.Post(service.URL, "application/json", bytes.NewReader(payload))
	errorutil.ErrCheck(err, false)

	if res != nil {
		defer ioutils.Close(res.Body)
	}

	body, err := ioutil.ReadAll(res.Body)
	errorutil.ErrCheck(err, false)

	log.Printf("Slack response: %s", string(body))
}
