package response

// CheckRun represents a Gitee check run
type CheckRun struct {
	ID          int64   `json:"id"`
	HeadSHA     string  `json:"head_sha"`
	URL         string  `json:"url"`
	HTMLURL     string  `json:"html_url"`
	DetailsURL  string  `json:"details_url"`
	Status      string  `json:"status"`
	Conclusion  string  `json:"conclusion"`
	StartedAt   string  `json:"started_at"`
	CompletedAt string  `json:"completed_at"`
	Output      *Output `json:"output"`
	Name        string  `json:"name"`
}

// Output represents the check run output
type Output struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Text    string `json:"text"`
}

// CheckAnnotation represents a check run annotation
type CheckAnnotation struct {
	Path            string `json:"path"`
	StartLine       int    `json:"start_line"`
	EndLine         int    `json:"end_line"`
	StartColumn     int    `json:"start_column,omitempty"`
	EndColumn       int    `json:"end_column,omitempty"`
	AnnotationLevel string `json:"annotation_level"`
	Title           string `json:"title,omitempty"`
	Message         string `json:"message"`
	RawDetails      string `json:"raw_details,omitempty"`
	BlobHref        string `json:"blob_href"`
}

// CreateCheckRunOptions is the request body for creating a check run
type CreateCheckRunOptions struct {
	Name          string `json:"name"`
	HeadSHA       string `json:"head_sha"`
	DetailsURL    string `json:"details_url,omitempty"`
	Status        string `json:"status,omitempty"`
	OutputTitle   string `json:"output[title]"`
	OutputSummary string `json:"output[summary]"`
	OutputText    string `json:"output[text,omitempty"`
}

// UpdateCheckRunOptions is the request body for updating a check run
type UpdateCheckRunOptions struct {
	DetailsURL    string `json:"details_url,omitempty"`
	Status        string `json:"status,omitempty"`
	StartedAt     string `json:"started_at,omitempty"`
	Conclusion    string `json:"conclusion,omitempty"`
	CompletedAt   string `json:"completed_at,omitempty"`
	OutputTitle   string `json:"output[title]"`
	OutputSummary string `json:"output[summary]"`
	OutputText    string `json:"output[text,omitempty"`
}

// ListCheckRunsOptions contains optional parameters for ListCommitCheckRuns
type ListCheckRunsOptions struct {
	Page    int
	PerPage int
}
