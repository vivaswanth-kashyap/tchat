// cli/chat.go
package cli

import (
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Launches the interactive chat interface (TUI)",
	Long:  `The 'chat' command starts the full-screen terminal user interface for messaging.`,
	Run:   launchTui,
}

func launchTui(cmd *cobra.Command, args []string) {

}
