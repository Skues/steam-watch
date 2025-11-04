package screens

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TextInputModel struct {
	TextInput textinput.Model
	err       error
}

func InitialTIModel() TextInputModel {
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
			return m, m.TextSubmitted

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
	return fmt.Sprintf(
		"Enter your SteamID:\n\n%s\n%s",
		m.TextInput.View(),
		"(esc to quit)",
	) + "\n"
}

func (m TextInputModel) TextSubmitted() tea.Msg {

	return m.TextInput.Value()
}
