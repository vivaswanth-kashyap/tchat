package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vivaswanth-kashyap/tchat/internal/db"
	"github.com/vivaswanth-kashyap/tchat/internal/models"
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

	if err := sendHttpMessage(recipient, message); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to send message: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Message sent to %s ✓\n", recipient)
}

func sendHttpMessage(recipient, message string) error {
	var currentUser models.User
	if err := db.DB.First(&currentUser).Error; err != nil {
		return fmt.Errorf("no authenticated user found. Please login first")
	}

	payload := map[string]string{
		"sender_id":         currentUser.ServerID,
		"receiver_username": recipient,
		"body":              message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	serverURL := os.Getenv("SERVER_URL")
	if serverURL == "" {
		serverURL = "http://localhost:8080"
	}

	resp, err := http.Post(serverURL+"/messages",
		"application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d): %s", resp.StatusCode, body)
	}

	return nil
}
