package screens

import (
	"encoding/json"
	"fmt"
	"os"
	"steam-watch/api"

	tea "github.com/charmbracelet/bubbletea"
)

type playerSummaryModel struct {
	details api.DetailedFriend
}

func ShowPlayerSummary(steamid string) tea.Model {
	var playerSummary playerSummaryModel
	fileText, err := os.ReadFile("data/playerSummary.json")
	if err != nil {
		api.GetPlayerSummary()
		// retrieve data yourself
	} else {
		json.Unmarshal(fileText, &playerSummary)

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
			return InitialMainMenu(), nil
		}
	}

	return m, nil

}
func (m playerSummaryModel) View() string {
	s := "Player Summary:\n\n"
	s += fmt.Sprintf("Name:%s\n\nCurrently:%v\n\n", m.details.FriendSummary.PlayerSummaryResponse.Players[0].PersonaName, m.details.FriendSummary.PlayerSummaryResponse.Players[0].CommunityVisibilityState)
	s += "\n\nPress q to quit"
	return s

}
