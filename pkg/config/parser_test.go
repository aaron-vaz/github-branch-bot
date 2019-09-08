package config

import (
	"os"
	"reflect"
	"testing"
)

func TestParseParams(t *testing.T) {
	tests := []struct {
		name        string
		envSupplier func()
		want        *Params
	}{
		{
			name: "Test Happy path",
			envSupplier: func() {
				os.Setenv("GITHUB_BASE_URL", "http://localhost.com")
				os.Setenv("GITHUB_TOKEN", "token")
				os.Setenv("GITHUB_ORGANISATION", "org")
				os.Setenv("BASE_BRANCH", "develop")
				os.Setenv("HEAD_BRANCH_PREFIX", "release")
				os.Setenv("WEBHOOK_URL", "http://localhost.com")
				os.Setenv("SLACK_COMMAND_TOKEN", "token")
			},
			want: &Params{
				GithubBaseURL:      "http://localhost.com",
				GithubToken:        "token",
				GithubOrganization: "org",
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"release"},
				WebhookURL:         "http://localhost.com",
				SlackCommandToken:  "token",
			},
		},

		{
			name: "Test Multiple params path",
			envSupplier: func() {
				os.Setenv("GITHUB_BASE_URL", "http://localhost.com")
				os.Setenv("GITHUB_TOKEN", "token")
				os.Setenv("GITHUB_ORGANISATION", "org")
				os.Setenv("BASE_BRANCH", "develop")
				os.Setenv("HEAD_BRANCH_PREFIX", "release,master")
				os.Setenv("WEBHOOK_URL", "http://localhost.com")
				os.Setenv("SLACK_COMMAND_TOKEN", "token")
			},
			want: &Params{
				GithubBaseURL:      "http://localhost.com",
				GithubToken:        "token",
				GithubOrganization: "org",
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"release", "master"},
				WebhookURL:         "http://localhost.com",
				SlackCommandToken:  "token",
			},
		},

		{
			name: "Test wrong delimeter path",
			envSupplier: func() {
				os.Setenv("GITHUB_BASE_URL", "http://localhost.com")
				os.Setenv("GITHUB_TOKEN", "token")
				os.Setenv("GITHUB_ORGANISATION", "org")
				os.Setenv("BASE_BRANCH", "develop")
				os.Setenv("HEAD_BRANCH_PREFIX", "release:master")
				os.Setenv("WEBHOOK_URL", "http://localhost.com")
				os.Setenv("SLACK_COMMAND_TOKEN", "token")
			},
			want: &Params{
				GithubBaseURL:      "http://localhost.com",
				GithubToken:        "token",
				GithubOrganization: "org",
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"release:master"},
				WebhookURL:         "http://localhost.com",
				SlackCommandToken:  "token",
			},
		},

		{
			name: "Test no environment variables path",
			envSupplier: func() {
				clearEnvs()
			},
			want: &Params{
				GithubBaseURL:      "http://localhost.com",
				GithubToken:        "",
				GithubOrganization: "",
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"master"},
				WebhookURL:         "http://localhost.com",
				SlackCommandToken:  "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// set environment variables
			tt.envSupplier()

			if got := ParseParams(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseParams() = %v, want %v", got, tt.want)
			}

			clearEnvs()
		})
	}
}

func clearEnvs() {
	os.Setenv("GITHUB_BASE_URL", "")
	os.Setenv("GITHUB_TOKEN", "")
	os.Setenv("GITHUB_ORGANISATION", "")
	os.Setenv("BASE_BRANCH", "")
	os.Setenv("HEAD_BRANCH_PREFIX", "")
	os.Setenv("WEBHOOK_URL", "")
	os.Setenv("SLACK_COMMAND_TOKEN", "")
}
