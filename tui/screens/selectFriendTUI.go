package screens

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type selectFriend struct {
	friend friendItem
}

func ShowSelectFriend(friend friendItem) tea.Model {
	var selectFriend selectFriend
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
			return ShowFriendList("asd"), nil
		}
	}

	return m, nil

}
func (m selectFriend) View() string {
	player := m.friend.FriendSummary.PlayerSummaryResponse.Players[0]
	s := fmt.Sprintf("\n%s Summary:\n", player.PersonaName)
	s += fmt.Sprintf("Game Status: %s\nTotal Game Count: %v\n", player.GameExtraInfo, m.friend.RecentGames.TotalCount)
	s += "\n\nPress q to quit"
	return s

}
