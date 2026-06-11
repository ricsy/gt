package response

import "encoding/json"

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

// Languages represents programming language bytes.
// swagger 定义为 object，但 live API 在无语言数据时会返回 {"languages":[]}.
type Languages struct {
	Languages map[string]int `json:"languages"`
}

// UnmarshalJSON 兼容空语言统计时的数组返回。
func (l *Languages) UnmarshalJSON(data []byte) error {
	type languagesObject struct {
		Languages map[string]int `json:"languages"`
	}
	var objectValue languagesObject
	if err := json.Unmarshal(data, &objectValue); err == nil && objectValue.Languages != nil {
		l.Languages = objectValue.Languages
		return nil
	}

	type languagesArray struct {
		Languages []any `json:"languages"`
	}
	var arrayValue languagesArray
	if err := json.Unmarshal(data, &arrayValue); err != nil {
		return err
	}

	l.Languages = map[string]int{}
	return nil
}
