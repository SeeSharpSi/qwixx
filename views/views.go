package views

import (
	"fmt"
	"seesharpsi/qwixx/game_logic"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

var tooltip lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
var hovering lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
var cardhovering lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("5")).Background(lipgloss.Color("0")).Bold(true)
var red lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Background(lipgloss.Color("0")).Bold(true)
var yellow lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("11")).Background(lipgloss.Color("0")).Bold(true)
var green lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("10")).Background(lipgloss.Color("0")).Bold(true)
var blue lipgloss.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Background(lipgloss.Color("0")).Bold(true)
var background lipgloss.Style = lipgloss.NewStyle().Background(lipgloss.Color("0"))
var bold lipgloss.Style = lipgloss.NewStyle().Background(lipgloss.Color("0")).Bold(true)
var cardstyle lipgloss.Style = lipgloss.NewStyle().
	Background(lipgloss.Color("0")).
	Padding(2).
	Border(lipgloss.OuterHalfBlockBorder()).
	BorderBackground(lipgloss.Color("0")).
	Margin(2)

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
	opts := [4]opt{}
	opts[0] = opt{
		display: "Join",
		val:     "gameselect",
	}
	opts[1] = opt{
		display: "Create",
		val:     "creategame",
	}
	opts[2] = opt{
		display: "Stats",
		val:     "stats",
	}
	opts[3] = opt{
		display: "Exit",
		val:     "exit",
	}
	return opts[pos[1]].val, [2]uint{0, 2}
}

func MenuRender(pos [2]uint, width int, height int) string {
	alignment := lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Center, lipgloss.Center)
	var s string
	opts := [4]opt{}
	opts[0] = opt{
		display: "Join",
		val:     "gameselect",
	}
	opts[1] = opt{
		display: "Create",
		val:     "creategame",
	}
	opts[2] = opt{
		display: "Stats",
		val:     "stats",
	}
	opts[3] = opt{
		display: "Exit",
		val:     "exit",
	}
	spacing := "   "
	opts[pos[1]].display = hovering.Render(opts[pos[1]].display)
	s += opts[0].display + spacing + opts[1].display + spacing + opts[2].display
	return alignment.Render(s)
}

func CardRender(pos [2]uint, width int, height int, card game_logic.Card) string {
	var s string

	s += "Red:    "
	for i, v := range card.Red {
		if v {
			if pos == [2]uint{0, uint(i)} {
				s += cardhovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += red.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{0, uint(i)} {
				s += cardhovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += red.Render(fmt.Sprintf("[%v]", i+2))
			}
		}
	}
	s += "\nYellow: "
	for i, v := range card.Yellow {
		if v {
			if pos == [2]uint{1, uint(i)} {
				s += cardhovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += yellow.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{1, uint(i)} {
				s += cardhovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += yellow.Render(fmt.Sprintf("[%v]", i+2))
			}
		}
	}
	s += "\nGreen:  "
	for i, v := range card.Green {
		if v {
			if pos == [2]uint{2, uint(i)} {
				s += cardhovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += green.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{2, uint(i)} {
				s += cardhovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += green.Render(fmt.Sprintf("[%v]", i+2))
			}
		}
	}
	s += "\nBlue:   "
	for i, v := range card.Blue {
		if v {
			if pos == [2]uint{3, uint(i)} {
				s += cardhovering.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += blue.Strikethrough(true).Render(fmt.Sprintf("[%v]", i+2))
			}
		} else {
			if pos == [2]uint{3, uint(i)} {
				s += cardhovering.Render(fmt.Sprintf("[%v]", i+2))
			} else {
				s += blue.Render(fmt.Sprintf("[%v]", i+2))
			}
		}

	}
	// opts[pos[1]].display = hovering.Render(opts[pos[1]].display)
	// s += opts[0].display + spacing + opts[1].display + spacing + opts[2].display
	s = cardstyle.Render(s)
	return s
}

func DiceRender(dice game_logic.Dice) string {
	var s string
	s = "DICE"
	s += "\n" + bold.Render("|") + bold.Render(fmt.Sprint(dice.White1)+bold.Render("|"))
	s += bold.Render(fmt.Sprint(dice.White2)) + bold.Render("|")
	s += red.Render(fmt.Sprint(dice.Red)) + bold.Render("|")
	s += yellow.Render(fmt.Sprint(dice.Yellow)) + bold.Render("|")
	s += green.Render(fmt.Sprint(dice.Green)) + bold.Render("|")
	s += blue.Render(fmt.Sprint(dice.Blue)) + bold.Render("|")
	return cardstyle.Align(lipgloss.Center).Render(s)
}

func CardInfo(pos [2]uint) (hovering string, maxPos [2]uint) {
	return "card", [2]uint{3, 10}
}

func GameSelectRender(games map[string]*game_logic.Game, ti textinput.Model, width int, height int) string {
	alignment := lipgloss.NewStyle().Width(width).Height(height).Align(lipgloss.Left, lipgloss.Center)
	var s string
	_, ok := games[ti.Value()]
	if ok {
		s += alignment.Render(green.Render(ti.View()))
	} else {
		s += alignment.Render(ti.View())
	}
	return s
}
