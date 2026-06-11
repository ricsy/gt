package util

import "testing"

func TestBuildQueryEncodesValues(t *testing.T) {
	got := BuildQuery(
		"q", "hello world",
		"labels", "bug&urgent",
		"empty", "",
	)

	want := "labels=bug%26urgent&q=hello+world"
	if got != want {
		t.Fatalf("BuildQuery() = %s, want %s", got, want)
	}
}
