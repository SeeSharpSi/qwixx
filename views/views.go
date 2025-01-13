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
	Hovering    string
}

type opt struct {
	display string
	val     string
}

func MenuInfo(pos [2]uint) (hovering string, maxPos [2]uint) {
	opts := [3]opt{}
	opts[0] = opt{
		display: "Play",
		val:     "board",
	}
	opts[1] = opt{
		display: "Stats",
		val:     "stats",
	}
	opts[2] = opt{
		display: "Exit",
		val:     "exit",
	}
	return opts[pos[1]].val, [2]uint{0, 2}
}

func MenuRender(pos [2]uint, width int, height int) string {
	alignment = lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Center, lipgloss.Center)
	var s string
	opts := [3]opt{}
	opts[0] = opt{
		display: "Play",
		val:     "board",
	}
	opts[1] = opt{
		display: "Stats",
		val:     "stats",
	}
	opts[2] = opt{
		display: "Exit",
		val:     "exit",
	}
	spacing := "   "
	opts[pos[1]].display = hovering.Render(opts[pos[1]].display)
	s += opts[0].display + spacing + opts[1].display + spacing + opts[2].display
	return alignment.Render(s)
}

func Stats(pos [2]uint, width int, height int) (string, ViewInfo) {
	return "", ViewInfo{}
}

func PCard(pos [2]int, card game_logic.Card) string {
	return ""
}
