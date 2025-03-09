package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

// messageMsg is a custom message type for our TUI
type messageMsg string

type model struct {
	message string
	msgChan chan string
}

func (m model) Init() tea.Cmd {
	// Return a command that listens for messages on the channel
	return listenForMessages(m.msgChan)
}

// listenForMessages is a command that listens for messages on the channel
func listenForMessages(ch chan string) tea.Cmd {
	return func() tea.Msg {
		// Wait for a message on the channel
		msg, ok := <-ch
		if !ok {
			// Channel closed
			return nil
		}
		// Convert the string to our message type
		return messageMsg(msg)
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case messageMsg:
		// Update the message and continue listening
		m.message = string(msg)
		return m, listenForMessages(m.msgChan)
	case tea.KeyMsg:
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, listenForMessages(m.msgChan)
}

func (m model) View() string {
	return fmt.Sprintf("Processing: %s\nPress 'q' to quit.", m.message)
}
