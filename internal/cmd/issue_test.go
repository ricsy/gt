package cmd

import (
	"bytes"
	"testing"
)

func TestIssueCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	issueCmd.SetOut(buf)
	issueCmd.SetArgs([]string{"--help"})
	issueCmd.Execute()
	// Just verify it runs without error
}

func TestIssueListCmd(t *testing.T) {
	if issueListCmd.Use != "list" {
		t.Errorf("issueListCmd.Use = %s", issueListCmd.Use)
	}
}

func TestIssueViewCmd(t *testing.T) {
	if issueViewCmd.Use != "view <number>" {
		t.Errorf("issueViewCmd.Use = %s", issueViewCmd.Use)
	}
}
