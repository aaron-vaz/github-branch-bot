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
		name      string
		message   string
		delivered bool
		expected  string
	}{
		{
			name:      "Test Happy path",
			message:   "Test Message",
			delivered: true,
			expected:  `{"text":"Test Message"}`,
		},
		{
			name:      "Test empty message path",
			message:   "",
			delivered: false,
			expected:  "",
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

			service := &SlackService{server.Client()}
			service.Notify(server.URL, tt.message)

			if received != tt.delivered {
				t.Errorf("Request delivery didnt match expected, want = %t, got = %t", tt.delivered, received)
			}
		})
	}
}

func TestSlackService_GenerateMessage(t *testing.T) {
	type args struct {
		repo    string
		base    string
		head    string
		aheadBy int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Happy Path",
			args: args{
				repo:    "test",
				base:    "develop",
				head:    "master",
				aheadBy: 5,
			},
			want: "master is ahead of develop by 5 commits\n",
		},

		{
			name: "Test branches up to date path",
			args: args{
				repo:    "test",
				base:    "develop",
				head:    "master",
				aheadBy: 0,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &SlackService{}
			if got := service.GenerateMessage(tt.args.repo, tt.args.base, tt.args.head, tt.args.aheadBy); got != tt.want {
				t.Errorf("SlackService.GenerateMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlackMessage_String(t *testing.T) {
	tests := []struct {
		name string
		sm   *SlackMessage
		want string
	}{
		{
			name: "Test happy path",
			sm: &SlackMessage{
				Org:      "Organisation",
				Messages: map[string][]string{"test repo": []string{"message 1", "message 2"}},
			},
			want: "*Organisation branch check summary:*\n\n*test repo*:\nmessage 1\nmessage 2\n",
		},
		{
			name: "Test 1 message path",
			sm: &SlackMessage{
				Org:      "Organisation",
				Messages: map[string][]string{"test repo": []string{"message 1"}},
			},
			want: "*Organisation branch check summary:*\n\n*test repo*:\nmessage 1\n",
		},
		{
			name: "Test no org path",
			sm: &SlackMessage{
				Org:      "",
				Messages: map[string][]string{"test repo": []string{"message 1"}},
			},
			want: "",
		},
		{
			name: "Test no messages path",
			sm: &SlackMessage{
				Org:      "Organisation",
				Messages: map[string][]string{},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sm.String(); got != tt.want {
				t.Errorf("SlackMessage.String() = %q, want %q", got, tt.want)
			}
		})
	}
}
