package api

import (
	"github.com/ricsy/gt/pkg/api/response"
)

// ListRepoNotificationsOptions is an alias for response.ListRepoNotificationsOptions
type ListRepoNotificationsOptions = response.ListRepoNotificationsOptions

// ListNotificationsOptions is an alias for response.ListNotificationsOptions
type ListNotificationsOptions = response.ListNotificationsOptions

// ListMessagesOptions is an alias for response.ListMessagesOptions
type ListMessagesOptions = response.ListMessagesOptions

// MarkNotificationsReadOptions is an alias for response.MarkNotificationsReadOptions
type MarkNotificationsReadOptions = response.MarkNotificationsReadOptions

// CreateMessageOptions is an alias for response.CreateMessageOptions
type CreateMessageOptions = response.CreateMessageOptions

// ListRepoNotifications lists notifications for a repo
func (c *Client) ListRepoNotifications(owner, repo string, opts response.ListRepoNotificationsOptions) (*response.UserNotificationList, error) {
	var result response.UserNotificationList
	err := c.DoFromEndpoint(RepoNotifications.List, []interface{}{owner, repo}, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MarkRepoNotificationsRead marks repo notifications as read
func (c *Client) MarkRepoNotificationsRead(owner, repo string, opts response.MarkNotificationsReadOptions) error {
	return c.DoFromEndpoint(RepoNotifications.Update, []interface{}{owner, repo}, opts, nil)
}

// ListNotifications lists all notifications for the user
func (c *Client) ListNotifications(opts response.ListNotificationsOptions) (*response.UserNotificationList, error) {
	var result response.UserNotificationList
	err := c.DoFromEndpoint(NotificationThreads.List, nil, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MarkAllNotificationsRead marks all notifications as read
func (c *Client) MarkAllNotificationsRead(opts response.MarkNotificationsReadOptions) error {
	return c.DoFromEndpoint(NotificationThreads.Update, nil, opts, nil)
}

// GetNotification gets a single notification
func (c *Client) GetNotification(id string) (*response.UserNotification, error) {
	var result response.UserNotification
	err := c.DoFromEndpoint(NotificationThreads.Get, []interface{}{id}, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MarkNotificationRead marks a single notification as read
func (c *Client) MarkNotificationRead(id string) error {
	return c.DoFromEndpoint(NotificationThreads.Patch, []interface{}{id}, nil, nil)
}

// GetNotificationCount gets notification count
func (c *Client) GetNotificationCount(unread *bool) (*response.UserNotificationCount, error) {
	var result response.UserNotificationCount
	err := c.DoFromEndpoint(NotificationCount.Get, nil, unread, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListMessages lists messages for the user
func (c *Client) ListMessages(opts response.ListMessagesOptions) (*response.UserMessageList, error) {
	var result response.UserMessageList
	err := c.DoFromEndpoint(Messages.List, nil, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateMessage creates a new message
func (c *Client) CreateMessage(opts response.CreateMessageOptions) (*response.UserMessage, error) {
	var result response.UserMessage
	err := c.DoFromEndpoint(Messages.Create, nil, opts, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MarkAllMessagesRead marks all messages as read
func (c *Client) MarkAllMessagesRead() error {
	return c.DoFromEndpoint(Messages.Update, nil, nil, nil)
}

// GetMessage gets a single message
func (c *Client) GetMessage(id string) (*response.UserMessage, error) {
	var result response.UserMessage
	err := c.DoFromEndpoint(Messages.Get, []interface{}{id}, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// MarkMessageRead marks a message as read
func (c *Client) MarkMessageRead(id string) error {
	return c.DoFromEndpoint(Messages.Patch, []interface{}{id}, nil, nil)
}
