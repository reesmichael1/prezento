package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type Pages struct {
	slides    []string
	paginator paginator.Model
	ready     bool
	viewport  viewport.Model
}

func NewPages(slides []string) Pages {
	pages := paginator.New()
	pages.Type = paginator.Dots
	pages.PerPage = 1
	pages.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	pages.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")
	pages.SetTotalPages(len(slides))

	return Pages{
		slides:    slides,
		paginator: pages,
	}
}

func (p Pages) Init() tea.Cmd {
	return nil
}

func (p Pages) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return p, tea.Quit
		}

	case tea.WindowSizeMsg:
		// headerHeight := lipgloss.Height(m.headerView())
		footerHeight := lipgloss.Height(p.paginator.View())

		// verticalMarginHeight := headerHeight + footerHeight
		if !p.ready {
			// Since this program is using the full size of the viewport we
			// need to wait until we've received the window dimensions before
			// we can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			p.viewport = viewport.New(msg.Width, msg.Height-footerHeight)
			// m.viewport.YPosition = 0
			p.ready = true
		} else {
			p.viewport.Width = msg.Width
			p.viewport.Height = msg.Height - footerHeight
		}
	}

	p.paginator, cmd = p.paginator.Update(msg)

	rendered, err := glamour.Render(p.slides[p.paginator.Page], "dark")
	if err != nil {
		panic(err)
	}
	p.viewport.SetContent(rendered)

	cmds = append(cmds, cmd)

	p.viewport, cmd = p.viewport.Update(msg)

	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p Pages) View() string {
	if !p.ready {
		return "\n  Initializing..."
	}

	var b strings.Builder

	b.WriteString(" " + lipgloss.NewStyle().Width(p.viewport.Width).Align(lipgloss.Center).Render(p.viewport.View()))
	b.WriteString("\n\n" + lipgloss.NewStyle().Width(p.viewport.Width).Align(lipgloss.Center).Render(p.paginator.View()))

	return b.String()
}
