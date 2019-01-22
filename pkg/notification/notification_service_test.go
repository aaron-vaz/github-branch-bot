package notification

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSlackService_Notify(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "Test Happy path",
			message:  "Test Message",
			expected: `{"text":"Test Message"}`,
		},
		{
			name:     "Test empty message path",
			message:  "",
			expected: `{"text":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			received := false
			server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				if req.Method != http.MethodPost {
					t.Error("Request was not made via POST")
				}

				body, _ := ioutil.ReadAll(req.Body)

				if strings.TrimSpace(tt.expected) != strings.TrimSpace(string(body)) {
					t.Errorf("Request payloads didn't match, expected = %s, got = %s", tt.expected, body)
				}

				received = true
			}))

			service := &SlackService{server.URL, server.Client()}
			service.Notify(tt.message)

			if received == false {
				t.Error("No message received")
			}
		})
	}
}
