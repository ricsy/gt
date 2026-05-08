package cmd

import (
	"bytes"
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func TestVersionCommand(t *testing.T) {
	buf := new(bytes.Buffer)

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("gt version", version)
		},
	}
	cmd.SetOut(buf)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}

	expected := "gt version " + version + "\n"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestRootCommandHelp(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := Execute()
	if err != nil {
		t.Errorf("Execute() error = %v", err)
	}
}

func TestRootCommandHasTimeoutFlag(t *testing.T) {
	flag := rootCmd.PersistentFlags().Lookup("timeout")
	if flag == nil {
		t.Fatal("expected root command to define a timeout flag")
	}
}

func TestCommandHTTPClientUsesRequestTimeout(t *testing.T) {
	originalTimeout := requestTimeout
	t.Cleanup(func() {
		requestTimeout = originalTimeout
	})

	requestTimeout = 7 * time.Second

	client := newCommandHTTPClient()
	if client.Timeout != requestTimeout {
		t.Fatalf("expected timeout %s, got %s", requestTimeout, client.Timeout)
	}
}
