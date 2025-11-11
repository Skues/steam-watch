package screens

import (
	"fmt"
	"steam-watch/api"
	"steam-watch/data"

	tea "github.com/charmbracelet/bubbletea"
)

type steamIDCheck struct {
	steamid string
	valid   bool
	choice  string
	cursor  int
}

var steamIDChoices = []string{"Go to Main Menu", "Load New SteamID", "Check SteamID", "Delete SteamID"}

func ShowSteamIDOptions(steamidInput string) tea.Model {
	tea.ClearScreen()

	if api.ValidateSteamID(steamidInput) {
		data.WriteSteamID(steamidInput)
		return steamIDCheck{steamid: steamidInput, valid: true}
	} else {
		return steamIDCheck{steamid: "", valid: false}
	}
}

func (m steamIDCheck) Init() tea.Cmd {
	return nil
}

func (m steamIDCheck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.String() {
		case "enter":
			if m.valid {
				switch chosen := steamIDChoices[m.cursor]; chosen {
				case "Go to Main Menu":
					return InitialMainMenu(m.steamid), nil
				case "Load New SteamID":
				case "Check SteamID":
				case "Delete SteamID":
				}
			} else {
				return InitialSteamIDInput(), nil
			}
		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(steamIDChoices)-1 {
				m.cursor++
			}

			// The "enter" key and the spacebar (a literal space) toggle
			// the selected state for the item that the cursor is pointing at.
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		return m, nil
	}
	return m, nil
}
func (m steamIDCheck) View() string {
	s := headerStyle.Render("Steam-Watch\n\n")
	if m.valid {
		for i, choice := range steamIDChoices {

			// Is the cursor pointing at this choice?
			if m.cursor == i {
				s += hoverOption.Render(fmt.Sprintf("%s %s", ">", choice)) + "\n"
			} else {
				s += menuOptions.Render(fmt.Sprintf("%s %s", " ", choice)) + "\n"
			}
		}
		s += helpStyle.Render(fmt.Sprintln("Press Enter to continue."))
	} else {
		s += fmt.Sprintln("You have entered an invalid SteamID.")
		s += fmt.Sprintln("Press Enter to try again.")
	}
	return s
}
