package tui

import (
	// "fmt"
	"steam/code/tui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	current tea.Model
	Input   string
}

func New() RootModel {
	return RootModel{current: screens.InitialChooseModel()}
}
func (m RootModel) Init() tea.Cmd {
	return m.current.Init()
}
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case string:
		m.Input = msg
		return screens.InitialSteamID(msg), nil
		// fmt.Println(m.input)
		return m, nil
	}

	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return m, cmd
}

func (m RootModel) View() string {
	return m.current.View() + "\n\nLast Input: " + m.Input
}
