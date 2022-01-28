package list

import (
	"get-service-version/entity"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type delegateKeyMap struct {
	choose key.Binding
}

var keyMap *delegateKeyMap

func init() {
	keyMap = &delegateKeyMap{
		choose: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "choose"),
		),
	}
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	list.Model
	OnUpdate OnUpdate
}

type OnEnter func(string) tea.Cmd
type OnUpdate func(msg tea.Msg, m *list.Model) tea.Cmd

func NewListModel(title string, items []list.Item, onEnter OnEnter) Model {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.UpdateFunc = newUpdateFunc(onEnter)
	l := list.New(items, d, 0, 0)
	l.Title = title
	return Model{
		Model: l,
	}
}

func NewListOptions() []list.Item {
	items := []list.Item{
		item{title: entity.BootstrapV3},
		item{title: entity.BootstrapV4},
	}
	return items
}

func newUpdateFunc(onEnter OnEnter) func(msg tea.Msg, m *list.Model) tea.Cmd {
	return func(msg tea.Msg, m *list.Model) tea.Cmd {
		var title string

		if i, ok := m.SelectedItem().(item); ok {
			title = i.Title()
		} else {
			return nil
		}

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(msg, keyMap.choose):
				return onEnter(title)
			}
		}
		return nil
	}

}
