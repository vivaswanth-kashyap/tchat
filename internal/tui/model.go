package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivaswanth-kashyap/tchat/internal/models"
)

type ViewType int
type InputType int

const (
	LoginView ViewType = iota
	ChatView
	ChannelListView
	SettingsView
)

const (
	NoInput InputType = iota
	UsernameInput
	PasswordInput
	MessageInput
	ChannelSearchInput
)

type Model struct {
	// Authentication state
	isAuthenticated bool
	authToken       string
	currentUser     models.User
	// UI state
	currentView  ViewType
	focusedInput InputType
	// Messaging state
	messages       []models.Message
	currentChannel string
	channels       []models.Channel
	// Input handling
	usernameInput textinput.Model
	passwordInput textinput.Model
	messageInput  textinput.Model
	// UI components
	viewport viewport.Model
	list     list.Model
	// Application state
	loading      bool
	error        error
	windowWidth  int
	windowHeight int
}

func initialModel() Model {
	// Username input
	usernameInput := textinput.New()
	usernameInput.Placeholder = "Username"
	usernameInput.Focus()
	usernameInput.CharLimit = 50
	usernameInput.Width = 30

	// Password input
	passwordInput := textinput.New()
	passwordInput.Placeholder = "Password"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.CharLimit = 50
	passwordInput.Width = 30

	// Message input
	messageInput := textinput.New()
	messageInput.Placeholder = "Type a message..."
	messageInput.CharLimit = 500
	messageInput.Width = 50

	// Viewport and list
	vp := viewport.New(80, 20)
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Channels"

	return Model{
		isAuthenticated: false,
		authToken:       "",
		currentUser:     models.User{},
		currentView:     LoginView,
		focusedInput:    UsernameInput,
		messages:        []models.Message{},
		currentChannel:  "",
		channels:        []models.Channel{},
		usernameInput:   usernameInput,
		passwordInput:   passwordInput,
		messageInput:    messageInput,
		viewport:        vp,
		list:            l,
		loading:         false,
		error:           nil,
		windowWidth:     80,
		windowHeight:    24,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.currentView == LoginView {
				if m.focusedInput == UsernameInput {
					m.focusedInput = PasswordInput
					m.usernameInput.Blur()
					m.passwordInput.Focus()
				} else {
					m.focusedInput = UsernameInput
					m.passwordInput.Blur()
					m.usernameInput.Focus()
				}
			}
		case "enter":
			if m.currentView == LoginView {
				return m, m.handleLogin()
			} else if m.currentView == ChatView && m.focusedInput == MessageInput {
				return m, m.handleSendMessage()
			}
		}

		// Update the focused input
		var cmd tea.Cmd
		switch m.focusedInput {
		case UsernameInput:
			m.usernameInput, cmd = m.usernameInput.Update(msg)
		case PasswordInput:
			m.passwordInput, cmd = m.passwordInput.Update(msg)
		case MessageInput:
			m.messageInput, cmd = m.messageInput.Update(msg)
		}
		return m, cmd

	case tea.WindowSizeMsg:
		m.windowWidth = msg.Width
		m.windowHeight = msg.Height

		// Resize viewport
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4

		// Resize list
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case LoginSuccessMsg:
		m.isAuthenticated = true
		m.authToken = msg.Token
		m.currentUser = msg.User
		m.currentView = ChatView
		m.focusedInput = MessageInput
		m.messageInput.Focus()
		return m, nil

	case LoginErrorMsg:
		m.error = msg.Error
		return m, nil
	}

	return m, nil
}

func (m Model) View() string {
	switch m.currentView {
	case LoginView:
		return m.renderLoginView()
	case ChatView:
		return m.renderChatView()
	case ChannelListView:
		return m.renderChannelListView()
	default:
		return "Loading..."
	}
}

func (m Model) renderLoginView() string {
	var s strings.Builder

	s.WriteString(" Welcome to TChat\n\n")

	if m.error != nil {
		s.WriteString(fmt.Sprintf("❌ Error: %v\n\n", m.error))
	}

	if m.loading {
		s.WriteString("⏳ Logging in...\n\n")
	}

	s.WriteString("Username:\n")
	s.WriteString(m.usernameInput.View() + "\n\n")

	s.WriteString("Password:\n")
	s.WriteString(m.passwordInput.View() + "\n\n")

	s.WriteString("Press Enter to login, Tab to switch fields, Ctrl+C to quit")

	return s.String()
}

func (m Model) renderChatView() string {
	var s strings.Builder

	// Header
	channelName := m.currentChannel
	if channelName == "" {
		channelName = "General"
	}
	s.WriteString(fmt.Sprintf("💬 Channel: #%s\n", channelName))
	s.WriteString(strings.Repeat("─", m.windowWidth) + "\n")

	// Messages viewport
	s.WriteString(m.viewport.View() + "\n")

	// Input area
	s.WriteString(strings.Repeat("─", m.windowWidth) + "\n")
	s.WriteString(m.messageInput.View())

	return s.String()
}

func (m Model) renderChannelListView() string {
	return m.list.View()
}

// Helper command functions
func (m Model) handleLogin() tea.Cmd {
	username := m.usernameInput.Value()
	password := m.passwordInput.Value()

	// Validate inputs
	if username == "" || password == "" {
		return func() tea.Msg {
			return LoginErrorMsg{Error: fmt.Errorf("username and password are required")}
		}
	}

	return func() tea.Msg {
		// TODO: Replace with actual auth service call
		// For now, simulate successful login
		if username == "demo" && password == "demo" {
			return LoginSuccessMsg{
				Token: "demo-jwt-token",
				User:  models.User{Username: username},
			}
		}
		return LoginErrorMsg{Error: fmt.Errorf("invalid credentials")}
	}
}

func (m Model) handleSendMessage() tea.Cmd {
	message := m.messageInput.Value()
	if message == "" {
		return nil
	}

	// Clear the input
	m.messageInput.SetValue("")

	return func() tea.Msg {
		// TODO: Send message to your backend
		return MessageSentMsg{Content: message}
	}
}

// Custom message types
type LoginSuccessMsg struct {
	Token string
	User  models.User
}

type LoginErrorMsg struct {
	Error error
}

type MessageSentMsg struct {
	Content string
}
