// gt is a CLI tool for Gitee
// API Documentation: https://help.gitee.com/openapi/v5
package main

import (
	"os"

	"github.com/ricsy/gt/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
