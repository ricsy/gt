package cmd

import "testing"

func TestParseGistFiles(t *testing.T) {
	files, err := parseGistFiles(`{"smoke.txt":{"content":"hello"}}`)
	if err != nil {
		t.Fatalf("parseGistFiles() returned error: %v", err)
	}

	if files["smoke.txt"]["content"] != "hello" {
		t.Fatalf("files[smoke.txt][content] = %q, want %q", files["smoke.txt"]["content"], "hello")
	}
}

func TestParseGistFilesRejectsEmptyObject(t *testing.T) {
	_, err := parseGistFiles(`{}`)
	if err == nil {
		t.Fatal("parseGistFiles() error = nil, want non-nil")
	}
}
