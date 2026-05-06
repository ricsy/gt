package api

import (
	"github.com/ricsy/gt/pkg/api/response"
)

// ListLicenses lists available licenses
func (c *Client) ListLicenses() ([]string, error) {
	var licenses []string
	err := c.DoFromEndpoint(Miscellaneous.ListLicenses, nil, nil, &licenses)
	if err != nil {
		return nil, err
	}
	return licenses, nil
}

// GetLicense gets a license by key
func (c *Client) GetLicense(license string) (*response.License, error) {
	var l response.License
	err := c.DoFromEndpoint(Miscellaneous.GetLicense, []interface{}{license}, nil, &l)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

// GetLicenseRaw gets a license raw content
func (c *Client) GetLicenseRaw(license string) (string, error) {
	var result string
	err := c.DoFromEndpoint(Miscellaneous.GetLicenseRaw, []interface{}{license}, nil, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// GetRepoLicense gets the license key for a repo
func (c *Client) GetRepoLicense(owner, repo string) (string, error) {
	var result map[string]string
	err := c.DoFromEndpoint(Miscellaneous.GetRepoLicense, []interface{}{owner, repo}, nil, &result)
	if err != nil {
		return "", err
	}
	return result["license"], nil
}

// ListGitignoreTemplates lists available gitignore templates
func (c *Client) ListGitignoreTemplates() ([]string, error) {
	var templates []string
	err := c.DoFromEndpoint(Miscellaneous.ListGitignoreTemplates, nil, nil, &templates)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// GetGitignoreTemplate gets a gitignore template by name
func (c *Client) GetGitignoreTemplate(name string) (*response.GitignoreTemplate, error) {
	var t response.GitignoreTemplate
	err := c.DoFromEndpoint(Miscellaneous.GetGitignoreTemplate, []interface{}{name}, nil, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

// GetGitignoreTemplateRaw gets a gitignore template raw content
func (c *Client) GetGitignoreTemplateRaw(name string) (string, error) {
	var result string
	err := c.DoFromEndpoint(Miscellaneous.GetGitignoreTemplateRaw, []interface{}{name}, nil, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}

// RenderMarkdown renders Markdown text
func (c *Client) RenderMarkdown(text string) (string, error) {
	var result string
	err := c.DoFromEndpoint(Miscellaneous.RenderMarkdown, nil, response.MarkdownRenderRequest{Text: text}, &result)
	if err != nil {
		return "", err
	}
	return result, nil
}
