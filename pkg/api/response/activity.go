package response

import "encoding/json"

// Event represents a Gitee activity event.
type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Actor     UserBasic       `json:"actor"`
	Repo      ProjectBasic    `json:"repo"`
	Org       Org             `json:"org"`
	Public    bool            `json:"public"`
	CreatedAt string          `json:"created_at"`
	Payload   json.RawMessage `json:"payload"`
}

// ListActivityOptions contains pagination options for activity list endpoints.
type ListActivityOptions struct {
	Page    int `json:"page,omitempty"`
	PerPage int `json:"per_page,omitempty"`
}
