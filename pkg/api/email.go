package api

import (
	"github.com/ricsy/gt/pkg/api/response"
)

// ListEmails lists all emails for the authenticated user
func (c *Client) ListEmails() ([]response.Email, error) {
	var emails []response.Email
	err := c.DoFromEndpoint(Emails.List, nil, nil, &emails)
	if err != nil {
		return nil, err
	}
	return emails, nil
}
