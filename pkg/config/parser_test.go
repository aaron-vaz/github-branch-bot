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
				os.Setenv("GITHUB_REPO", "repo")
				os.Setenv("BASE_BRANCH", "develop")
				os.Setenv("HEAD_BRANCH_PREFIX", "release")
				os.Setenv("WEBHOOK_URL", "http://localhost.com")
			},
			want: &Params{
				GithubBaseURL:      "http://localhost.com",
				GithubToken:        "token",
				GithubOrganization: "org",
				GithubRepo:         []string{"repo"},
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"release"},
				WebhookURL:         "http://localhost.com",
			},
		},

		{
			name: "Test Multiple params path",
			envSupplier: func() {
				os.Setenv("GITHUB_BASE_URL", "http://localhost.com")
				os.Setenv("GITHUB_TOKEN", "token")
				os.Setenv("GITHUB_ORGANISATION", "org")
				os.Setenv("GITHUB_REPO", "repo,repo2")
				os.Setenv("BASE_BRANCH", "develop")
				os.Setenv("HEAD_BRANCH_PREFIX", "release,master")
				os.Setenv("WEBHOOK_URL", "http://localhost.com")
			},
			want: &Params{
				GithubBaseURL:      "http://localhost.com",
				GithubToken:        "token",
				GithubOrganization: "org",
				GithubRepo:         []string{"repo", "repo2"},
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"release", "master"},
				WebhookURL:         "http://localhost.com",
			},
		},

		{
			name: "Test wrong delimeter path",
			envSupplier: func() {
				os.Setenv("GITHUB_BASE_URL", "http://localhost.com")
				os.Setenv("GITHUB_TOKEN", "token")
				os.Setenv("GITHUB_ORGANISATION", "org")
				os.Setenv("GITHUB_REPO", "repo:repo2")
				os.Setenv("BASE_BRANCH", "develop")
				os.Setenv("HEAD_BRANCH_PREFIX", "release:master")
				os.Setenv("WEBHOOK_URL", "http://localhost.com")
			},
			want: &Params{
				GithubBaseURL:      "http://localhost.com",
				GithubToken:        "token",
				GithubOrganization: "org",
				GithubRepo:         []string{"repo:repo2"},
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"release:master"},
				WebhookURL:         "http://localhost.com",
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
				GithubRepo:         []string{""},
				BaseBranch:         "develop",
				HeadBranchPrefixes: []string{"master"},
				WebhookURL:         "http://localhost.com",
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
	os.Setenv("GITHUB_REPO", "")
	os.Setenv("BASE_BRANCH", "")
	os.Setenv("HEAD_BRANCH_PREFIX", "")
	os.Setenv("WEBHOOK_URL", "")
}
