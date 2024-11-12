package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type ProfileModel struct {
	list           list.Model
	CurrentSeleted string
}

func (m ProfileModel) Init() tea.Cmd {
	return nil
}

func (m ProfileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	selectedItem := m.list.SelectedItem()
	if selectedItem != nil {
		m.CurrentSeleted = selectedItem.(item).Title()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ProfileModel) View() string {
	return docStyle.Render(m.list.View())
}

func InitProfileModel() ProfileModel {
	items := []list.Item{
		item{title: "localSetup", desc: "localhost:9092"},
		item{title: "dockerSetup", desc: "localhost:9092"},
		item{title: "K8sDev", desc: "localhost:9092"},
	}

	ist := list.New(items, list.NewDefaultDelegate(), 0, 0)
	ist.Title = "Profiles"
	ist.SetSize(50, 50)
	return ProfileModel{list: ist}
}
