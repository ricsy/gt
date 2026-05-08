package cmd

import "testing"

func TestGitDataCommand(t *testing.T) {
	if gitDataCmd.Use != "gitdata" {
		t.Errorf("expected use 'gitdata', got %s", gitDataCmd.Use)
	}
	if len(gitDataCmd.Commands()) != 3 {
		t.Errorf("expected 3 subcommands, got %d", len(gitDataCmd.Commands()))
	}
}

func TestGitDataCommandFlags(t *testing.T) {
	if gitDataBlobCmd.Flags().Lookup("repo") == nil {
		t.Error("expected blob --repo flag")
	}
	if gitDataTreeCmd.Flags().Lookup("repo") == nil {
		t.Error("expected tree --repo flag")
	}
	if gitDataTreeCmd.Flags().Lookup("recursive") == nil {
		t.Error("expected tree --recursive flag")
	}
	if gitDataMetricsCmd.Flags().Lookup("repo") == nil {
		t.Error("expected metrics --repo flag")
	}
}
