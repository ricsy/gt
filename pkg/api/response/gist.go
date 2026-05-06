package response

import "time"

// Gist represents a Gitee gist
type Gist struct {
	URL         string     `json:"url"`
	ForksURL    string     `json:"forks_url"`
	CommitsURL  string     `json:"commits_url"`
	ID          string     `json:"id"`
	Description string     `json:"description"`
	Public      bool       `json:"public"`
	Owner       *UserBasic `json:"owner"`
	User        *UserBasic `json:"user"`
	Files       any        `json:"files"`
	Truncated   bool       `json:"truncated"`
	HTMLURL     string     `json:"html_url"`
	Comments    int        `json:"comments"`
	CommentsURL string     `json:"comments_url"`
	GitPullURL  string     `json:"git_pull_url"`
	GitPushURL  string     `json:"git_push_url"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Forks       any        `json:"forks"`
	History     any        `json:"history"`
}

// GistComment represents a comment on a gist
type GistComment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateGistOptions is the request body for creating a gist
type CreateGistOptions struct {
	Files       any    `json:"files"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

// UpdateGistOptions is the request body for updating a gist
type UpdateGistOptions struct {
	Files       any    `json:"files,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateGistCommentOptions is the request body for creating a gist comment
type CreateGistCommentOptions struct {
	Body string `json:"body"`
}

// UpdateGistCommentOptions is the request body for updating a gist comment
type UpdateGistCommentOptions struct {
	Body string `json:"body"`
}
