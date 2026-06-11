package api

import "testing"

func TestBuildSearchQueryAddsQuestionMark(t *testing.T) {
	got := buildSearchQuery("q", "gt", "", "", 1, 1)

	if got != "?page=1&per_page=1&q=gt" {
		t.Fatalf("buildSearchQuery() = %q, want %q", got, "?page=1&per_page=1&q=gt")
	}
}
