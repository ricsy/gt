package response

// ListMilestonesOptions contains optional parameters for ListMilestones
type ListMilestonesOptions struct {
	State     string // open, closed, all
	Sort      string // due_on
	Direction string // asc, desc
	Page      int
	PerPage   int
}

// CreateMilestoneOptions is the request body for creating a milestone
type CreateMilestoneOptions struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
	DueOn       string `json:"due_on"`
}

// UpdateMilestoneOptions is the request body for updating a milestone
type UpdateMilestoneOptions struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	State       string `json:"state,omitempty"`
	DueOn       string `json:"due_on,omitempty"`
}
