package views

import (
	"fmt"
	"seesharpsi/qwixx/game_logic"

	"github.com/charmbracelet/lipgloss"
)

var tooltip lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
var hovering lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Background(lipgloss.Color("0"))
var red lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Background(lipgloss.Color("0"))
var yellow lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Background(lipgloss.Color("0"))
var green lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Background(lipgloss.Color("0"))
var blue lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("0"))
var background lipgloss.Style = lipgloss.NewStyle().Background(lipgloss.Color("0"))
var cardstyle lipgloss.Style = lipgloss.NewStyle().Bold(true).Background(lipgloss.Color("0")).Padding(2).Border(lipgloss.OuterHalfBlockBorder()).BorderBackground(lipgloss.Color("0"))
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

func CardRender(pos [2]uint, width int, height int, card game_logic.Card) string {
	alignment = lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Center, lipgloss.Center)
	var s string

	s += "Red:    "
	for i, v := range card.Red {
		if v {
			if pos == [2]uint{0, uint(i)} {
				s += hovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += red.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{0, uint(i)} {
				s += hovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += red.Render(fmt.Sprintf("[%v]", i+2))
			}
		}
	}
	s += "\nYellow: "
	for i, v := range card.Yellow {
		if v {
			if pos == [2]uint{1, uint(i)} {
				s += hovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += yellow.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{1, uint(i)} {
				s += hovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += yellow.Render(fmt.Sprintf("[%v]", i+2))
			}
		}
	}
	s += "\nGreen:  "
	for i, v := range card.Green {
		if v {
			if pos == [2]uint{2, uint(i)} {
				s += hovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += green.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{2, uint(i)} {
				s += hovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += green.Render(fmt.Sprintf("[%v]", i+2))
			}
		}
	}
	s += "\nBlue:   "
	for i, v := range card.Blue {
		if v {
			if pos == [2]uint{3, uint(i)} {
				s += hovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += blue.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{3, uint(i)} {
				s += hovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += blue.Render(fmt.Sprintf("[%v]", i+2))
			}
		}

	}
	// opts[pos[1]].display = hovering.Render(opts[pos[1]].display)
	// s += opts[0].display + spacing + opts[1].display + spacing + opts[2].display
	s = cardstyle.Render(s)
	s = alignment.Render(s)
	return s
}

func CardInfo(pos [2]uint) (hovering string, maxPos [2]uint) {
	return "card", [2]uint{3, 10}
}

func Stats(pos [2]uint, width int, height int) (string, ViewInfo) {
	return "", ViewInfo{}
}

func PCard(pos [2]int, card game_logic.Card) string {
	return ""
}
