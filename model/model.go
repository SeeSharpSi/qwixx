package model

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"seesharpsi/qwixx/game_logic"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	"github.com/muesli/termenv"
)

type box struct {
	pos [][]int
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "enter":
			m.App.send("hello")
		case "j":
			m.Hovering[0] += 1
		case "k":
			m.Hovering[0] -= 1
		case "l":
			m.Hovering[1] += 1
		case "h":
			m.Hovering[1] -= 1
		}
	case string:
		print(msg)
	}
	return m, nil
}

type Model struct {
	CurrentView string
	Hovering    [2]int
	*App
	Styles
	Turn     string
	Player   string
	Term     string
	Profile  string
	Width    int
	Height   int
	Bg       string
	Messages []string
	Id       string
	Err      error
}

type Styles struct {
	Viewport      viewport.Model
	Textarea      textarea.Model
	SenderStyle   lipgloss.Style
	TxtStyle      lipgloss.Style
	QuitStyle     lipgloss.Style
	ToolTipStyle  lipgloss.Style
	HoveringStyle lipgloss.Style
}

type Cursor struct {
	Box [][]string
}

func (m Model) Init() tea.Cmd {
	return nil
}

func InitialModel() Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`placeholder text`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return Model{
		CurrentView: "menu",
		Hovering:    [2]int{2, 2},
		Styles: Styles{
			Textarea:      ta,
			Viewport:      vp,
			QuitStyle:     lipgloss.NewStyle().Foreground(lipgloss.Color("11")),
			ToolTipStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
			HoveringStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		},
		Messages: []string{},
		Err:      nil,
	}
}

func (m Model) View() string {
	var s string
	var posX = m.Hovering[0] % 10
	var posY = m.Hovering[1] % 10
	if m.CurrentView == "menu" {
		for i := range 10 {
			for j := range 10 {
				if posX == i && posY == j {
					s += m.HoveringStyle.Render("X")
				} else {
					s += "O"
				}
			}
			s += "\n"
		}
	} else {
		s = fmt.Sprintf("Your term is %s\nYour window size is %dx%d\nBackground: %s\nColor Profile: %s", m.Term, m.Width, m.Height, m.Bg, m.Profile)
	}
	return m.TxtStyle.Render(s) + "\n\n" + m.ToolTipStyle.Render("Press C to see card") + "\n\n" + m.QuitStyle.Render("Press 'q' to quit\n")
}

//app stuff

const (
	host = "localhost"
	port = "23234"
)

// app contains a wish server and the list of running programs.
type App struct {
	game_logic.Game
	*ssh.Server
	progs []*tea.Program
}

// I need the app to send a msg that adds a card and player when someone joins the game
// or actually, when a message is sent it also sends "player" attached; it said player doesn't exist card is appended

// send dispatches a message to all running programs.
func (a *App) send(msg tea.Msg) {
	for _, p := range a.progs {
		go p.Send(msg)
	}
}

func (a *App) send2(msg string) {
	for _, p := range a.progs {
		go p.Send(msg)
	}
}

func (a *App) Start() {
	var err error
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = a.ListenAndServe(); err != nil {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := a.Shutdown(ctx); err != nil {
		log.Error("Could not stop server", "error", err)
	}
}

func NewApp() *App {
	a := new(App)
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath("./id_ed25519"),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(a.ProgramHandler, termenv.ANSI256),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	a.Server = s
	return a
}

func (a *App) ProgramHandler(s ssh.Session) *tea.Program {
	model := InitialModel()
	model.App = a
	model.Player = s.User()

	p := tea.NewProgram(model, bubbletea.MakeOptions(s)...)
	a.progs = append(a.progs, p)

	return p
}
