package response

// Email represents a Gitee email address
type Email struct {
	Email    string `json:"email"`
	State    string `json:"state"`
	Verified bool   `json:"verified"`
}
