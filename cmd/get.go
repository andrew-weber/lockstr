package cmd

import (
	"fmt"
	"strings"

	entries "github.com/andrew-weber/lockstr/lib"
	"github.com/atotto/clipboard"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip04"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a password",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("No key provided. Try 'lockstr get <key>'")
			return
		}

		entry := entries.GetEntry(strings.TrimSpace(string(args[0])))

		if entry == nil {
			fmt.Println("No entry found")
			return
		}

		sk := viper.GetString("KEY")
		pub, _ := nostr.GetPublicKey(sk)

		shared, _ := nip04.ComputeSharedSecret(pub, sk)
		pwd, _ := nip04.Decrypt(entry.Password, shared)

		clipboard.WriteAll(pwd)
		fmt.Println("Password copied to clipboard")
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
