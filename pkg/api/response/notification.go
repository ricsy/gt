package response

// UserNotificationSubject represents the subject of a notification
type UserNotificationSubject struct {
	Title            string `json:"title"`
	URL              string `json:"url"`
	LatestCommentURL string `json:"latest_comment_url"`
	Type             string `json:"type"`
}

// UserNotificationNamespace represents a namespace in a notification
type UserNotificationNamespace struct {
	Name    string `json:"name"`
	HTMLURL string `json:"html_url"`
	Type    string `json:"type"`
}

// UserNotification represents a Gitee notification
type UserNotification struct {
	ID         int64                       `json:"id"`
	Content    string                      `json:"content"`
	Type       string                      `json:"type"`
	Unread     bool                        `json:"unread"`
	Mute       bool                        `json:"mute"`
	UpdatedAt  string                      `json:"updated_at"`
	URL        string                      `json:"url"`
	HTMLURL    string                      `json:"html_url"`
	Actor      UserBasic                   `json:"actor"`
	Repository ProjectBasic                `json:"repository"`
	Subject    UserNotificationSubject     `json:"subject"`
	Namespaces []UserNotificationNamespace `json:"namespaces"`
}

// UserNotificationList represents a list of notifications
type UserNotificationList struct {
	TotalCount int64              `json:"total_count"`
	List       []UserNotification `json:"list"`
}

// UserNotificationCount represents notification count
type UserNotificationCount struct {
	TotalCount        int64 `json:"total_count"`
	NotificationCount int64 `json:"notification_count"`
	MessageCount      int64 `json:"message_count"`
}

// UserMessage represents a Gitee message
type UserMessage struct {
	ID        int64     `json:"id"`
	Sender    UserBasic `json:"sender"`
	Unread    bool      `json:"unread"`
	Content   string    `json:"content"`
	UpdatedAt string    `json:"updated_at"`
	URL       string    `json:"url"`
	HTMLURL   string    `json:"html_url"`
}

// UserMessageList represents a list of messages
type UserMessageList struct {
	TotalCount int64         `json:"total_count"`
	List       []UserMessage `json:"list"`
}

// ListRepoNotificationsOptions contains optional parameters for listing repo notifications
type ListRepoNotificationsOptions struct {
	Unread        *bool  `json:"unread,omitempty"`
	Participating *bool  `json:"participating,omitempty"`
	Type          string `json:"type,omitempty"`
	Since         string `json:"since,omitempty"`
	Before        string `json:"before,omitempty"`
	IDs           string `json:"ids,omitempty"`
	Page          int    `json:"page,omitempty"`
	PerPage       int    `json:"per_page,omitempty"`
}

// ListNotificationsOptions contains optional parameters for listing notifications
type ListNotificationsOptions struct {
	Unread        *bool  `json:"unread,omitempty"`
	Participating *bool  `json:"participating,omitempty"`
	Type          string `json:"type,omitempty"`
	Since         string `json:"since,omitempty"`
	Before        string `json:"before,omitempty"`
	IDs           string `json:"ids,omitempty"`
	Page          int    `json:"page,omitempty"`
	PerPage       int    `json:"per_page,omitempty"`
}

// ListMessagesOptions contains optional parameters for listing messages
type ListMessagesOptions struct {
	Unread  *bool  `json:"unread,omitempty"`
	Since   string `json:"since,omitempty"`
	Before  string `json:"before,omitempty"`
	IDs     string `json:"ids,omitempty"`
	Page    int    `json:"page,omitempty"`
	PerPage int    `json:"per_page,omitempty"`
}

// MarkNotificationsReadOptions contains options for marking notifications as read
type MarkNotificationsReadOptions struct {
	IDs string `json:"ids,omitempty"`
}

// CreateMessageOptions contains options for creating a message
type CreateMessageOptions struct {
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}
