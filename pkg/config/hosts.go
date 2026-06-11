package config

import "fmt"

// HostAuth holds authentication info for a host
type HostAuth struct {
	Token string `yaml:"token"`
	User  string `yaml:"user"`
}

const DefaultHost = "gitee.com"

// ApiUrl returns the base API URL for a given host.
func ApiUrl(host string) string {
	return WebURL(host) + "/api/v5"
}

// WebURL returns the browser-facing base URL for a given host.
func WebURL(host string) string {
	return "https://" + host
}

// RepoGitHTTPSURL returns the HTTPS clone URL for a repository.
func RepoGitHTTPSURL(host, owner, repo string) string {
	return fmt.Sprintf("%s/%s/%s.git", WebURL(host), owner, repo)
}
