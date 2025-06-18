package cli

import (
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

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "tchat read @Username",
	Long:  "tchat read @Username outputs to stdout can be used to pipe",
	Run:   readMessage,
}

func readMessage(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		_ = cmd.Usage()
		os.Exit(1)
	}

	recipient := strings.TrimPrefix(args[0], "@")
	if recipient == "" {
		fmt.Fprintf(os.Stderr, "Invalid recipient")
		os.Exit(1)
	}

	if err := readHttpMessage(recipient); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read message: %v\n", err)
	}
}

func readHttpMessage(recipient string) error {
	var currentUser models.User
	if err := db.DB.First(&currentUser).Error; err != nil {
		return fmt.Errorf("No authenticated User found. Please Login first")
	}

	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		serverUrl = "http://localhost:8080"
	}

	url := fmt.Sprintf("%s/messages/last?sender_id=%s&receiver_username=%s", serverUrl, currentUser.ServerID, recipient)

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("No messages found between you and %s\n", recipient)
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("server error (%d):%s", resp.StatusCode, body)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if message, ok := response["message"].(map[string]interface{}); ok {
		body := message["body"].(string)
		fmt.Printf(body)
	} else {
		fmt.Printf("No message data found\n")
	}

	return nil
}
