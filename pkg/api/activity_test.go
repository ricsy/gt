package api

import "testing"

func TestBuildActivityQuery(t *testing.T) {
	got := buildActivityQuery(ListActivityOptions{
		Page:    1,
		PerPage: 20,
	})

	want := "?page=1&per_page=20"
	if got != want {
		t.Fatalf("buildActivityQuery() = %q, want %q", got, want)
	}
}
