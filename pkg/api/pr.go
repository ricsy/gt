package api

import "fmt"

// PRState represents the state of a pull request.
type PRState string

const (
	PRStateOpen   PRState = "open"
	PRStateClosed PRState = "closed"
	PRStateMerged PRState = "merged"
	PRStateAll    PRState = "all"
)

// PullRequest represents a Gitee pull request
type PullRequest struct {
	ID      int64  `json:"id"`
	Number  int    `json:"number"`
	State   string `json:"state"`
	Title   string `json:"title"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
	User    struct {
		ID    int64  `json:"id"`
		Login string `json:"login"`
	} `json:"user"`
	Labels       []Label `json:"labels"`
	Assignee     *User   `json:"assignee"`
	Assignees    []User  `json:"assignees"`
	Comments     int     `json:"comments"`
	Commits      int     `json:"commits"`
	Additions    int     `json:"additions"`
	Deletions    int     `json:"deletions"`
	FilesChanged int     `json:"files_changed"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	ClosedAt     string  `json:"closed_at"`
	MergedAt     string  `json:"merged_at"`
	Merged       bool    `json:"merged"`
	Mergeable    bool    `json:"mergeable"`
	Head         struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	} `json:"head"`
	Base struct {
		Ref string `json:"ref"`
		Sha string `json:"sha"`
	} `json:"base"`
}

// ListPRs lists pull requests in a repository
func (c *Client) ListPRs(owner, repo, state string) ([]PullRequest, error) {
	path := fmt.Sprintf(apiPathPRs+"?state=%s", owner, repo, state)
	var prs []PullRequest
	err := c.Do("GET", path, nil, &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

// GetPR gets a single pull request
func (c *Client) GetPR(owner, repo string, number int) (*PullRequest, error) {
	path := fmt.Sprintf(apiPathPRs+"/%d", owner, repo, number)
	var pr PullRequest
	err := c.Do("GET", path, nil, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// CreatePRRequest is the request body for creating a PR
type CreatePRRequest struct {
	Title string `json:"title"`
	Body  string `json:"body,omitempty"`
	Head  string `json:"head"`
	Base  string `json:"base"`
}

// CreatePR creates a new pull request
func (c *Client) CreatePR(owner, repo, title, body, head, base string) (*PullRequest, error) {
	path := fmt.Sprintf(apiPathPRs, owner, repo)
	req := CreatePRRequest{Title: title, Body: body, Head: head, Base: base}
	var pr PullRequest
	err := c.Do("POST", path, req, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// MergePR merges a pull request
func (c *Client) MergePR(owner, repo string, number int) error {
	path := fmt.Sprintf(apiPathPRs+"/%d/merge", owner, repo, number)
	return c.Do("PUT", path, nil, nil)
}

// UpdatePRState updates a pull request's state (open/closed)
func (c *Client) UpdatePRState(owner, repo string, number int, state PRState) (*PullRequest, error) {
	path := fmt.Sprintf(apiPathPRs+"/%d", owner, repo, number)
	updateReq := map[string]string{"state": string(state)}
	var pr PullRequest
	err := c.Do("PATCH", path, updateReq, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// CreatePRComment adds a comment to a pull request
func (c *Client) CreatePRComment(owner, repo string, number int, body string) error {
	path := fmt.Sprintf(apiPathPRs+"/%d/comments", owner, repo, number)
	commentReq := map[string]string{"body": body}
	return c.Do("POST", path, commentReq, nil)
}
