package api

import "fmt"

// Repository represents a Gitee repository
type Repository struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
		Name  string `json:"name"`
	} `json:"owner"`
	Description     string `json:"description"`
	Private         bool   `json:"private"`
	Fork            bool   `json:"fork"`
	HTMLURL         string `json:"html_url"`
	SSHURL          string `json:"ssh_url"`
	CloneURL        string `json:"clone_url"`
	DefaultBranch   string `json:"default_branch"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
	Homepage        string `json:"homepage"`
	Language        string `json:"language"`
	StarCount       int    `json:"stargazers_count"`
	WatchCount      int    `json:"watchers_count"`
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
}

// ListRepos lists user repositories
func (c *Client) ListRepos() ([]Repository, error) {
	var repos []Repository
	err := c.Do("GET", apiPathUserRepos, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// GetRepo gets a single repository
func (c *Client) GetRepo(owner, repo string) (*Repository, error) {
	var result Repository
	path := fmt.Sprintf(apiPathRepos, owner, repo)
	err := c.Do("GET", path, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateRepoOptions contains options for creating a repository
type CreateRepoOptions struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	AutoInit    bool   `json:"auto_init"`
}

// CreateRepo creates a new repository
func (c *Client) CreateRepo(opts CreateRepoOptions) (*Repository, error) {
	var result Repository
	err := c.Do("POST", "/user/repos", opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserRepos lists repositories for a specific user
func (c *Client) ListUserRepos(username string) ([]Repository, error) {
	var repos []Repository
	path := fmt.Sprintf(apiPathUsersRepos, username)
	err := c.Do("GET", path, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}
