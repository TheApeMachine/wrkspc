package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Layer struct {
	model tea.Model
}

func NewLayer(model tea.Model) *Layer {
	return &Layer{model}
}

func (layer *Layer) Init() tea.Cmd {
	return layer.model.Init()
}

func (layer *Layer) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return layer.model.Update(msg)
}

func (layer *Layer) View() string {
	return layer.model.View()
}
