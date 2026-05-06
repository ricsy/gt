package api

// Webhook represents a Gitee webhook
type Webhook struct {
	ID                  int64  `json:"id"`
	URL                 string `json:"url"`
	Password            string `json:"password"`
	Result              string `json:"result"`
	ProjectID           int64  `json:"project_id"`
	ResultCode          int    `json:"result_code"`
	CreatedAt           string `json:"created_at"`
	PushEvents          bool   `json:"push_events"`
	TagPushEvents       bool   `json:"tag_push_events"`
	IssuesEvents        bool   `json:"issues_events"`
	NoteEvents          bool   `json:"note_events"`
	MergeRequestsEvents bool   `json:"merge_requests_events"`
	Title               string `json:"title"`
}

// ListWebhooks lists webhooks for a repository
func (c *Client) ListWebhooks(owner, repo string) ([]Webhook, error) {
	var webhooks []Webhook
	err := c.DoFromEndpoint(Webhooks.List, []interface{}{owner, repo}, nil, &webhooks)
	if err != nil {
		return nil, err
	}
	return webhooks, nil
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
	err := c.DoFromEndpoint(Webhooks.Create, []interface{}{owner, repo}, opts, &webhook)
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
