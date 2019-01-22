package config

import (
	"os"
	"strings"
)

// Params represents the configuration params that will be used by the services
type Params struct {
	GithubBaseURL      string
	GithubOrganization string
	GithubRepo         []string
	BaseBranch         string
	HeadBranchPrefixes []string
	WebhookURL         string
}

// ParseParams read the configuration parameters from environment variables and creates a Params struct to return
func ParseParams() *Params {
	return &Params{
		GithubBaseURL:      getEnv("GITHUB_BASE_URL", "http://localhost.com"),
		GithubOrganization: getEnv("GITHUB_ORGANISATION", ""),
		GithubRepo:         splitEnv("GITHUB_REPO", "", ","),
		BaseBranch:         getEnv("BASE_BRANCH", "develop"),
		HeadBranchPrefixes: splitEnv("HEAD_BRANCH_PREFIX", "master", ","),
		WebhookURL:         getEnv("GITHUB_BASE_URL", "http://localhost.com"),
	}
}

func splitEnv(key, fallback, delimeter string) []string {
	return strings.Split(getEnv(key, fallback), delimeter)
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}