package screens

import tea "github.com/charmbracelet/bubbletea"

type gameListModel struct {
}

func ShowGamesList(steamid string) tea.Model {
	return gameListModel{}
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
		case "b":
			// go back to the main menu
			return InitialMainMenu(), nil
		}
	}
	return m, nil

}
func (m gameListModel) View() string {
	return "Nah"

}
