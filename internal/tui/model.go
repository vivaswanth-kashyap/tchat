package tui

import (
	"log"
	"sync"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/vivaswanth-kashyap/tchat/internal/app"
	"github.com/vivaswanth-kashyap/tchat/internal/models"
)

type ChatModel struct {
	appClient *app.Client
	logger    *log.Logger
	config    *app.Config

	// UI Components\
	viewport      viewport.Model
	messageInput  textinput.Model
	statusSpinner spinner.Model
	channelList   list.Model
	userList      list.Model

	// Application State
	messages       []models.Message
	currentChannel models.Channel
	currentUser    models.User
	dmRecipient    models.User

	// Connection and Loading Status
	connected      bool
	connecting     bool
	loadingHistory bool
	errorMessage   string

	// Sync for go routines
	mu sync.Mutex

	//view management
	mode  AppMode // Current mode of the app (e.g., ChatMode, ConfigMode)
	focus FocusableComponent

	//Dimensions for layout
	width, height int
}

type AppMode int

const (
	ChatMode AppMode = iota
)

// UI element currently has focus
type FocusableComponent int

const (
	InputFocus FocusableComponent = iota
)

// NewChatModel initializes a new ChatModel with default states and components.

// func NewChatModel(appClient *app.Client, appConfig *app.Config) ChatModel {
// 	ti := textinput.New()
// 	ti.Placeholder = "Type a message or /command..."
// 	ti.Focus()         // Start with input field focused
// 	ti.CharLimit = 200 // Max characters for a message
// 	ti.Width = 80      // Default width, will be updated by Viewport

// 	vp := viewport.New(0, 0) // Width/height will be set in Update
// 	vp.YOffset = 0           // Start at top of history
// 	// Optional: vp.KeyMap.PgUp / vp.KeyMap.PgDown for navigation

// 	s := spinner.New()
// 	s.Spinner = spinner.Dot // Choose a spinner style
// 	s.Style = styles.Theme.SpinnerStyle

// 	// Initialize with placeholder data
// 	initialMessages := []models.Message{
// 		{ID: "sys-welcome", Content: "Welcome to Tchat! Connecting...", Sender: "System", Timestamp: time.Now()},
// 	}

// 	return ChatModel{
// 		appClient: appClient,
// 		logger:    log.Default(), // Consider using charm.sh/log for better logging
// 		config:    appConfig,

// 		viewport:      vp,
// 		messageInput:  ti,
// 		statusSpinner: s,
// 		// channelList:    list.New(...), // Initialize if you're using it
// 		// userList:       list.New(...),

// 		messages:       initialMessages,
// 		connected:      false,
// 		connecting:     true, // Assume connecting on startup
// 		loadingHistory: false,
// 		errorMessage:   "",

// 		mode:  ChatMode,
// 		focus: InputFocus,
// 	}
// }

// Ensure you define custom message types in internal/tui/messages.go
// Example:
// type MsgReceived models.Message
// type ConnectionStatusMsg bool
// type ErrorMsg string
