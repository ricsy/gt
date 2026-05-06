package response

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
