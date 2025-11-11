package screens_test

import (
	tea "github.com/charmbracelet/bubbletea"
	"steam-watch/tui/screens"
	"testing"
)

func TestShowSteamIDOptions(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		steamidInput string
		want         tea.Model
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := screens.ShowSteamIDOptions(tt.steamidInput)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("ShowSteamIDOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}
