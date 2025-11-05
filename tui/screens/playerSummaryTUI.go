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
}

func ShowPlayerSummary(steamid string) tea.Model {
	var playerSummary playerSummaryModel
	fileText, err := os.ReadFile("data/playerSummary.json")
	if err != nil {
		api.GetPlayerSummary("a")
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
		case "b":
			// go back to the main menu
			return m, tea.Quit
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
