package screens

import (
	"encoding/json"
	"fmt"
	"os"
	"steam-watch/api"

	tea "github.com/charmbracelet/bubbletea"
)

type playerSummaryModel struct {
	details api.PlayerSummary
	steamid string
}

func ShowPlayerSummary(steamid string) tea.Model {
	var playerSummary playerSummaryModel
	playerSummary.steamid = steamid
	fileText, err := os.ReadFile("data/playerSummary.json")
	if err != nil {
		api.GetPlayerSummary(steamid)
		// retrieve data yourself
	} else {
		json.Unmarshal(fileText, &playerSummary.details)

	}
	return playerSummary
}
func (m playerSummaryModel) Init() tea.Cmd {
	return nil

}
func (m playerSummaryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			// go back to the main menu
			return InitialMainMenu(m.steamid), nil
		}
	}

	return m, nil

}
func (m playerSummaryModel) View() string {
	player := m.details.PlayerSummaryResponse.Players[0]
	s := "Player Summary:\n\n"
	s += fmt.Sprintf("Name: %s\n\nCurrently: %v\n\n", player.PersonaName, api.PersonaStateStr(player.PersonaState))
	s += "\n\nPress q to quit"
	return s

}
