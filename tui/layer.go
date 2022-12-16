package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/davecgh/go-spew/spew"
	"github.com/theapemachine/wrkspc/tui/layers"
)

type LayerType []byte

func (lt LayerType) Bytes() []byte {
	return []byte(lt)
}

var (
	LOGO    LayerType = []byte("logo")
	SPINNER LayerType = []byte("spinner")
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
	spew.Dump(msg)
	return layer.model.Update(msg)
}

func (layer *Layer) View() string {
	return layer.model.View()
}

var core = map[string]tea.Model{
	"logo":    layers.NewLogo(),
	"spinner": layers.NewSpinner(),
}
