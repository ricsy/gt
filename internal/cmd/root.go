package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/ricsy/gt/pkg/api"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

var version = "0.1.0-beta"
var requestTimeout = api.DefaultTimeout
var commandHost string
var commandEnvFile string

var rootCmd = &cobra.Command{
	Use:          "gt",
	Short:        "gt is a CLI for Gitee",
	Long:         `Gitee's command line tool.`,
	Version:      version,
	SilenceUsage: true,
}

func Execute() error {
	return rootCmd.Execute()
}

func newCommandAPIClient(host, token string) *api.Client {
	return api.NewClientWithTimeout(host, token, requestTimeout)
}

func newCommandHTTPClient() *http.Client {
	timeout := requestTimeout
	if timeout <= 0 {
		timeout = api.DefaultTimeout
	}
	return &http.Client{Timeout: timeout}
}

func resolveCommandEnvFile() string {
	if commandEnvFile != "" {
		return commandEnvFile
	}
	return strings.TrimSpace(os.Getenv("GT_ENV_FILE"))
}

func loadCommandEnvFile() error {
	return config.LoadEnvFile(resolveCommandEnvFile())
}

func init() {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		return loadCommandEnvFile()
	}
	rootCmd.PersistentFlags().DurationVar(&requestTimeout, "timeout", api.DefaultTimeout, "HTTP request timeout")
	rootCmd.PersistentFlags().StringVar(&commandHost, "host", "", "Git host (defaults to configured host or gitee.com)")
	rootCmd.PersistentFlags().StringVar(&commandEnvFile, "env-file", "", "Environment file to load before executing the command")

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gt version", version)
		},
	})
}
