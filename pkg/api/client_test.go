package api

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("gitee.com", "test-token")
	if client == nil {
		t.Error("NewClient returned nil")
	}
}
