package api

import "fmt"

const (
	StateOpen   = "open"
	StateClosed = "closed"
)

// Issue represents a Gitee issue
type Issue struct {
	ID      int64  `json:"id"`
	Number  string `json:"number"`
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

// ListIssues lists issues in a repository
func (c *Client) ListIssues(owner, repo, state string, page, perPage int) ([]Issue, error) {
	path := fmt.Sprintf(apiPathIssues+"?state=%s&page=%d&per_page=%d", owner, repo, state, page, perPage)
	var issues []Issue
	err := c.Do("GET", path, nil, &issues)
	if err != nil {
		return nil, err
	}
	return issues, nil
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
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Repo  string `json:"repo,omitempty"`
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

// UpdateIssueState updates an issue's state (open/closed)
func (c *Client) UpdateIssueState(owner, repo, number, state string) (*Issue, error) {
	path := fmt.Sprintf(apiPathIssueUpdate, owner, number)
	req := map[string]string{"state": state, "repo": repo}
	var issue Issue
	err := c.Do("PATCH", path, req, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// CloseIssue closes an issue
func (c *Client) CloseIssue(owner, repo, number string) (*Issue, error) {
	return c.UpdateIssueState(owner, repo, number, StateClosed)
}

// ReopenIssue reopens a closed issue
func (c *Client) ReopenIssue(owner, repo, number string) (*Issue, error) {
	return c.UpdateIssueState(owner, repo, number, StateOpen)
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
