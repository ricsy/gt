package response

// Blob represents a Git blob object.
type Blob struct {
	SHA      string `json:"sha"`
	Size     int    `json:"size"`
	URL      string `json:"url"`
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
}

// Tree represents a Git tree object.
type Tree struct {
	SHA       string      `json:"sha"`
	URL       string      `json:"url"`
	Tree      []TreeEntry `json:"tree"`
	Truncated bool        `json:"truncated"`
}

// TreeEntry represents an item in a Git tree.
type TreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Size int    `json:"size"`
	SHA  string `json:"sha"`
	URL  string `json:"url"`
}

// GetTreeOptions contains options for fetching a Git tree.
type GetTreeOptions struct {
	Recursive int `json:"recursive,omitempty"`
}

// GiteeMetrics represents the Gitee repository metric response.
type GiteeMetrics struct {
	Data       string       `json:"data"`
	TotalScore int          `json:"total_score"`
	CreatedAt  string       `json:"created_at"`
	Repo       ProjectBasic `json:"repo"`
}
