package api

import "github.com/ricsy/gt/pkg/api/response"

// CreateGistOptions is an alias for response.CreateGistOptions
type CreateGistOptions = response.CreateGistOptions

// UpdateGistOptions is an alias for response.UpdateGistOptions
type UpdateGistOptions = response.UpdateGistOptions

// CreateGistCommentOptions is an alias for response.CreateGistCommentOptions
type CreateGistCommentOptions = response.CreateGistCommentOptions

// UpdateGistCommentOptions is an alias for response.UpdateGistCommentOptions
type UpdateGistCommentOptions = response.UpdateGistCommentOptions

// ListGists lists gists for the authenticated user
func (c *Client) ListGists(opts *ListGistsOptions) ([]response.Gist, error) {
	var gists []response.Gist
	err := c.DoFromEndpoint(Gists.List, nil, opts, &gists)
	if err != nil {
		return nil, err
	}
	return gists, nil
}

// CreateGist creates a new gist
func (c *Client) CreateGist(opts response.CreateGistOptions) (*response.Gist, error) {
	var gist response.Gist
	err := c.DoFromEndpoint(Gists.Create, nil, opts, &gist)
	if err != nil {
		return nil, err
	}
	return &gist, nil
}

// GetGist gets a single gist
func (c *Client) GetGist(id string) (*response.Gist, error) {
	var gist response.Gist
	err := c.DoFromEndpoint(Gists.Get, []interface{}{id}, nil, &gist)
	if err != nil {
		return nil, err
	}
	return &gist, nil
}

// UpdateGist updates a gist
func (c *Client) UpdateGist(id string, opts response.UpdateGistOptions) (*response.Gist, error) {
	var gist response.Gist
	err := c.DoFromEndpoint(Gists.Update, []interface{}{id}, opts, &gist)
	if err != nil {
		return nil, err
	}
	return &gist, nil
}

// DeleteGist deletes a gist
func (c *Client) DeleteGist(id string) error {
	return c.DoFromEndpoint(Gists.Delete, []interface{}{id}, nil, nil)
}

// ListStarredGists lists starred gists for the authenticated user
func (c *Client) ListStarredGists(opts *ListGistsOptions) ([]response.Gist, error) {
	var gists []response.Gist
	err := c.DoFromEndpoint(Gists.Starred, nil, opts, &gists)
	if err != nil {
		return nil, err
	}
	return gists, nil
}

// IsGistStarred checks if a gist is starred
func (c *Client) IsGistStarred(id string) (bool, error) {
	err := c.DoFromEndpoint(Gists.Star, []interface{}{id}, nil, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// StarGist stars a gist
func (c *Client) StarGist(id string) error {
	return c.DoFromEndpoint(Gists.StarPut, []interface{}{id}, nil, nil)
}

// UnstarGist unstars a gist
func (c *Client) UnstarGist(id string) error {
	return c.DoFromEndpoint(Gists.StarDel, []interface{}{id}, nil, nil)
}

// ListGistForks lists forks of a gist
func (c *Client) ListGistForks(id string) ([]response.Gist, error) {
	var gists []response.Gist
	err := c.DoFromEndpoint(Gists.Forks, []interface{}{id}, nil, &gists)
	if err != nil {
		return nil, err
	}
	return gists, nil
}

// ForkGist forks a gist
func (c *Client) ForkGist(id string) (*response.Gist, error) {
	var gist response.Gist
	err := c.DoFromEndpoint(Gists.Fork, []interface{}{id}, nil, &gist)
	if err != nil {
		return nil, err
	}
	return &gist, nil
}

// ListGistCommits lists commits of a gist
func (c *Client) ListGistCommits(id string) ([]response.Gist, error) {
	var gists []response.Gist
	err := c.DoFromEndpoint(Gists.Commits, []interface{}{id}, nil, &gists)
	if err != nil {
		return nil, err
	}
	return gists, nil
}

// ListGistComments lists comments on a gist
func (c *Client) ListGistComments(gistID string) ([]response.GistComment, error) {
	var comments []response.GistComment
	err := c.DoFromEndpoint(Gists.Comments, []interface{}{gistID}, nil, &comments)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

// GetGistComment gets a single comment on a gist
func (c *Client) GetGistComment(gistID string, commentID int64) (*response.GistComment, error) {
	var comment response.GistComment
	err := c.DoFromEndpoint(Gists.Comment, []interface{}{gistID, commentID}, nil, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// CreateGistComment creates a comment on a gist
func (c *Client) CreateGistComment(gistID string, opts response.CreateGistCommentOptions) (*response.GistComment, error) {
	var comment response.GistComment
	err := c.DoFromEndpoint(Gists.CreateComment, []interface{}{gistID}, opts, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// UpdateGistComment updates a comment on a gist
func (c *Client) UpdateGistComment(gistID string, commentID int64, opts response.UpdateGistCommentOptions) (*response.GistComment, error) {
	var comment response.GistComment
	err := c.DoFromEndpoint(Gists.UpdateComment, []interface{}{gistID, commentID}, opts, &comment)
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

// DeleteGistComment deletes a comment from a gist
func (c *Client) DeleteGistComment(gistID string, commentID int64) error {
	return c.DoFromEndpoint(Gists.DeleteComment, []interface{}{gistID, commentID}, nil, nil)
}

// ListGistsOptions contains the optional parameters for ListGists
type ListGistsOptions struct {
	Since   string `json:"since,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}
