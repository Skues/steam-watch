package screens

import tea "github.com/charmbracelet/bubbletea"

type friendListModel struct {
}

func ShowFriendList(steamid string) tea.Model {
	return friendListModel{}
}

func (m friendListModel) Init() tea.Cmd {
	return nil

}
func (m friendListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
func (m friendListModel) View() string {
	return "Nah"

}
