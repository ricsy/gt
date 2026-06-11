package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// SearchReposOptions is an alias for response.SearchReposOptions
type SearchReposOptions = response.SearchReposOptions

// SearchIssuesOptions is an alias for response.SearchIssuesOptions
type SearchIssuesOptions = response.SearchIssuesOptions

// SearchUsersOptions is an alias for response.SearchUsersOptions
type SearchUsersOptions = response.SearchUsersOptions

// SearchRepos searches repositories
func (c *Client) SearchRepos(opts SearchReposOptions) ([]Repository, error) {
	path := Search.SearchRepos.Path
	query := buildSearchQuery("q", opts.Q, opts.Sort, opts.Order, opts.Page, opts.PerPage,
		"owner", opts.Owner,
		"fork", opts.Fork,
		"language", opts.Language)
	var repos []Repository
	err := c.Do("GET", path+query, nil, &repos)
	if err != nil {
		return nil, err
	}
	return repos, nil
}

// SearchIssues searches issues
func (c *Client) SearchIssues(opts SearchIssuesOptions) ([]Issue, error) {
	path := Search.SearchIssues.Path
	query := buildSearchQuery("q", opts.Q, opts.Sort, opts.Order, opts.Page, opts.PerPage,
		"repo", opts.Repo,
		"language", opts.Language,
		"label", opts.Label,
		"state", opts.State,
		"author", opts.Author,
		"assignee", opts.Assignee)
	var issues []Issue
	err := c.Do("GET", path+query, nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

// SearchUsers searches users
func (c *Client) SearchUsers(opts SearchUsersOptions) ([]User, error) {
	path := Search.SearchUsers.Path
	query := buildSearchQuery("q", opts.Q, opts.Sort, opts.Order, opts.Page, opts.PerPage)
	var users []User
	err := c.Do("GET", path+query, nil, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func buildSearchQuery(qKey, qVal, sort, order string, page, perPage int, kvPairs ...interface{}) string {
	params := []string{qKey, qVal}
	for i := 0; i < len(kvPairs); i += 2 {
		if i+1 < len(kvPairs) {
			switch v := kvPairs[i+1].(type) {
			case string:
				if v != "" {
					params = append(params, kvPairs[i].(string), v)
				}
			case *bool:
				if v != nil {
					params = append(params, kvPairs[i].(string), strconv.FormatBool(*v))
				}
			}
		}
	}
	if sort != "" {
		params = append(params, "sort", sort)
	}
	if order != "" {
		params = append(params, "order", order)
	}
	if page > 0 {
		params = append(params, "page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params = append(params, "per_page", strconv.Itoa(perPage))
	}
	query := util.BuildQuery(params...)
	if query == "" {
		return ""
	}
	return "?" + query
}
