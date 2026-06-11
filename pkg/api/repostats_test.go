package api

import (
	"encoding/json"
	"testing"

	"github.com/ricsy/gt/pkg/api/response"
)

func TestBuildListBranchesQueryAddsQuestionMark(t *testing.T) {
	got := buildListBranchesQuery(ListBranchesOptions{
		Sort:      "updated",
		Direction: "desc",
		Page:      1,
		PerPage:   20,
	})

	want := "?direction=desc&page=1&per_page=20&sort=updated"
	if got != want {
		t.Fatalf("buildListBranchesQuery() = %q, want %q", got, want)
	}
}

func TestBuildListForksQueryAddsQuestionMark(t *testing.T) {
	got := buildListForksQuery(ListForksOptions{
		Sort:    "newest",
		Page:    2,
		PerPage: 10,
	})

	want := "?page=2&per_page=10&sort=newest"
	if got != want {
		t.Fatalf("buildListForksQuery() = %q, want %q", got, want)
	}
}

func TestLanguagesUnmarshalAllowsEmptyArray(t *testing.T) {
	payload := []byte(`{"languages":[]}`)

	var languages response.Languages
	if err := json.Unmarshal(payload, &languages); err != nil {
		t.Fatalf("json.Unmarshal(languages) returned error: %v", err)
	}

	if languages.Languages == nil {
		t.Fatal("languages.Languages is nil")
	}
	if len(languages.Languages) != 0 {
		t.Fatalf("len(languages.Languages) = %d, want 0", len(languages.Languages))
	}
}

func TestTrafficDataUnmarshalAllowsStringBucket(t *testing.T) {
	payload := []byte(`{"counts":[{"bucket":"1781136000","ip":0,"pull":1,"push":1,"download_zip":0}],"summary":{"ip":0,"pull":1,"push":1,"download_zip":0}}`)

	var traffic response.TrafficData
	if err := json.Unmarshal(payload, &traffic); err != nil {
		t.Fatalf("json.Unmarshal(traffic) returned error: %v", err)
	}

	if len(traffic.Counts) != 1 {
		t.Fatalf("len(traffic.Counts) = %d, want 1", len(traffic.Counts))
	}
	if traffic.Counts[0].Bucket != "1781136000" {
		t.Fatalf("traffic.Counts[0].Bucket = %q, want %q", traffic.Counts[0].Bucket, "1781136000")
	}
}
