package api

import "testing"

func TestBuildMilestonesQueryAddsQuestionMark(t *testing.T) {
	got := buildMilestonesQuery(ListMilestonesOptions{
		State:     "open",
		Sort:      "due_on",
		Direction: "asc",
	})

	if got != "?direction=asc&sort=due_on&state=open" {
		t.Fatalf("buildMilestonesQuery() = %q, want %q", got, "?direction=asc&sort=due_on&state=open")
	}
}
