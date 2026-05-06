package response

// License represents a license template
type License struct {
	License string `json:"license"`
	Source  string `json:"source"`
}

// GitignoreTemplate represents a gitignore template
type GitignoreTemplate struct {
	Name   string `json:"name"`
	Source string `json:"source"`
}

// MarkdownRenderRequest is the request body for rendering Markdown
type MarkdownRenderRequest struct {
	Text string `json:"text"`
}
