package models

import (
	"TeaQueue/utils"
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

	profileSelectedItem := m.list.SelectedItem()
	if profileSelectedItem != nil {
		m.CurrentSeleted = profileSelectedItem.(item).Title()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ProfileModel) View() string {
	return docStyle.Render(m.list.View())
}

func InitProfileModel(config utils.Config) ProfileModel {
	items := []list.Item{}

	for key, value := range config.Profile {
		items = append(items, item{title: key, desc: value.Desc})
	}

	ist := list.New(items, list.NewDefaultDelegate(), 0, 0)
	ist.Title = "Profiles"
	ist.SetSize(50, 50)
	return ProfileModel{list: ist}
}
