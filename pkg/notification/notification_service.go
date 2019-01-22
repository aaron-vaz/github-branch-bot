package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/aaron-vaz/golang-utils/pkg/util"
)

// SlackService provides operations that allow you to post notifications to slack
type SlackService struct {
	URL string
	*http.Client
}

// GenerateMessage build a mesage that will be posted to the slack channel
func (service *SlackService) GenerateMessage(repo, base, head string, aheadBy int) string {
	message := "*%s*:\n"
	message += "%s is ahead of %s by %d commits\n"
	return fmt.Sprintf(message, repo, head, base, aheadBy)
}

// Notify sends slack message in the form of a json payload to the URL provided
func (service *SlackService) Notify(message string) {
	payload, err := json.Marshal(map[string]string{"text": message})
	util.ErrCheck(err, false)

	res, err := service.Post(service.URL, "application/json", bytes.NewReader(payload))
	util.ErrCheck(err, false)

	if res != nil {
		defer util.Close(res.Body)
	}

	body, err := ioutil.ReadAll(res.Body)
	util.ErrCheck(err, false)

	log.Printf("Slack response: %s", string(body))
}
