package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/ricsy/gt/pkg/auth"
	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

func newApiCmd() *cobra.Command {
	var bodyFlag string
	var rawFlag bool

	cmd := &cobra.Command{
		Use:   "api <method> <path>",
		Short: "Make a raw API request",
		Long: `Make a raw request to the Gitee API.

Examples:
  gt api GET /user
  gt api POST /repos/owner/repo/issues --body '{"title":"New Issue"}'
  gt api GET /user --raw
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			method := strings.ToUpper(args[0])
			path := args[1]

			// Get auth token
			token, err := auth.GetToken(config.DefaultHost)
			if err != nil {
				return fmt.Errorf("authentication required: %w", err)
			}

			// Build URL
			url := fmt.Sprintf("https://gitee.com/api/v5%s", path)

			// Prepare body
			var bodyReader io.Reader
			if bodyFlag != "" {
				bodyReader = bytes.NewBufferString(bodyFlag)
			}

			// Create request
			req, err := http.NewRequest(method, url, bodyReader)
			if err != nil {
				return err
			}

			req.Header.Set("Authorization", "token "+token)
			req.Header.Set("Content-Type", "application/json")

			// Execute request
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			defer func() { _ = resp.Body.Close() }()

			// Read response
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if rawFlag {
				// Output raw response
				cmd.Println(string(respBody))
			} else {
				// Pretty print JSON
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, respBody, "", "  "); err != nil {
					// If not valid JSON, just output raw
					cmd.Println(string(respBody))
				} else {
					cmd.Println(prettyJSON.String())
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&bodyFlag, "body", "", "Request body (JSON)")
	cmd.Flags().BoolVar(&rawFlag, "raw", false, "Output raw response without formatting")

	return cmd
}

func init() {
	rootCmd.AddCommand(newApiCmd())
}
