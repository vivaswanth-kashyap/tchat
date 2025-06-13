package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var sendCmd = &cobra.Command{
	Use:   "send",
	Short: "tchat send @Username 'message'",
	Long:  "tchat send @Username 'message' or pipe the output of other unix commands",
	Run:   sendMessage,
}

func sendMessage(cmd *cobra.Command, args []string) {
	if len(args) < 1 || len(args) > 2 {
		_ = cmd.Usage()
		os.Exit(1)
	}

	recipient := strings.TrimPrefix(args[0], "@")
	if recipient == "" {
		fmt.Fprintln(os.Stderr, "invalid recipient")
		os.Exit(1)
	}

	var message string
	if len(args) == 2 {
		message = args[1]
	} else {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to read stdin: %v\n", err)
			os.Exit(1)
		}
		message = strings.TrimSpace(string(data))
	}

	if message == "" {
		fmt.Fprintln(os.Stderr, "message cannot be empty")
		os.Exit(1)
	}

	fmt.Printf("Sending message %q: %q\n", recipient, message)
}
