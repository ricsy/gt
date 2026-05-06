package response

// Label represents a Gitee label
type Label struct {
	ID           int64  `json:"id"`
	Color        string `json:"color"`
	Name         string `json:"name"`
	RepositoryID int64  `json:"repository_id"`
	URL          string `json:"url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// ProjectLabel represents a Gitee project label
type ProjectLabel struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Ident string `json:"ident"`
}

// CreateLabelOptions is the request body for creating a label
type CreateLabelOptions struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// UpdateLabelOptions is the request body for updating a label
type UpdateLabelOptions struct {
	Name  string `json:"name,omitempty"`
	Color string `json:"color,omitempty"`
}
