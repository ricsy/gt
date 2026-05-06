package response

// UserBasic represents basic user info
type UserBasic struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}

// User is an alias for UserBasic
type User = UserBasic
