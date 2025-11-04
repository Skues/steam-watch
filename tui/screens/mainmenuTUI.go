package screens

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
)

var steamid = ""

type choice struct {
	description string
	title       string
	action      func() tea.Model
}

var mainMenuOptions = []choice{{"View Player Summary", "PlayerSummary", func() tea.Model { return ShowPlayerSummary(steamid) }}, {"View Friend List", "FriendList", func() tea.Model { return ShowFriendList(steamid) }}, {"View Games List", "GamesList", func() tea.Model { return ShowGamesList(steamid) }}, {"SteamID Options", "SteamIDOptions", func() tea.Model { return ShowSteamIDOptions("test") }}}

type mainMenuModel struct {
	choice choice
	cursor int
}

func InitialMainMenu(steamidinput string) mainMenuModel {
	steamid = steamidinput
	return mainMenuModel{}
}
func (m mainMenuModel) Init() tea.Cmd {
	return nil
}
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ":
			m.choice = mainMenuOptions[m.cursor]
			return m.choice.action(steamid), nil

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(mainMenuOptions)-1 {
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

func (m mainMenuModel) View() string {
	// The header
	s := "Steam-Watch\n\nChoose an option:\n\n"

	// Iterate over our choices
	for i, choice := range mainMenuOptions {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice.description)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}
