package cmd

import "testing"

func TestActivityCommand(t *testing.T) {
	if activityCmd.Use != "activity" {
		t.Errorf("expected use 'activity', got %s", activityCmd.Use)
	}
	if len(activityCmd.Commands()) != 3 {
		t.Errorf("expected 3 subcommands, got %d", len(activityCmd.Commands()))
	}
}

func TestActivityWatchCommand(t *testing.T) {
	if activityWatchCmd.Use != "watch" {
		t.Errorf("expected use 'watch', got %s", activityWatchCmd.Use)
	}
	if len(activityWatchCmd.Commands()) != 3 {
		t.Errorf("expected 3 subcommands, got %d", len(activityWatchCmd.Commands()))
	}
}

func TestActivityCommandFlags(t *testing.T) {
	if activityEventsCmd.Flags().Lookup("repo") == nil {
		t.Error("expected events --repo flag")
	}
	if activityEventsCmd.Flags().Lookup("org") == nil {
		t.Error("expected events --org flag")
	}
	if activityEventsCmd.Flags().Lookup("public") == nil {
		t.Error("expected events --public flag")
	}
	if activityEventsCmd.Flags().Lookup("received") == nil {
		t.Error("expected events --received flag")
	}
	if activityWatchRepoCmd.Flags().Lookup("repo") == nil {
		t.Error("expected watch repo --repo flag")
	}
	if activitySubscribersCmd.Flags().Lookup("repo") == nil {
		t.Error("expected subscribers --repo flag")
	}
}
