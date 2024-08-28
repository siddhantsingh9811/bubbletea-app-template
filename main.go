package main

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

var (
	helpStyle               = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	bunnyStyle              = lipgloss.NewStyle().Foreground(lipgloss.Color("#00AFFF"))
	bannerStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("#5F5FD7"))
	unselectedListStyle     = lipgloss.NewStyle().Margin(1, 1).Border(lipgloss.RoundedBorder())
	unselectedViewportStyle = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1).Margin(1, 2, 0, 0)
	selectedListStyle       = lipgloss.NewStyle().Margin(1, 1).Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#00AFFF"))
	selectedViewportStyle   = lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).Padding(0, 1).Margin(1, 2, 0, 0).BorderForeground(lipgloss.Color("#00AFFF"))
)

const (
	host = "localhost"
	port = "23234"
)

type sessionState uint
type page uint
type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type model struct {
	state    sessionState
	page     page
	list     list.Model
	viewport viewport.Model
	ready    bool
	cahche   [5]string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) changePage() {
	m.page = page(m.list.Index())
	switch m.page {
	case 0:
		m.viewport.SetContent(m.getHome())
	case 1:
		m.viewport.SetContent(m.getAbout())
	case 2:
		m.viewport.SetContent(m.cahche[2])
	case 3:
		m.viewport.SetContent(m.getContact())
	case 4:
		m.viewport.SetContent(m.cahche[4])
	}
}

func (m *model) changeFocus() {
	if m.state == 0 {
		m.state = 1
	} else if m.state == 1 {
		m.state = 0
	}
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if k := msg.String(); k == "ctrl+c" || k == "q" || k == "Q" {
			return m, tea.Quit
		}
		if k := msg.String(); k == "enter" || k == "right" {
			m.changePage()
			m.changeFocus()
		}
		if k := msg.String(); k == "left" {
			if m.state == 1 {
				m.changeFocus()
			}
		}
	case tea.WindowSizeMsg:
		h, v := unselectedListStyle.GetFrameSize()
		m.list.SetSize(((24 * msg.Width) / 100), msg.Height-v)
		m.list.SetShowStatusBar(false)
		m.list.SetShowHelp(false)

		if !m.ready {
			m.viewport = viewport.New((65*msg.Width)/100, msg.Height-v)
			m.viewport.YPosition = h

			m.ready = true
			m.viewport.Style = lipgloss.NewStyle()

			// 			asciiArt := `
			// -----------------------------
			// Welcome to ssh l0calhost.xyz
			// -----------------------------
			//    (\__/) ||
			//    (•ㅅ•) ||
			//  /    づ
			// `
			banner := `
-----------------------------
Welcome to ssh l0calhost.xyz
-----------------------------`
			banner = lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, bannerStyle.Render(banner))
			// banner += "\n"

			bunny := `
(\__/) ||
(•ㅅ•) ||
/    づ`
			bunny = lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, bunnyStyle.Render(bunny))
			bunny += "\n"

			// text := lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, asciiArt)
			// text += "\n"
			c := lipgloss.PlaceHorizontal(m.viewport.Width, lipgloss.Center, helpStyle.Render("Navigation: Arrow Keys + Enter • Quit: Ctrl + C or q"))
			text := banner + bunny + c
			text = lipgloss.PlaceVertical(20, lipgloss.Center, text)
			// 		text += str
			m.viewport.SetContent(text)

		} else {
			m.viewport.Width = (65 * msg.Width) / 100
			m.viewport.Height = msg.Height - v
		}

	}

	var cmd tea.Cmd
	if m.state == 0 {
		m.list, cmd = m.list.Update(msg)

	} else if m.state == 1 {
		m.viewport, cmd = m.viewport.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	// s1 := unselectedListStyle.Render(m.list.View())
	// s2 := borderStyle.Render(m.viewport.View())
	var s string
	s1 := unselectedListStyle.Render(m.list.View())
	s2 := unselectedViewportStyle.Render(m.viewport.View())

	if m.state == 0 {
		s1 = selectedListStyle.Render(m.list.View())
		s2 = unselectedViewportStyle.Render(m.viewport.View())
	} else if m.state == 1 {
		s1 = unselectedListStyle.Render(m.list.View())
		s2 = selectedViewportStyle.Render(m.viewport.View())

	}

	s = lipgloss.JoinHorizontal(0, s1, s2)
	return s
}

func newModel() model {
	items := []list.Item{
		item{title: "Home"},
		item{title: "About"},
		item{title: "Projects"},
		item{title: "Contact"},
		item{title: "Resume"},
	}
	d := list.NewDefaultDelegate()
	c := lipgloss.Color("#00AFFF")
	d.Styles.SelectedTitle = d.Styles.SelectedTitle.Foreground(c).BorderLeftForeground(c)
	d.Styles.SelectedDesc = d.Styles.SelectedDesc.Foreground(c).BorderLeftForeground(c)
	m := model{state: 0, page: 0, list: list.New(items, d, 0, 0)}

	m.list.Title = "l0calhost.xyz"
	m.list.Styles.TitleBar = lipgloss.NewStyle().Padding(1, 1)

	m.cahche[2] = m.getProjects()
	m.cahche[4] = m.getResume()
	return m
}

func teaHandler(s ssh.Session) (tea.Model, []tea.ProgramOption) {
	m := newModel()
	// p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	//render both the list and viewport in the same window

	return m, []tea.ProgramOption{tea.WithAltScreen(), tea.WithMouseCellMotion()}
}

// func main(){
// 	m := newModel()
// 	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
// 	if _, err := p.Run(); err != nil {
// 		os.Exit(1)
// 	}
// }

func main() {
	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.Middleware(teaHandler),
			activeterm.Middleware(), // Bubble Tea apps usually require a PTY.
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err = s.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Error("Could not start server", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer func() { cancel() }()
	if err := s.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		log.Error("Could not stop server", "error", err)
	}
}
