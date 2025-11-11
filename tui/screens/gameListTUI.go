package screens

import (
	"encoding/json"
	"fmt"

	// "fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	"steam-watch/api"
)

var gameListDocStyle = lipgloss.NewStyle().Margin(1, 2)

type gameListModel struct {
	details api.RecentGames
	list    list.Model
	steamid string
}
type gameItem struct {
	AppID           int    `json:"appid"`
	Name            string `json:"name"`
	Playtime2Week   int    `json:"playtime_2weeks"`
	PlaytimeForever int    `json:"playtime_forever"`
	ImgIcon         string `json:"img_icon_url"`
}

func (i gameItem) Title() string { return i.Name }
func (i gameItem) Description() string {
	return fmt.Sprintf("Playtime 2 week: %v\nPlaytime Overall:%v", i.Playtime2Week/60, i.PlaytimeForever/60)
}
func (i gameItem) FilterValue() string { return i.Name }

func ShowGamesList(steamid string) tea.Model {
	var gameList gameListModel
	gameList.steamid = steamid
	fileText, err := os.ReadFile("data/gamesList.json")
	if err != nil {
		api.GetOwnedGames(steamid)
		// retrieve data yourself
	} else {
		json.Unmarshal(fileText, &gameList.details)

	}
	items := []list.Item{}
	for _, game := range gameList.details.RecentGamesResponse.Games {
		items = append(items, gameItem{game.AppID, game.Name, game.Playtime2Week, game.PlaytimeForever, game.ImgIcon})
	}
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(1)
	delegate.ShowDescription = true
	delegate.SetHeight(3)
	gameList.list = list.New(items, delegate, 0, 0)
	gameList.list.Title = "All Owned Games List"
	gameList.list.SetSize(50, 20)
	return gameList
}
func (m gameListModel) Init() tea.Cmd {
	return nil

}
func (m gameListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
	case tea.WindowSizeMsg:
		h, v := gameListDocStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd

}
func (m gameListModel) View() string {
	return gameListDocStyle.Render(m.list.View())
	// s := "Game List:\n\n"
	// for i, game := range m.details.RecentGamesResponse.Games {
	// 	s += fmt.Sprintf("\n%v\n%s:\nPlaytime 2 Weeks: %v hours\nPlaytime Overall: %v hours\n", i+1, game.Name, game.Playtime2Week/60, game.PlaytimeForever/60)
	// }
	//
	// s += fmt.Sprintln("\n\nPress q to quit.")
	// return s

}
