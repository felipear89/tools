package model

import (
	"fmt"
	"get-service-version/entity"
	"get-service-version/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

const (
	Listing    = iota
	Loading    = iota
	Displaying = iota
)

type Model struct {
	List    list.Model
	Spinner spinner.Model
	State   *entity.State
}

func (m Model) Init() tea.Cmd {
	m.State.Selected = Listing
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.List.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	switch m.State.Selected {
	case Listing:
		m.List.Model, cmd = m.List.Update(msg)
		return m, cmd
	case Loading:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	case Displaying:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.State.Selected = Listing
			}
		}
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.State.Selected {
	case Listing:
		return docStyle.Render(m.List.View())
	case Loading:
		str := fmt.Sprintf("\n\n   %s Loading ...\n\n", m.Spinner.View())
		return str
	case Displaying:
		return m.State.Screen
	}
	return ""
}

func (m Model) Loading() func() tea.Msg {
	m.State.Selected = Loading
	m.Spinner.Start()
	return m.Spinner.Tick
}

func (m Model) Ready(screen string) {
	m.State.Selected = Displaying
	m.Spinner.Finish()
	m.State.Screen = screen
}
