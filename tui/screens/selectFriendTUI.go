package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	selectFriendStyle = lipgloss.NewStyle()
	titleStyle        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("205")).PaddingLeft(2).MarginBottom(1)
	summaryStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("53")).PaddingLeft(4).MarginBottom(0)
	helpStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true).PaddingLeft(2).MarginTop(1)
	borderStyle       = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderLeft(true).BorderForeground(lipgloss.Color("228"))
)

type selectFriend struct {
	friend  friendItem
	steamid string
}

func ShowSelectFriend(steamid string, friend friendItem) tea.Model {
	var selectFriend selectFriend
	selectFriend.steamid = steamid
	selectFriend.friend = friend
	return selectFriend
}
func (m selectFriend) Init() tea.Cmd {
	return nil

}
func (m selectFriend) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			return ShowFriendList(m.steamid), nil
		}
	}

	return m, nil

}
func (m selectFriend) View() string {
	tea.ClearScreen()
	player := m.friend.FriendSummary.PlayerSummaryResponse.Players[0]

	var gameInfo string
	if player.GameExtraInfo == "" {
		gameInfo = "Not in game"
	} else {
		gameInfo = player.GameExtraInfo
	}

	s := lipgloss.JoinVertical(lipgloss.Left, titleStyle.Render(fmt.Sprintf("%s Summary:", player.PersonaName)),
		summaryStyle.Render(fmt.Sprintf("Game Status: %s\nTotal Game Count: %v", gameInfo, m.friend.RecentGames.TotalCount)),
		helpStyle.Render("Press q to quit"))

	content := borderStyle.Render(s)
	return content

}
