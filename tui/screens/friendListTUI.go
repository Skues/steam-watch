package screens

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"steam-watch/api"

	"github.com/charmbracelet/bubbles/list"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type friendListModel struct {
	details api.DetailedFriendList
	list    list.Model
}

type friendItem struct {
	FriendDetails api.Friend
	FriendSummary api.PlayerSummary
	RecentGames   api.RecentGamesResult
}

func (i friendItem) Title() string {
	return i.FriendSummary.PlayerSummaryResponse.Players[0].PersonaName
}
func (i friendItem) Description() string {
	player := i.FriendSummary.PlayerSummaryResponse.Players[0]
	return fmt.Sprintf("Name: %s\nFriend since: %s\nCurrently: %s\nRelationship: %s\nLast Logoff: %s\n", player.PersonaName, api.UnixToTime(i.FriendDetails.FriendSince), api.PersonaStateStr(player.PersonaState), i.FriendDetails.Relationship, api.UnixToTime(player.LastLogoff))
}

func (i friendItem) FilterValue() string {
	return i.FriendSummary.PlayerSummaryResponse.Players[0].PersonaName
}

func ShowFriendList(steamid string) tea.Model {
	var friendList friendListModel
	fileText, err := os.ReadFile("data/friendList.json")

	if err != nil {
		api.GetOwnedGames("a")
		log.Fatalln("HELLO", err)
		// retrieve data yourself
	} else {
		json.Unmarshal(fileText, &friendList.details)
	}

	items := []list.Item{}
	for _, friend := range friendList.details.DetailedFriendList {
		items = append(items, friendItem{friend.FriendDetails, friend.FriendSummary, friend.RecentGames})

	}
	delegate := list.NewDefaultDelegate()
	delegate.SetSpacing(1)
	delegate.ShowDescription = true
	delegate.SetHeight(5)
	friendList.list = list.New(items, delegate, 0, 0)
	friendList.list.Title = "Friend List:"
	friendList.list.SetSize(50, 20)
	return friendList
}

func (m friendListModel) Init() tea.Cmd {
	return nil

}
func (m friendListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "enter":
			i, ok := m.list.SelectedItem().(friendItem)
			if ok {
				fmt.Println(i.FriendSummary.PlayerSummaryResponse.Players[0].PersonaName)
				return ShowSelectFriend(i), nil
			}
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "b":
			// go back to the main menu
			return m, tea.Quit

		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}
func (m friendListModel) View() string {
	return docStyle.Render(m.list.View())
}
