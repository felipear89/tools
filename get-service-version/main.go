package main

import (
	"fmt"
	"get-service-version/display"
	"get-service-version/entity"
	"get-service-version/list"
	"get-service-version/model"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
)

func main() {

	m := &model.Model{
		State:   &entity.State{},
		Spinner: newSpinner(),
	}
	onUpdate := func(title string) tea.Cmd {
		go func() {
			display.Display(m, title)
		}()
		return m.Loading()
	}
	m.List = list.NewListModel("Select the service name", list.NewListOptions(), onUpdate)

	p := tea.NewProgram(m, tea.WithAltScreen())

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func newSpinner() spinner.Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return s
}
