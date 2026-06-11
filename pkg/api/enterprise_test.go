package api

import "testing"

func TestBuildEnterprisePullRequestsQueryIncludesSupportedFilters(t *testing.T) {
	got := buildEnterprisePullRequestsQuery(ListEnterprisePullRequestsOptions{
		IssueNumber:     "ISSUE-1",
		Repo:            "demo",
		ProgramID:       12,
		State:           "open",
		Head:            "user:feature",
		Base:            "main",
		Sort:            "updated",
		Since:           "2026-06-11T00:00:00Z",
		Direction:       "asc",
		MilestoneNumber: 34,
		Labels:          "bug,help wanted",
		Page:            2,
		PerPage:         50,
	})

	want := "base=main&direction=asc&head=user%3Afeature&issue_number=ISSUE-1&labels=bug%2Chelp+wanted&milestone_number=34&page=2&per_page=50&program_id=12&repo=demo&since=2026-06-11T00%3A00%3A00Z&sort=updated&state=open"
	if got != want {
		t.Fatalf("buildEnterprisePullRequestsQuery() = %q, want %q", got, want)
	}
}
