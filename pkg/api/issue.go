package api

import (
	"fmt"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// IssueState is an alias for response.IssueState
type IssueState = response.IssueState

const (
	IssueStateOpen        = response.IssueStateOpen
	IssueStateClosed      = response.IssueStateClosed
	IssueStateProgressing = response.IssueStateProgressing
	IssueStateRejected    = response.IssueStateRejected
	IssueStateAll         = response.IssueStateAll
)

// Issue is an alias for response.Issue
type Issue = response.Issue

// Label is an alias for response.Label
type Label = response.Label

// User is an alias for response.UserBasic
type User = response.UserBasic

// IssueComment is an alias for response.IssueComment
type IssueComment = response.IssueComment

// Milestone is an alias for response.Milestone
type Milestone = response.Milestone

// Project is an alias for response.Project
type Project = response.Project

// ListIssuesOptions is an alias for response.ListIssuesOptions
type ListIssuesOptions = response.ListIssuesOptions

// CreateIssueRequest is an alias for response.CreateIssueRequest
type CreateIssueRequest = response.CreateIssueRequest

// UpdateIssueRequest is an alias for response.UpdateIssueRequest
type UpdateIssueRequest = response.UpdateIssueRequest

// ListIssues lists issues in a repository
func (c *Client) ListIssues(owner, repo string, opts ListIssuesOptions) ([]Issue, error) {
	path := Issues.List.Build(owner, repo)
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
	if opts.Since != "" {
		params = append(params, "since", opts.Since)
	}
	if opts.Schedule != "" {
		params = append(params, "schedule", opts.Schedule)
	}
	if opts.Deadline != "" {
		params = append(params, "deadline", opts.Deadline)
	}
	if opts.CreatedAt != "" {
		params = append(params, "created_at", opts.CreatedAt)
	}
	if opts.FinishedAt != "" {
		params = append(params, "finished_at", opts.FinishedAt)
	}
	return util.BuildQuery(params...)
}

// GetIssue gets a single issue
func (c *Client) GetIssue(owner, repo, number string) (*Issue, error) {
	var issue Issue
	err := c.DoFromEndpoint(Issues.List, []interface{}{owner, repo, number}, nil, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(owner, repo, title, body string) (*Issue, error) {
	req := CreateIssueRequest{Title: title, Body: body, Repo: repo}
	var issue Issue
	err := c.DoFromEndpoint(Issues.Create, []interface{}{owner}, req, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// UpdateIssue updates an issue
func (c *Client) UpdateIssue(owner, repo, number string, req UpdateIssueRequest) (*Issue, error) {
	var issue Issue
	err := c.DoFromEndpoint(Issues.Update, []interface{}{owner, number}, req, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// UpdateIssueState updates an issue's state (open/closed/progressing)
// Note: rejected state is not supported by the API for updates
func (c *Client) UpdateIssueState(owner, repo, number string, state IssueState) (*Issue, error) {
	req := map[string]string{"state": string(state), "repo": repo}
	var issue Issue
	err := c.DoFromEndpoint(Issues.Update, []interface{}{owner, number}, req, &issue)
	if err != nil {
		return nil, err
	}
	return &issue, nil
}

// ListIssueComments lists comments for an issue
func (c *Client) ListIssueComments(owner, repo, number string) ([]IssueComment, error) {
	var comments []IssueComment
	err := c.DoFromEndpoint(Issues.List, []interface{}{owner, repo, number, "comments"}, nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// CreateIssueComment creates a comment on an issue
func (c *Client) CreateIssueComment(owner, repo, number, body string) (*IssueComment, error) {
	req := map[string]string{"body": body}
	var comment IssueComment
	err := c.DoFromEndpoint(Issues.List, []interface{}{owner, repo, number, "comments"}, req, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
