package model

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"seesharpsi/qwixx/game_logic"
	"seesharpsi/qwixx/views"

	"github.com/charmbracelet/bubbles/textinput"
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

type user struct {
	Player  string
	Gamekey string
}

type Model struct {
	*App
	views.ViewInfo
	Game *game_logic.Game
	Pos  [2]uint
	Turn string
	// remaining skips
	Player   string
	Term     string
	Profile  string
	Width    int
	Height   int
	Bg       string
	Messages []string
	Id       string
	Err      error
	Styles
}

type Styles struct {
	Viewport  viewport.Model
	TextInput textinput.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
		m.Width = msg.Width
	case tea.KeyMsg:
		switch m.ViewInfo.CurrentView {
		case "menu":
			m, cmd = m.menuUpdate(msg)
		case "card":
			m, cmd = m.cardUpdate(msg)
		case "gameselect":
			m, cmd = m.gameSelectUpdate(msg)
		}
	case string:
		print(msg)
	case views.ViewInfo:
		m.ViewInfo = msg
	}
	return m, cmd
}

func (m Model) menuUpdate(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		m.ViewInfo.CurrentView, _ = views.MenuInfo(m.Pos)
		if m.ViewInfo.CurrentView == "creategame" {
			m.Game, m.Game.Key = m.App.createNewGame(m)
			m.Pos = [2]uint{0, 0}
			m.ViewInfo.CurrentView = "card"
		}
	case "j":
		m.Pos[0] += 1
		if m.Pos[0] > m.MaxPos[0] {
			m.Pos[0] -= 1
		}
	case "k":
		if m.Pos[0] > 0 {
			m.Pos[0] -= 1
		}
	case "l":
		m.Pos[1] += 1
		if m.Pos[1] > m.MaxPos[1] {
			m.Pos[1] -= 1
		}
	case "h":
		if m.Pos[1] > 0 {
			m.Pos[1] -= 1
		}
	}
	return m, nil
}

func (m Model) cardUpdate(msg tea.KeyMsg) (Model, tea.Cmd) {
	_, m.MaxPos = views.CardInfo(m.Pos)
	var b bool
	turn := false
	if m.Player == m.Game.Players[m.Game.Turn] {
		turn = true
	}
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "enter":
		switch m.Pos[0] {
		case 0:
			playerCard := m.Game.Cards[m.Player]
			playerCard.Red, b = playerCard.Red.TryMark(int(m.Pos[1]), turn, [3]game_logic.Die{m.Game.Dice.White1, m.Game.Dice.White2, m.Game.Dice.Red})
			if !b {
				return m, nil
			}
			m.Game.Gone[m.Player] = true
			m.Game.Cards[m.Player] = playerCard
			if m.Game.TurnCheck() {
				m.Game.ResetTurns()
				m.Game.Dice.Roll()
				m.App.send("")
			}
		case 1:
			playerCard := m.Game.Cards[m.Player]
			playerCard.Yellow, b = playerCard.Yellow.TryMark(int(m.Pos[1]), turn, [3]game_logic.Die{m.Game.Dice.White1, m.Game.Dice.White2, m.Game.Dice.Yellow})
			if !b {
				return m, nil
			}
			m.Game.Gone[m.Player] = true
			m.Game.Cards[m.Player] = playerCard
			if m.Game.TurnCheck() {
				m.Game.ResetTurns()
				m.Game.Dice.Roll()
				m.App.send("")
			}
		case 2:
			playerCard := m.Game.Cards[m.Player]
			playerCard.Green, b = playerCard.Green.TryMark(int(m.Pos[1]), turn, [3]game_logic.Die{m.Game.Dice.White1, m.Game.Dice.White2, m.Game.Dice.Green})
			if !b {
				return m, nil
			}
			m.Game.Gone[m.Player] = true
			m.Game.Cards[m.Player] = playerCard
			if m.Game.TurnCheck() {
				m.Game.ResetTurns()
				m.Game.Dice.Roll()
				m.App.send("")
			}
		case 3:
			playerCard := m.Game.Cards[m.Player]
			playerCard.Blue, b = playerCard.Blue.TryMark(int(m.Pos[1]), turn, [3]game_logic.Die{m.Game.Dice.White1, m.Game.Dice.White2, m.Game.Dice.Blue})
			if !b {
				return m, nil
			}
			m.Game.Gone[m.Player] = true
			m.Game.Cards[m.Player] = playerCard
			if m.Game.TurnCheck() {
				m.Game.ResetTurns()
				m.Game.Dice.Roll()
				m.App.send("")
			}
		}
	case "p":
		b := m.TrySkip()
		if !b {
			return m, nil
		}
		m.Game.Gone[m.Player] = true
		if m.Game.TurnCheck() {
			m.Game.ResetTurns()
			m.Game.Dice.Roll()
			m.App.send("")
		}
	case "j":
		m.Pos[0] += 1
		if m.Pos[0] > m.MaxPos[0] {
			m.Pos[0] -= 1
		}
	case "k":
		if m.Pos[0] > 0 {
			m.Pos[0] -= 1
		}
	case "l":
		m.Pos[1] += 1
		if m.Pos[1] > m.MaxPos[1] {
			m.Pos[1] -= 1
		}
	case "h":
		if m.Pos[1] > 0 {
			m.Pos[1] -= 1
		}
	}
	return m, nil
}

