package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// Blob is an alias for response.Blob.
type Blob = response.Blob

// Tree is an alias for response.Tree.
type Tree = response.Tree

// GetTreeOptions is an alias for response.GetTreeOptions.
type GetTreeOptions = response.GetTreeOptions

// GiteeMetrics is an alias for response.GiteeMetrics.
type GiteeMetrics = response.GiteeMetrics

// GetBlob gets a Git blob by SHA.
func (c *Client) GetBlob(owner, repo, sha string) (*Blob, error) {
	var blob Blob
	err := c.DoFromEndpoint(GitData.Blob, []interface{}{owner, repo, sha}, nil, &blob)
	if err != nil {
		return nil, err
	}
	return &blob, nil
}

// GetTree gets a Git tree by SHA, branch, or commit.
func (c *Client) GetTree(owner, repo, sha string, opts GetTreeOptions) (*Tree, error) {
	var tree Tree
	path := GitData.Tree.Build(owner, repo, sha)
	if opts.Recursive > 0 {
		path += "?" + util.BuildQuery("recursive", strconv.Itoa(opts.Recursive))
	}
	err := c.Do("GET", path, nil, &tree)
	if err != nil {
		return nil, err
	}
	return &tree, nil
}

// GetGiteeMetrics gets Gitee metrics for a repository.
func (c *Client) GetGiteeMetrics(owner, repo string) (*GiteeMetrics, error) {
	var metrics GiteeMetrics
	err := c.DoFromEndpoint(GitData.Metrics, []interface{}{owner, repo}, nil, &metrics)
	if err != nil {
		return nil, err
	}
	return &metrics, nil
}
