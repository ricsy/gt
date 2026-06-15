package api

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

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

// ReviewPRRequest is an alias for response.ReviewPRRequest
type ReviewPRRequest = response.ReviewPRRequest

// TestPRRequest is an alias for response.TestPRRequest
type TestPRRequest = response.TestPRRequest

// UpdatePRRequest is an alias for response.UpdatePRRequest
type UpdatePRRequest = response.UpdatePRRequest

// ListPRsOptions is an alias for response.ListPRsOptions
type ListPRsOptions = response.ListPRsOptions

// ListPRs lists pull requests in a repository
func (c *Client) ListPRs(owner, repo string, opts ListPRsOptions) ([]PullRequest, error) {
	path := PRs.List.Build(owner, repo)
	query := buildPRQuery(opts)
	if query != "" {
		path += "?" + query
	}
	var prs []PullRequest
	err := c.Do("GET", path, nil, &prs)
	if err != nil {
		return nil, err
	}
	return prs, nil
}

// buildPRQuery builds query string from ListPRsOptions
func buildPRQuery(opts ListPRsOptions) string {
	params := []string{
		"state", opts.State,
		"head", opts.Head,
		"base", opts.Base,
		"sort", opts.Sort,
		"direction", opts.Direction,
		"labels", opts.Labels,
		"author", opts.Author,
		"assignee", opts.Assignee,
		"tester", opts.Tester,
	}
	if opts.Since != "" {
		params = append(params, "since", opts.Since)
	}
	if opts.MilestoneNumber > 0 {
		params = append(params, "milestone_number", fmt.Sprintf("%d", opts.MilestoneNumber))
	}
	params = append(params, paginationParams(opts.Page, opts.PerPage)...)
	return util.BuildQuery(params...)
}

// GetPR gets a single pull request
func (c *Client) GetPR(owner, repo string, number int) (*PullRequest, error) {
	var pr PullRequest
	err := c.DoFromEndpoint(PRs.Get, []interface{}{owner, repo, number}, nil, &pr)
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

// ReviewPR marks a pull request review as passed
func (c *Client) ReviewPR(owner, repo string, number int, req ReviewPRRequest) error {
	return c.DoFromEndpoint(PRs.Review, []interface{}{owner, repo, number}, req, nil)
}

// TestPR marks a pull request test as passed
func (c *Client) TestPR(owner, repo string, number int, req TestPRRequest) error {
	return c.DoFromEndpoint(PRs.Test, []interface{}{owner, repo, number}, req, nil)
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

// ListPRComments lists comments on a pull request
func (c *Client) ListPRComments(owner, repo string, number int) ([]response.PullRequestComment, error) {
	var comments []response.PullRequestComment
	err := c.DoFromEndpoint(PRComments.List, []interface{}{owner, repo, number}, nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetPRComment gets a single comment on a pull request
func (c *Client) GetPRComment(owner, repo string, number, commentID int) (*response.PullRequestComment, error) {
	var comment response.PullRequestComment
	err := c.DoFromEndpoint(PRComments.Get, []interface{}{owner, repo, commentID}, nil, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// ListPRCommits lists commits in a pull request
func (c *Client) ListPRCommits(owner, repo string, number int) ([]response.PullRequestCommit, error) {
	var commits []response.PullRequestCommit
	err := c.DoFromEndpoint(PRCommits.List, []interface{}{owner, repo, number}, nil, &commits)
	if err != nil {
		return nil, err
	}
	return commits, nil
}

// ListPRFiles lists changed files in a pull request
func (c *Client) ListPRFiles(owner, repo string, number int) ([]response.PullRequestFile, error) {
	var files []response.PullRequestFile
	err := c.DoFromEndpoint(PRFiles.List, []interface{}{owner, repo, number}, nil, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}
