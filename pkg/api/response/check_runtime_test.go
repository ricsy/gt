package response

import (
	"encoding/json"
	"testing"
)

func TestCheckRunListUnmarshalMatchesLiveResponseShape(t *testing.T) {
	payload := []byte(`{
		"total_count": 1,
		"check_runs": [
			{
				"id": 26691196,
				"head_sha": "6ad580ae6d4405fb129f0ce2fac180c1c1e318e3",
				"url": "https://gitee.com/api/v5/repos/yimeng555/example/check-runs/26691196",
				"html_url": "https://gitee.com/yimeng555/example/run?id=26691196",
				"details_url": null,
				"status": "queued",
				"conclusion": null,
				"started_at": "2026-06-11T20:57:56+08:00",
				"completed_at": null,
				"output": {
					"title": null,
					"summary": null,
					"text": null,
					"annotations_count": 0,
					"annotations_url": "https://gitee.com/api/v5/repos/yimeng555/example/check-runs/26691196/annotations"
				},
				"name": "smoke-check"
			}
		]
	}`)

	var result CheckRunList
	if err := json.Unmarshal(payload, &result); err != nil {
		t.Fatalf("json.Unmarshal(check run list) returned error: %v", err)
	}

	if len(result.CheckRuns) != 1 {
		t.Fatalf("len(result.CheckRuns) = %d, want 1", len(result.CheckRuns))
	}
}
