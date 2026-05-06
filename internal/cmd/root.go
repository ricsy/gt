package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.1.0-alpha"

var rootCmd = &cobra.Command{
	Use:     "gt",
	Short:   "gt is a CLI for Gitee",
	Long:    `Gitee's command line tool.`,
	Version: version,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gt version", version)
		},
	})
}
