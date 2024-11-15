package main

import (
	"fmt"
	"os"

	"TeaQueue/models"
	"TeaQueue/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type modelID int

const (
	profileId modelID = iota
	topicId
)

type MainModel struct {
	currentModel modelID
	profileModel models.ProfileModel
	topicsModel  models.ServerModel
}

func (m *MainModel) Init() tea.Cmd {
	return nil
}

func (m *MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var updatedData tea.Model
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if m.currentModel == topicId {
				m.currentModel = profileId
			}
			return m, nil
		case "enter":
			switch m.currentModel {
			case profileId:
				updatedData, cmd = m.profileModel.Update(msg)
				m.profileModel = updatedData.(models.ProfileModel)
				m.currentModel = topicId
				m.topicsModel.SelectedProfile = m.profileModel.CurrentSeleted
			case topicId:
				updatedData, cmd = m.topicsModel.Update(msg)
				m.topicsModel = updatedData.(models.ServerModel)

			}

		case "ctrl+c":
			return m, tea.Quit
		}
	}

	switch m.currentModel {
	case profileId:
		updatedData, cmd = m.profileModel.Update(msg)
		m.profileModel = updatedData.(models.ProfileModel)
	case topicId:
		updatedData, cmd = m.topicsModel.Update(msg)
		m.topicsModel = updatedData.(models.ServerModel)
	}

	return m, cmd
}

func (m *MainModel) View() string {
	res := ""
	switch m.currentModel {
	case profileId:
		res = m.profileModel.View()
	case topicId:
		res = m.topicsModel.View()
	default:
		return ""
	}

	return res
}

func newModel() tea.Model {
	config, err := utils.ReadConfig()
	if err != nil {
		panic(err)
	}

	return &MainModel{
		currentModel: profileId,
		profileModel: models.InitProfileModel(config),
		topicsModel:  models.InitServerModel(config),
	}
}

func main() {
	if _, err := tea.NewProgram(newModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
