package response

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

type Output struct {
	Title   string `json:"title"`
	Summary string `json:"summary"`
	Text    string `json:"text"`
}

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

type CreateCheckRunOptions struct {
	PullRequestID                int      `json:"pull_request_id,omitempty"`
	DetailsURL                   string   `json:"details_url,omitempty"`
	Status                       string   `json:"status,omitempty"`
	StartedAt                    string   `json:"started_at,omitempty"`
	Conclusion                   string   `json:"conclusion,omitempty"`
	CompletedAt                  string   `json:"completed_at,omitempty"`
	OutputTitle                  string   `json:"output[title]"`
	OutputSummary                string   `json:"output[summary]"`
	OutputText                   string   `json:"output[text],omitempty"`
	OutputAnnotationsPath        []string `json:"output[annotations][path]"`
	OutputAnnotationsStartLine   []int    `json:"output[annotations][start_line]"`
	OutputAnnotationsEndLine     []int    `json:"output[annotations][end_line]"`
	OutputAnnotationsStartColumn []int    `json:"output[annotations][start_column],omitempty"`
	OutputAnnotationsEndColumn   []int    `json:"output[annotations][end_column],omitempty"`
	OutputAnnotationsLevel       []string `json:"output[annotations][annotation_level]"`
	OutputAnnotationsMessage     []string `json:"output[annotations][message]"`
	OutputAnnotationsTitle       []string `json:"output[annotations][title],omitempty"`
	OutputAnnotationsRawDetails  []string `json:"output[annotations][raw_details],omitempty"`
	OutputImagesAlt              []string `json:"output[images][alt]"`
	OutputImagesURL              []string `json:"output[images][image_url]"`
	OutputImagesCaption          []string `json:"output[images][caption],omitempty"`
	ActionsLabel                 []string `json:"actions[label]"`
	ActionsDescription           []string `json:"actions[description]"`
	ActionsIdentifier            []string `json:"actions[identifier]"`
	Name                         string   `json:"name"`
	HeadSHA                      string   `json:"head_sha"`
}

type UpdateCheckRunOptions struct {
	PullRequestID                int      `json:"pull_request_id,omitempty"`
	DetailsURL                   string   `json:"details_url,omitempty"`
	Status                       string   `json:"status,omitempty"`
	StartedAt                    string   `json:"started_at,omitempty"`
	Conclusion                   string   `json:"conclusion,omitempty"`
	CompletedAt                  string   `json:"completed_at,omitempty"`
	OutputTitle                  string   `json:"output[title]"`
	OutputSummary                string   `json:"output[summary]"`
	OutputText                   string   `json:"output[text],omitempty"`
	OutputAnnotationsPath        []string `json:"output[annotations][path]"`
	OutputAnnotationsStartLine   []int    `json:"output[annotations][start_line]"`
	OutputAnnotationsEndLine     []int    `json:"output[annotations][end_line]"`
	OutputAnnotationsStartColumn []int    `json:"output[annotations][start_column],omitempty"`
	OutputAnnotationsEndColumn   []int    `json:"output[annotations][end_column],omitempty"`
	OutputAnnotationsLevel       []string `json:"output[annotations][annotation_level]"`
	OutputAnnotationsMessage     []string `json:"output[annotations][message]"`
	OutputAnnotationsTitle       []string `json:"output[annotations][title],omitempty"`
	OutputAnnotationsRawDetails  []string `json:"output[annotations][raw_details],omitempty"`
	OutputImagesAlt              []string `json:"output[images][alt]"`
	OutputImagesURL              []string `json:"output[images][image_url]"`
	OutputImagesCaption          []string `json:"output[images][caption],omitempty"`
	ActionsLabel                 []string `json:"actions[label]"`
	ActionsDescription           []string `json:"actions[description]"`
	ActionsIdentifier            []string `json:"actions[identifier]"`
	Name                         string   `json:"name,omitempty"`
}

type ListCheckRunsOptions struct {
	Page    int
	PerPage int
}