func (m Model) gameSelectUpdate(msg tea.KeyMsg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg.String() {
	case "enter":
		v, ok := m.App.Games[m.TextInput.Value()]
		if ok {
			m.Game = v
			m.CurrentView = "card"
			m.Game.Players = append(m.Game.Players, m.Player)
			m.App.users = append(m.App.users, user{m.Player, m.Game.Key})
			m.Game.Skips[m.Player] = 4
		} else {
		}
	case "ctrl+c":
		return m, tea.Quit
	default:
		m.TextInput, cmd = m.TextInput.Update(msg)

	}
	return m, cmd
}

func (m Model) View() string {
	var s string = " "
	turn := false
	switch m.ViewInfo.CurrentView {
	case "menu":
		s = views.MenuRender(m.Pos, m.Width, m.Height)
	case "gameselect":
		s = views.GameSelectRender(m.App.Games, m.TextInput, m.Width, m.Height)
	case "stats":
		s = views.MenuRender(m.Pos, m.Width, m.Height)
	case "card":
		if m.Player == m.Game.Players[m.Game.Turn] {
			turn = true
		}
		s = views.DiceRender(m.Game.Dice, turn)
		s += "\n" + views.CardRender(m.Pos, m.Width, m.Height, m.Game.Skips[m.Player], m.Game.Cards[m.Player])
		s += "\n" + m.Game.Key
	case "exit":
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
	alignment := lipgloss.NewStyle().Width(m.Width).Height(m.Height).Align(lipgloss.Center, lipgloss.Center)
	return alignment.Render(s)
}

//app stuff

const (
	host = "localhost"
	port = "23234"
)

// app contains a wish server and the list of running programs.
type App struct {
	Games map[string]*game_logic.Game
	*ssh.Server
	progs []*tea.Program
	users []user
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
	a.Games = make(map[string]*game_logic.Game)
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath("./id_ed25519"),
		wish.WithMiddleware(
			handleDisconnectMiddleware(a),
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

func handleDisconnectMiddleware(a *App) wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(s ssh.Session) {
			// Get the session context
			ctx := s.Context()

			// Run the disconnection handler in a goroutine
			go func() {
				select {
				case <-ctx.Done():
					// The user disconnected
					handleDisconnect(s, a)
				}
			}()

			// Call the next handler in the chain
			next(s)
		}
	}
}

// Log or handle a disconnection
func handleDisconnect(s ssh.Session, a *App) {
	user := s.User()
	for _, v := range a.users {
		if v.Player == user {
			game := a.Games[v.Gamekey]
			if game != nil {
				game.Players = remove(game.Players, user)
				delete(game.Gone, user)
				a.Games[v.Gamekey] = game
			}
		}
	}

	for _, v := range a.users {
		game := a.Games[v.Gamekey]
		if game != nil {
			if len(game.Players) == 0 {
				delete(a.Games, v.Gamekey)
			}
		}
	}
}

func remove(l []string, item string) []string {
	for i, other := range l {
		if other == item {
			return append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func (a *App) ProgramHandler(s ssh.Session) *tea.Program {
	Pty, _, _ := s.Pty()
	ti := textinput.New()
	ti.Placeholder = "Enter game id..."
	ti.Focus()
	ti.Width = 20
	ti.CharLimit = 4
	player := s.User()
	model := Model{
		Player:   player,
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
		Styles: Styles{
			TextInput: ti,
		},
		Game: &game_logic.Game{},
	}

	p := tea.NewProgram(model, bubbletea.MakeOptions(s)...)
	a.progs = append(a.progs, p)
	return p
}

func (m *Model) TrySkip() bool {
	if m.Game.Skips[m.Player] > 0 {
		m.Game.Skips[m.Player]--
		return true
	}
	return false
}

func (a *App) createNewGame(m Model) (*game_logic.Game, string) {
	player := m.Player
	init_row := [11]bool{false, false, false, false, false, false, false, false, false, false, false}
	card := game_logic.Card{
		Player: player,
		Red:    init_row,
		Yellow: init_row,
		Green:  init_row,
		Blue:   init_row,
		Skips:  4,
	}
	cards := make(map[string]game_logic.Card)
	cards[player] = card
	skips := make(map[string]int)
	skips[player] = 4
	gamekey := keygen(4)
	game := game_logic.Game{
		Key:     gamekey,
		Cards:   cards,
		Players: []string{player},
		Gone:    make(map[string]bool),
		Skips:   skips,
		Turn:    0,
	}
	game.Dice.Roll()
	a.Games[gamekey] = &game
	m.Game = &game
	m.Game.ResetTurns()
	m.App.users = append(m.App.users, user{m.Player, m.Game.Key})
	return &game, gamekey
}

func keygen(n int) string {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)

	var src = rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}
