package api

import "github.com/ricsy/gt/pkg/api/response"

// PRState is an alias for response.PRState
type PRState = response.PRState

const (
	PRStateOpen   = response.PRStateOpen
	PRStateClosed = response.PRStateClosed
	PRStateMerged = response.PRStateMerged
	PRStateAll    = response.PRStateAll
)

// PullRequest is an alias for response.PullRequest
type PullRequest = response.PullRequest

// CreatePRRequest is an alias for response.CreatePRRequest
type CreatePRRequest = response.CreatePRRequest

// MergePRRequest is an alias for response.MergePRRequest
type MergePRRequest = response.MergePRRequest

// UpdatePRRequest is an alias for response.UpdatePRRequest
type UpdatePRRequest = response.UpdatePRRequest

// ListPRs lists pull requests in a repository
func (c *Client) ListPRs(owner, repo, state string) ([]PullRequest, error) {
	var prs []PullRequest
	err := c.DoFromEndpoint(PRs.List, []interface{}{owner, repo}, nil, &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

// GetPR gets a single pull request
func (c *Client) GetPR(owner, repo string, number int) (*PullRequest, error) {
	var pr PullRequest
	err := c.DoFromEndpoint(PRs.List, []interface{}{owner, repo, number}, nil, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// CreatePR creates a new pull request
func (c *Client) CreatePR(owner, repo, title, body, head, base string) (*PullRequest, error) {
	req := CreatePRRequest{Title: title, Body: body, Head: head, Base: base}
	var pr PullRequest
	err := c.DoFromEndpoint(PRs.Create, []interface{}{owner, repo}, req, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// MergePR merges a pull request
func (c *Client) MergePR(owner, repo string, number int, req MergePRRequest) error {
	return c.DoFromEndpoint(PRs.Merge, []interface{}{owner, repo, number}, req, nil)
}

// UpdatePR updates a pull request
func (c *Client) UpdatePR(owner, repo string, number int, req UpdatePRRequest) (*PullRequest, error) {
	var pr PullRequest
	err := c.DoFromEndpoint(PRs.Update, []interface{}{owner, repo, number}, req, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// UpdatePRState updates a pull request's state (open/closed)
func (c *Client) UpdatePRState(owner, repo string, number int, state PRState) (*PullRequest, error) {
	updateReq := map[string]string{"state": string(state)}
	var pr PullRequest
	err := c.DoFromEndpoint(PRs.Update, []interface{}{owner, repo, number}, updateReq, &pr)
	if err != nil {
		return nil, err
	}
	return &pr, nil
}

// CreatePRComment adds a comment to a pull request
func (c *Client) CreatePRComment(owner, repo string, number int, body string) error {
	commentReq := map[string]string{"body": body}
	return c.DoFromEndpoint(PRs.Comment, []interface{}{owner, repo, number}, commentReq, nil)
}
