package api

import "fmt"

// Webhook represents a Gitee webhook
type Webhook struct {
	ID         int64  `json:"id"`
	URL        string `json:"url"`
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
	CREATEDAT  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Title      string `json:"title"`
	Type       string `json:"type"`
	LastTestAt string `json:"last_test_at"`
}

// ListWebhooks lists webhooks for a repository
func (c *Client) ListWebhooks(owner, repo string) ([]Webhook, error) {
	var webhooks []Webhook
	path := fmt.Sprintf(apiPathWebhooks, owner, repo)
	err := c.Do("GET", path, nil, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
}

// GetWebhook gets a single webhook
func (c *Client) GetWebhook(owner, repo string, id int64) (*Webhook, error) {
	var webhook Webhook
	path := fmt.Sprintf(apiPathWebhooks+"/%d", owner, repo, id)
	err := c.Do("GET", path, nil, &webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// CreateWebhookOptions contains options for creating a webhook
type CreateWebhookOptions struct {
	URL                 string `json:"url"`
	Title               string `json:"title,omitempty"`
	EncryptionType      int    `json:"encryption_type,omitempty"`
	Password            string `json:"password,omitempty"`
	PushEvents          *bool  `json:"push_events,omitempty"`
	TagPushEvents       *bool  `json:"tag_push_events,omitempty"`
	IssuesEvents        *bool  `json:"issues_events,omitempty"`
	NoteEvents          *bool  `json:"note_events,omitempty"`
	MergeRequestsEvents *bool  `json:"merge_requests_events,omitempty"`
}

// CreateWebhook creates a new webhook
func (c *Client) CreateWebhook(owner, repo string, opts CreateWebhookOptions) (*Webhook, error) {
	var webhook Webhook
	path := fmt.Sprintf(apiPathWebhooks, owner, repo)
	err := c.Do("POST", path, opts, &webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// UpdateWebhookOptions contains options for updating a webhook
type UpdateWebhookOptions struct {
	URL                 string `json:"url"`
	Title               string `json:"title,omitempty"`
	EncryptionType      int    `json:"encryption_type,omitempty"`
	Password            string `json:"password,omitempty"`
	PushEvents          *bool  `json:"push_events,omitempty"`
	TagPushEvents       *bool  `json:"tag_push_events,omitempty"`
	IssuesEvents        *bool  `json:"issues_events,omitempty"`
	NoteEvents          *bool  `json:"note_events,omitempty"`
	MergeRequestsEvents *bool  `json:"merge_requests_events,omitempty"`
}

// UpdateWebhook updates a webhook
func (c *Client) UpdateWebhook(owner, repo string, id int64, opts UpdateWebhookOptions) (*Webhook, error) {
	var webhook Webhook
	path := fmt.Sprintf(apiPathWebhooks+"/%d", owner, repo, id)
	err := c.Do("PATCH", path, opts, &webhook)
	if err != nil {
		return nil, err
	}
	return &webhook, nil
}

// DeleteWebhook deletes a webhook
func (c *Client) DeleteWebhook(owner, repo string, id int64) error {
	path := fmt.Sprintf(apiPathWebhooks+"/%d", owner, repo, id)
	return c.Do("DELETE", path, nil, nil)
}

// TestWebhook tests a webhook
func (c *Client) TestWebhook(owner, repo string, id int64) error {
	path := fmt.Sprintf(apiPathWebhooks+"/%d/tests", owner, repo, id)
	return c.Do("POST", path, nil, nil)
}
