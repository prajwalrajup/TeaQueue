package models

import (
	"TeaQueue/utils"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

var serverConfigFile utils.Config

type serverItem struct {
	title, desc string
}

func (i serverItem) Title() string       { return i.title }
func (i serverItem) Description() string { return i.desc }
func (i serverItem) FilterValue() string { return i.title }

type ServerModel struct {
	list            list.Model
	SelectedProfile string
	SelectedServer  string
}

func (m ServerModel) Init() tea.Cmd {
	return nil
}

func (m ServerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			return m, nil
		case "ctrl+c":
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	for profileName, profile := range serverConfigFile.Profile {
		switch m.SelectedProfile {
		case profileName:
			currentServerItems := []list.Item{}
			for serverName, server := range profile.Servers {
				currentServerItems = append(currentServerItems, serverItem{title: serverName, desc: server.Desc})
			}
			m.list.SetItems(currentServerItems)
		}
	}

	serverSelectedItem := m.list.SelectedItem()
	if serverSelectedItem != nil {
		m.SelectedServer = serverSelectedItem.(serverItem).Title()
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m ServerModel) View() string {
	return docStyle.Render(m.list.View())
}

func InitServerModel(config utils.Config) ServerModel {
	serversItems := []list.Item{}

	serverConfigFile = config
	serversList := list.New(serversItems, list.NewDefaultDelegate(), 0, 0)
	serversList.Title = "Topics"
	serversList.SetSize(50, 50)
	return ServerModel{list: serversList}
}
