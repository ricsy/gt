package cmd

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/ricsy/gt/pkg/config"
	"github.com/spf13/cobra"
)

var apiHTTPClient = &http.Client{Timeout: 30 * time.Second}

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

			client, err := getClient()
			if err != nil {
				return err
			}
			token := client.Token()

			url := config.ApiUrl(config.DefaultHost) + path

			var bodyReader io.Reader
			if bodyFlag != "" {
				bodyReader = bytes.NewBufferString(bodyFlag)
			}

			req, err := http.NewRequest(method, url, bodyReader)
			if err != nil {
				return err
			}

			req.Header.Set("Authorization", "token "+token)
			req.Header.Set("Content-Type", "application/json")

			resp, err := apiHTTPClient.Do(req)
			if err != nil {
				return err
			}
			defer func() { _ = resp.Body.Close() }()

			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			if rawFlag {
				cmd.Println(string(respBody))
			} else {
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, respBody, "", "  "); err != nil {
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
