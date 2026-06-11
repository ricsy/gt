package api

import "testing"

func TestIssueLabelRemoveEndpointBuildsNamePath(t *testing.T) {
	got := IssueLabels.Remove.Build("owner", "repo", "ISSUE-1", "bug")
	want := "/repos/owner/repo/issues/ISSUE-1/labels/bug"
	if got != want {
		t.Fatalf("IssueLabels.Remove.Build() = %q, want %q", got, want)
	}
}
