package api

import (
	"strconv"

	"github.com/ricsy/gt/pkg/api/response"
	"github.com/ricsy/gt/pkg/util"
)

// Webhook is an alias for response.Webhook
type Webhook = response.Webhook

// CreateWebhookOptions is an alias for response.CreateWebhookOptions
type CreateWebhookOptions = response.CreateWebhookOptions

// UpdateWebhookOptions is an alias for response.UpdateWebhookOptions
type UpdateWebhookOptions = response.UpdateWebhookOptions

// ListWebhooksOptions contains optional parameters for ListWebhooks
type ListWebhooksOptions struct {
	Page    int
	PerPage int
}

// ListWebhooks lists webhooks for a repository
func (c *Client) ListWebhooks(owner, repo string, opts ListWebhooksOptions) ([]Webhook, error) {
	path := Webhooks.List.Build(owner, repo)
	query := buildWebhookQuery(opts)
	if query != "" {
		path += "?" + query
	}
	var webhooks []Webhook
	err := c.Do("GET", path, nil, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

func buildWebhookQuery(opts ListWebhooksOptions) string {
	params := []string{}
	if opts.Page > 0 {
		params = append(params, "page", strconv.Itoa(opts.Page))
	}
	if opts.PerPage > 0 {
		params = append(params, "per_page", strconv.Itoa(opts.PerPage))
	}
	return util.BuildQuery(params...)
}

// GetWebhook gets a single webhook
func (c *Client) GetWebhook(owner, repo string, id int64) (*Webhook, error) {
	var webhook Webhook
	err := c.DoFromEndpoint(Webhooks.Get, []interface{}{owner, repo, id}, nil, &webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// CreateWebhook creates a new webhook
func (c *Client) CreateWebhook(owner, repo string, opts CreateWebhookOptions) (*Webhook, error) {
	var webhook Webhook
	err := c.DoFromEndpoint(Webhooks.Create, []interface{}{owner, repo}, opts, &webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// UpdateWebhook updates a webhook
func (c *Client) UpdateWebhook(owner, repo string, id int64, opts UpdateWebhookOptions) (*Webhook, error) {
	var webhook Webhook
	err := c.DoFromEndpoint(Webhooks.Update, []interface{}{owner, repo, id}, opts, &webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// DeleteWebhook deletes a webhook
func (c *Client) DeleteWebhook(owner, repo string, id int64) error {
	return c.DoFromEndpoint(Webhooks.Delete, []interface{}{owner, repo, id}, nil, nil)
}

// TestWebhook tests a webhook
func (c *Client) TestWebhook(owner, repo string, id int64) error {
	return c.DoFromEndpoint(Webhooks.Test, []interface{}{owner, repo, id}, nil, nil)
}
