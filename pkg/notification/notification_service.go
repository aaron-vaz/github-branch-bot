package notification

import (
	"bytes"
	"encoding/json"
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
