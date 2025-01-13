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
	"seesharpsi/qwixx/views"

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

type Model struct {
	views.ViewInfo
	Pos [2]uint
	*App
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
	Border        lipgloss.Style
	Alignment     lipgloss.Style
	SenderStyle   lipgloss.Style
	TxtStyle      lipgloss.Style
	QuitStyle     lipgloss.Style
	ToolTipStyle  lipgloss.Style
	HoveringStyle lipgloss.Style
}

func (m Model) Init() tea.Cmd {
	return nil
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
			fmt.Printf("\nhovering: %+v\n", m.ViewInfo)
		case "j":
			m.Pos[0] += 1
			if m.Pos[0] > m.MaxPos[0] {
				m.Pos[0] -= 1
			}
			m.CurrentView, m.MaxPos = views.MenuInfo(m.Pos)
		case "k":
			if m.Pos[0] > 0 {
				m.Pos[0] -= 1
			}
			m.CurrentView, m.MaxPos = views.MenuInfo(m.Pos)
		case "l":
			m.Pos[1] += 1
			if m.Pos[1] > m.MaxPos[1] {
				m.Pos[1] -= 1
			}
			m.CurrentView, m.MaxPos = views.MenuInfo(m.Pos)
		case "h":
			if m.Pos[1] > 0 {
				m.Pos[1] -= 1
			}
			m.CurrentView, m.MaxPos = views.MenuInfo(m.Pos)
		}
	case string:
		print(msg)
	case views.ViewInfo:
		m.ViewInfo = msg
	}
	return m, nil
}

func (m Model) View() string {
	var s string = " "
	switch m.ViewInfo.CurrentView {
	case "menu":
		fmt.Println("menu time")
		s = views.MenuRender(m.Pos, m.Width, m.Height)
	case "stats":
		fmt.Println("stats time")
		s = views.MenuRender(m.Pos, m.Width, m.Height)
	case "board":
		fmt.Println("board time")
		s = views.MenuRender(m.Pos, m.Width, m.Height)
	case "exit":
		fmt.Println("exit time")
		s = views.MenuRender(m.Pos, m.Width, m.Height)
	default:
		fmt.Println("hovering unknown value")
		m.App.send(views.ViewInfo{
			MaxPos:      [2]uint{0, 2},
			CurrentView: "menu",
		})
	}
	// s += "\nposX: " + fmt.Sprint(posX) + "\nposY: " + fmt.Sprint(posY)
	s += "\nm.Pos[0]: " + fmt.Sprint(m.Pos[0]) + "\nm.Pos[1]: " + fmt.Sprint(m.Pos[1])
	return s
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
	Pty, _, _ := s.Pty()
	model := Model{
		Player:   s.User(),
		App:      a,
		Pos:      [2]uint{0, 0},
		Width:    Pty.Window.Width,
		Height:   Pty.Window.Height,
		Messages: []string{},
		Err:      nil,
		ViewInfo: views.ViewInfo{
			CurrentView: "menu",
			MaxPos:      [2]uint{0, 2},
		},
	}

	p := tea.NewProgram(model, bubbletea.MakeOptions(s)...)
	a.progs = append(a.progs, p)

	return p
}

// returns the terminal width * a decimal percentage
func (m Model) screenpadX(pad float64) int {
	tmp := int(float64(m.Width) * pad)
	print("\npadX: ", tmp)
	return tmp
}

// returns the terminal height * a decimal percentage
func (m Model) screenpadY(pad float64) int {
	tmp := int(float64(m.Height) * pad)
	print("\npadY: ", tmp)
	return tmp
}
