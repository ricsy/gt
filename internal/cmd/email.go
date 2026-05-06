package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var emailCmd = &cobra.Command{
	Use:   "email",
	Short: "Manage emails",
	Long:  `Commands for managing Gitee emails`,
}

var emailListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all emails",
	RunE:  emailList,
}

func init() {
	emailCmd.AddCommand(emailListCmd)
	rootCmd.AddCommand(emailCmd)
}

func emailList(cmd *cobra.Command, args []string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	emails, err := client.ListEmails()
	if err != nil {
		return fmt.Errorf("failed to list emails: %w", err)
	}

	if len(emails) == 0 {
		fmt.Println("No emails found")
		return nil
	}

	for _, e := range emails {
		verified := ""
		if e.Verified {
			verified = " [verified]"
		}
		state := e.State
		if state == "" {
			state = "unknown"
		}
		fmt.Printf("%s (%s)%s\n", e.Email, state, verified)
	}
	return nil
}
