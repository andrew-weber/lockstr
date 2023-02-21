package cmd

import (
	"fmt"

	entries "github.com/andrew-weber/lockstr/lib"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List password keys",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		entries := entries.ListEntries()

		for _, entry := range entries {
			fmt.Println(entry.Key)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
