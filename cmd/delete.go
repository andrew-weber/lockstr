package cmd

import (
	"strings"

	entries "github.com/andrew-weber/lockstr/lib"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Remove a password",
	Long:  `Remove a password from the password manager`,
	Run: func(cmd *cobra.Command, args []string) {
		entries.DeleteEntry(strings.TrimSpace(string(args[0])))
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
