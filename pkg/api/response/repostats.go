package response

// Contributor represents a repository contributor
type Contributor struct {
	Email         string `json:"email"`
	Name          string `json:"name"`
	Contributions int    `json:"contributions"`
}

// GetTrafficDataOptions contains optional parameters for GetTrafficData
type GetTrafficDataOptions struct {
	StartDay string `json:"start_day,omitempty"`
	EndDay   string `json:"end_day,omitempty"`
}

// TrafficDataDesc represents daily traffic data
type TrafficDataDesc struct {
	Bucket      int `json:"bucket"`
	IP          int `json:"ip"`
	Pull        int `json:"pull"`
	Push        int `json:"push"`
	DownloadZip int `json:"download_zip"`
}

// TrafficDataSummary represents traffic data summary
type TrafficDataSummary struct {
	IP          int `json:"ip"`
	Pull        int `json:"pull"`
	Push        int `json:"push"`
	DownloadZip int `json:"download_zip"`
}

// TrafficData represents repository traffic data
type TrafficData struct {
	Counts  []TrafficDataDesc  `json:"counts"`
	Summary TrafficDataSummary `json:"summary"`
}

// Languages represents programming language bytes
type Languages struct {
	Languages map[string]int `json:"languages"`
}
