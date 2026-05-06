package api

import (
	"fmt"

	"github.com/ricsy/gt/pkg/util"
)

// IssueState represents the state of an issue.
type IssueState string

const (
	IssueStateOpen        IssueState = "open"
	IssueStateClosed      IssueState = "closed"
	IssueStateProgressing IssueState = "progressing"
	IssueStateRejected    IssueState = "rejected"
	IssueStateAll         IssueState = "all"
)

// Issue represents a Gitee issue
type Issue struct {
	ID      int64  `json:"id"`
	Number  int64  `json:"number"`
	State   string `json:"state"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	} `json:"user"`
	Labels      []Label `json:"labels"`
	Assignee    *User   `json:"assignee"`
	Assignees   []User  `json:"assignees"`
	Comments    int     `json:"comments"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	ClosedAt    string  `json:"closed_at"`
	PullRequest *struct {
		URL     string `json:"url"`
		HTMLURL string `json:"html_url"`
	} `json:"pull_request"`
}

// Label represents a Gitee label
type Label struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

// User represents a Gitee user
type User struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// ListIssuesOptions contains the optional parameters for ListIssues
type ListIssuesOptions struct {
	State        string
	Labels       string
	Sort         string // created, updated
	Direction    string // asc, desc
	Milestone    string
	Assignee     string
	Creator      string
	Program      string
	Q            string
	SecurityHole *bool
	Page         int
	PerPage      int
}

// ListIssues lists issues in a repository
func (c *Client) ListIssues(owner, repo string, opts ListIssuesOptions) ([]Issue, error) {
	path := fmt.Sprintf(apiPathIssues, owner, repo)
	query := buildQuery(opts)
	if query != "" {
		path += "?" + query
	}
	var issues []Issue
	err := c.Do("GET", path, nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
}

// buildQuery builds query string from ListIssuesOptions
func buildQuery(opts ListIssuesOptions) string {
	params := []string{
		"state", opts.State,
		"labels", opts.Labels,
		"sort", opts.Sort,
		"direction", opts.Direction,
		"milestone", opts.Milestone,
		"assignee", opts.Assignee,
		"creator", opts.Creator,
		"program", opts.Program,
		"q", opts.Q,
	}
	if opts.SecurityHole != nil {
		params = append(params, "security_hole", fmt.Sprintf("%t", *opts.SecurityHole))
	}
	if opts.Page > 0 {
		params = append(params, "page", fmt.Sprintf("%d", opts.Page))
	}
	if opts.PerPage > 0 {
		params = append(params, "per_page", fmt.Sprintf("%d", opts.PerPage))
	}
	return util.BuildQuery(params...)
}

// GetIssue gets a single issue
func (c *Client) GetIssue(owner, repo, number string) (*Issue, error) {
	path := fmt.Sprintf(apiPathIssues+"/%s", owner, repo, number)
	var issue Issue
	err := c.Do("GET", path, nil, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// CreateIssueRequest is the request body for creating an issue
type CreateIssueRequest struct {
	Title         string `json:"title"`
	Body          string `json:"body,omitempty"`
	Repo          string `json:"repo,omitempty"`
	IssueType     string `json:"issue_type,omitempty"`
	Assignee      string `json:"assignee,omitempty"`
	Collaborators string `json:"collaborators,omitempty"`
	Milestone     int    `json:"milestone,omitempty"`
	Labels        string `json:"labels,omitempty"`
	Program       string `json:"program,omitempty"`
	SecurityHole  bool   `json:"security_hole,omitempty"`
	Branch        string `json:"branch,omitempty"`
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(owner, repo, title, body string) (*Issue, error) {
	path := fmt.Sprintf(apiPathIssueCreate, owner)
	req := CreateIssueRequest{Title: title, Body: body, Repo: repo}
	var issue Issue
	err := c.Do("POST", path, req, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// UpdateIssueRequest is the request body for updating an issue
type UpdateIssueRequest struct {
	Repo          string `json:"repo,omitempty"`
	Title         string `json:"title,omitempty"`
	Body          string `json:"body,omitempty"`
	State         string `json:"state,omitempty"`
	Assignee      string `json:"assignee,omitempty"`
	Collaborators string `json:"collaborators,omitempty"`
	Milestone     int    `json:"milestone,omitempty"`
	Labels        string `json:"labels,omitempty"`
	Program       string `json:"program,omitempty"`
	SecurityHole  bool   `json:"security_hole,omitempty"`
	Branch        string `json:"branch,omitempty"`
}

// UpdateIssue updates an issue
func (c *Client) UpdateIssue(owner, repo, number string, req UpdateIssueRequest) (*Issue, error) {
	path := fmt.Sprintf(apiPathIssueUpdate, owner, number)
	var issue Issue
	err := c.Do("PATCH", path, req, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// UpdateIssueState updates an issue's state (open/closed/progressing)
// Note: rejected state is not supported by the API for updates
func (c *Client) UpdateIssueState(owner, repo, number string, state IssueState) (*Issue, error) {
	path := fmt.Sprintf(apiPathIssueUpdate, owner, number)
	req := map[string]string{"state": string(state), "repo": repo}
	var issue Issue
	err := c.Do("PATCH", path, req, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// IssueComment represents a comment on an issue
type IssueComment struct {
	ID        int64  `json:"id"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

// ListIssueComments lists comments for an issue
func (c *Client) ListIssueComments(owner, repo, number string) ([]IssueComment, error) {
	path := fmt.Sprintf(apiPathIssues+"/%s/comments", owner, repo, number)
	var comments []IssueComment
	err := c.Do("GET", path, nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreateIssueComment creates a comment on an issue
func (c *Client) CreateIssueComment(owner, repo, number, body string) (*IssueComment, error) {
	path := fmt.Sprintf(apiPathIssues+"/%s/comments", owner, repo, number)
	req := map[string]string{"body": body}
	var comment IssueComment
	err := c.Do("POST", path, req, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
