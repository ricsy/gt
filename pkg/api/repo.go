package api

import "github.com/ricsy/gt/pkg/api/response"

// Repository is an alias for response.Repository
type Repository = response.Repository

// CreateRepoOptions is an alias for response.CreateRepoOptions
type CreateRepoOptions = response.CreateRepoOptions

// UpdateRepoOptions is an alias for response.UpdateRepoOptions
type UpdateRepoOptions = response.UpdateRepoOptions

// ListRepos lists user repositories
func (c *Client) ListRepos() ([]Repository, error) {
	var repos []Repository
	err := c.DoFromEndpoint(UserRepos.List, nil, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// GetRepo gets a single repository
func (c *Client) GetRepo(owner, repo string) (*Repository, error) {
	var result Repository
	err := c.DoFromEndpoint(Repo.Get, []interface{}{owner, repo}, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateRepo creates a new repository
func (c *Client) CreateRepo(opts CreateRepoOptions) (*Repository, error) {
	var result Repository
	err := c.DoFromEndpoint(UserRepos.List, nil, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListUserRepos lists repositories for a specific user
func (c *Client) ListUserRepos(username string) ([]Repository, error) {
	var repos []Repository
	err := c.DoFromEndpoint(UserRepos.List, []interface{}{username}, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// UpdateRepo updates a repository
func (c *Client) UpdateRepo(owner, repo string, opts UpdateRepoOptions) (*Repository, error) {
	var result Repository
	err := c.DoFromEndpoint(Repo.Get, []interface{}{owner, repo}, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
