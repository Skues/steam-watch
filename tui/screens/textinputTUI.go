package screens

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModel struct {
	TextInput textinput.Model
	err       error
	reason    string
}
type SteamIDInput struct {
	SteamID string
}

func InitialSteamIDInput() TextInputModel {
	ti := textinput.New()
	ti.Placeholder = "111111111111111"
	ti.Focus()
	ti.CharLimit = 120
	ti.Width = 20

	return TextInputModel{
		TextInput: ti,
		err:       nil,
	}

}

func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			valid := m.checkInput()
			if valid {
				return ShowSteamIDOptions(strings.TrimSpace(m.TextInput.Value())), nil
			} else {
				return m, nil
			}

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		m.err = msg
		return m, nil
	}

	m.TextInput, cmd = m.TextInput.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	s := fmt.Sprintf(
		"Enter your SteamID:\n\n%s",
		m.TextInput.View(),
	) + "\n"
	if m.reason != "" {
		s += fmt.Sprintln(m.reason)
	}
	s += fmt.Sprintln("(esc to quit)")
	return s
}
func (m *TextInputModel) checkInput() bool {
	if utf8.RuneCountInString(m.TextInput.Value()) != 17 {
		m.reason = "Invalid length of input, a SteamID is 17 characters long."
		return false
	} else if _, err := strconv.Atoi(m.TextInput.Value()); err != nil {
		m.reason = "Invalid input type, must be a number input."
		return false
	}
	m.reason = ""
	return true

}
