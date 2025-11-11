package tui

import (
	// "fmt"
	"steam-watch/tui/screens"

	tea "github.com/charmbracelet/bubbletea"
)

type RootModel struct {
	current tea.Model
	Input   string
	SteamID string
}

func New(steamid string) *RootModel {
	root := &RootModel{SteamID: steamid}
	if steamid != "" {
		root.current = screens.InitialMainMenu(steamid)
	} else {
		root.current = screens.InitialSteamIDInput()
	}
	return root
}
func (m RootModel) Init() tea.Cmd {
	return m.current.Init()
}
func (m RootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case screens.SteamIDInput:
		m.Input = msg.SteamID
		return screens.ShowSteamIDOptions(m.SteamID), nil
	}

	var cmd tea.Cmd
	m.current, cmd = m.current.Update(msg)
	return m, cmd
}

func (m RootModel) View() string {
	return m.current.View()
}
