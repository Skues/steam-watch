package screens

import (
	"fmt"
	"steam-watch/api"

	tea "github.com/charmbracelet/bubbletea"
)

type steamIDCheck struct {
	steamid string
	valid   bool
}

func ShowSteamIDOptions(steamid string) tea.Model {
	return steamIDCheck{steamid: steamid, valid: api.ValidateSteamID(steamid)}
}

func (m steamIDCheck) Init() tea.Cmd {
	return nil
}

func (m steamIDCheck) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:

		switch msg.Type {
		case tea.KeyEnter:
			return m, tea.Quit

		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case error:
		return m, nil
	}
	return m, nil
}
func (m steamIDCheck) View() string {
	s := "Steam-Watch\n\n"
	s += fmt.Sprintf("SteamID: %s\n\nValid: %v", m.steamid, m.valid)
	return s
}
