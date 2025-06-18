package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tchat",
	Short: "messaging application within your terminal",
	Long:  "TUI and cli based messaging and file sharing application within the terminal",
}

func init() {
	rootCmd.AddCommand(sendCmd)
	rootCmd.AddCommand(readCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
