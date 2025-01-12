package views

import (
	"seesharpsi/qwixx/game_logic"

	"github.com/charmbracelet/lipgloss"
)

var tooltip lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
var hovering lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
var alignment lipgloss.Style

type ViewInfo struct {
	CurrentView string
	MaxPos      [2]uint
}

func Menu(pos [2]uint, width int, height int) (string, ViewInfo) {
	alignment = lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Center, lipgloss.Center)
	var s string
	opts := make(map[uint]string)
	opts[0] = "Play"
	opts[1] = "Stats"
	opts[2] = "Exit"
	spacing := "   "
	opts[pos[0]] = hovering.Render(opts[pos[0]])
	s += opts[0] + spacing + opts[1] + spacing + opts[2]
	//for i, v := range opts {
	//	s += v
	//	if opts[i+1] != "" {
	//		s += spacing
	//	}
	//}
	return alignment.Render(s), ViewInfo{}
}

func PCard(pos [2]int, card game_logic.Card) string {
	return ""
}
