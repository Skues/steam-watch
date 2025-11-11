package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

var choices = []string{"Load New SteamID", "Check SteamID", "Delete SteamID"}

type chooseOptionModel struct {
	choice string
	cursor int
}

func InitialChooseModel() chooseOptionModel {
	return chooseOptionModel{}
}
func (m chooseOptionModel) Init() tea.Cmd {
	return nil
}
func (m chooseOptionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ":
			m.choice = choices[m.cursor]
			if m.choice == "Load New SteamID" {
				return InitialSteamIDInput(), nil
			}
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(choices)-1 {
				m.cursor++
			}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
		}
	case string:
		fmt.Println(msg)
		return m, tea.Quit
	}
	return m, nil
}

func (m chooseOptionModel) View() string {
	// The header
	s := "Steam-Watch\n\nChoose an option:\n\n"

	// Iterate over our choices
	for i, choice := range choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
