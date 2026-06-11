package cmd

import "testing"

func TestNotificationCommandFlags(t *testing.T) {
	if notificationListCmd.Flags().Lookup("ids") == nil {
		t.Error("expected repo notification list --ids flag")
	}
	if notificationListCmd.Flags().Lookup("page") == nil {
		t.Error("expected repo notification list --page flag")
	}
	if notificationListCmd.Flags().Lookup("per-page") == nil {
		t.Error("expected repo notification list --per-page flag")
	}
	if threadListCmd.Flags().Lookup("ids") == nil {
		t.Error("expected thread list --ids flag")
	}
	if threadListCmd.Flags().Lookup("page") == nil {
		t.Error("expected thread list --page flag")
	}
	if messageListCmd.Flags().Lookup("ids") == nil {
		t.Error("expected message list --ids flag")
	}
	if messageListCmd.Flags().Lookup("per-page") == nil {
		t.Error("expected message list --per-page flag")
	}
	if messageMarkAllReadCmd.Flags().Lookup("ids") == nil {
		t.Error("expected message mark-all-read --ids flag")
	}
}
