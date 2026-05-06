package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// CreateCheckRunOptions is an alias for response.CreateCheckRunOptions
type CreateCheckRunOptions = response.CreateCheckRunOptions

// UpdateCheckRunOptions is an alias for response.UpdateCheckRunOptions
type UpdateCheckRunOptions = response.UpdateCheckRunOptions

// ListCheckRunsOptions is an alias for response.ListCheckRunsOptions
type ListCheckRunsOptions = response.ListCheckRunsOptions

// CreateCheckRun creates a new check run
func (c *Client) CreateCheckRun(owner, repo string, opts CreateCheckRunOptions) (*response.CheckRun, error) {
	var checkRun response.CheckRun
	err := c.DoFromEndpoint(CheckRuns.Create, []interface{}{owner, repo}, opts, &checkRun)
	if err != nil {
		return nil, err
	}
	return &checkRun, nil
}

// GetCheckRun gets a single check run
func (c *Client) GetCheckRun(owner, repo string, checkRunID int64) (*response.CheckRun, error) {
	var checkRun response.CheckRun
	err := c.DoFromEndpoint(CheckRuns.Get, []interface{}{owner, repo, checkRunID}, nil, &checkRun)
	if err != nil {
		return nil, err
	}
	return &checkRun, nil
}

// UpdateCheckRun updates a check run
func (c *Client) UpdateCheckRun(owner, repo string, checkRunID int64, opts UpdateCheckRunOptions) (*response.CheckRun, error) {
	var checkRun response.CheckRun
	err := c.DoFromEndpoint(CheckRuns.Update, []interface{}{owner, repo, checkRunID}, opts, &checkRun)
	if err != nil {
		return nil, err
	}
	return &checkRun, nil
}

// GetCheckRunAnnotations gets annotations for a check run
func (c *Client) GetCheckRunAnnotations(owner, repo string, checkRunID int64, opts ListCheckRunsOptions) ([]response.CheckAnnotation, error) {
	var annotations []response.CheckAnnotation
	path := CheckRuns.GetAnnotations.Build(owner, repo, checkRunID)
	query := buildCheckRunsQuery(opts)
	err := c.Do("GET", path+query, nil, &annotations)
	if err != nil {
		return nil, err
	}
	return annotations, nil
}

// ListCommitCheckRuns lists check runs for a commit
func (c *Client) ListCommitCheckRuns(owner, repo, ref string, opts ListCheckRunsOptions) ([]response.CheckRun, error) {
	var checkRuns []response.CheckRun
	path := CheckRuns.GetCommitCheckRuns.Build(owner, repo, ref)
	query := buildCheckRunsQuery(opts)
	err := c.Do("GET", path+query, nil, &checkRuns)
	if err != nil {
		return nil, err
	}
	return checkRuns, nil
}

func buildCheckRunsQuery(opts ListCheckRunsOptions) string {
	var params []string
	if opts.Page > 0 {
		params = append(params, "page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params = append(params, "per_page", strconv.Itoa(opts.PerPage))
	}
	if len(params) == 0 {
		return ""
	}
	return "?" + util.BuildQuery(params...)
}
