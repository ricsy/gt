package response

import (
	"encoding/json"
	"testing"
)

func TestIssueUnmarshalMatchesLiveResponseShape(t *testing.T) {
	payload := []byte(`{
		"id": 1,
		"number": "IJTYTH",
		"state": "open",
		"title": "smoke issue",
		"scheduled_time": 0
	}`)

	var issue Issue
	if err := json.Unmarshal(payload, &issue); err != nil {
		t.Fatalf("json.Unmarshal(issue) returned error: %v", err)
	}

	if issue.Number != "IJTYTH" {
		t.Fatalf("issue.Number = %q, want %q", issue.Number, "IJTYTH")
	}
}

func TestIssueCommentUnmarshalMatchesLiveResponseShape(t *testing.T) {
	payload := []byte(`{
		"id": 1,
		"body": "comment",
		"source": {},
		"target": {}
	}`)

	var comment IssueComment
	if err := json.Unmarshal(payload, &comment); err != nil {
		t.Fatalf("json.Unmarshal(comment) returned error: %v", err)
	}
}
