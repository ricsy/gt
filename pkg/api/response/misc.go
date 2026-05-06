package response

// License represents a license template
type License struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	SPDXID  string `json:"spdx_id"`
	URL     string `json:"url"`
	HTMLURL string `json:"html_url"`
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
