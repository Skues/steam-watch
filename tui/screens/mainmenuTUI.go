package screens

import (
	"encoding/json"
	"fmt"

	// "log"
	"os"
	"steam-watch/api"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var steamid = ""

var (
	headerStyle = lipgloss.NewStyle()
	menuOptions = lipgloss.NewStyle()
	hoverOption = lipgloss.NewStyle().Foreground(lipgloss.Color("168"))
)

type choice struct {
	description string
	title       string
	action      func(steamid string) tea.Model
}

var mainMenuOptions = []choice{{"View Player Summary", "PlayerSummary", func(steamid string) tea.Model { return ShowPlayerSummary(steamid) }}, {"View Friend List", "FriendList", func(steamid string) tea.Model { return ShowFriendList(steamid) }}, {"View Games List", "GamesList", func(steamid string) tea.Model { return ShowGamesList(steamid) }}, {"SteamID Options", "SteamIDOptions", func(steamid string) tea.Model { return ShowSteamIDOptions(steamid) }}}

type mainMenuModel struct {
	choice  choice
	cursor  int
	steamid string
}

func InitialMainMenu(steamidinput string) mainMenuModel {
	SetData(steamidinput)
	return mainMenuModel{steamid: steamidinput}
}
func (m mainMenuModel) Init() tea.Cmd {
	return nil
}
func (m mainMenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter", " ":
			m.choice = mainMenuOptions[m.cursor]
			return m.choice.action(m.steamid), nil

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(mainMenuOptions)-1 {
				m.cursor++
			}
		}
	case string:
		return m, tea.Quit
	}
	return m, nil
}

func (m mainMenuModel) View() string {
	// The header
	s := headerStyle.Render("Steam-Watch\n\nChoose an option:\n\n")

	// Iterate over our choices
	for i, choice := range mainMenuOptions {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Render the row
		if cursor == ">" {
			s += hoverOption.Render(fmt.Sprintf("%s %s", cursor, choice.description)) + "\n"

		} else {
			s += menuOptions.Render(fmt.Sprintf("%s %s", cursor, choice.description)) + "\n"
		}
	}

	// The footer
	s += helpStyle.Render("\nPress q to quit.\n")

	// Send the UI for rendering
	return s
}

func SetData(steamid string) {
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		file, err := os.Create("data/playerSummary.json")
		if err != nil {
			// log.Fatalln(err)
		}
		defer wg.Done()
		defer file.Close()
		playerSummary := api.GetPlayerSummary(steamid)
		encoder := json.NewEncoder(file)
		err = encoder.Encode(playerSummary)
		if err != nil {
			// log.Println(err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

	}()
	go func() {
		file, _ := os.Create("data/friendList.json")
		defer wg.Done()
		defer file.Close()
		friendList := api.FriendListData(steamid)

		encoder := json.NewEncoder(file)
		err := encoder.Encode(friendList)
		if err != nil {
			// log.Println(err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	go func() {
		file, _ := os.Create("data/gamesList.json")
		defer wg.Done()
		defer file.Close()
		ownedGames := api.GetOwnedGames(steamid)
		encoder := json.NewEncoder(file)
		err := encoder.Encode(ownedGames)
		if err != nil {
			// log.Println(err)
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}()
	wg.Wait()
}
