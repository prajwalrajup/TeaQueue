package models

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type topicItem struct {
	title, desc string
}

func (i topicItem) Title() string       { return i.title }
func (i topicItem) Description() string { return i.desc }
func (i topicItem) FilterValue() string { return i.title }

type TopicsModel struct {
	list list.Model
}

func (m TopicsModel) Init() tea.Cmd {
	return nil
}

func (m TopicsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m TopicsModel) View() string {
	return docStyle.Render(m.list.View())
}

func InitTopicsModel() TopicsModel {
	topicsItems := []list.Item{
		topicItem{title: "localSetup", desc: "localhost:9092"},
		topicItem{title: "dockerSetup", desc: "localhost:9092"},
		topicItem{title: "K8sDev", desc: "localhost:9092"},
	}
	topicsList := list.New(topicsItems, list.NewDefaultDelegate(), 0, 0)
	topicsList.Title = "Topics"
	topicsList.SetSize(50, 50)
	return TopicsModel{list: topicsList}
}
