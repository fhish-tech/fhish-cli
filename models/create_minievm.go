package models

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type CreateMiniEVMModel struct {
	inputs  []textinput.Model
	focused int
	err     error
	done    bool
}

func NewCreateMiniEVMModel() CreateMiniEVMModel {
	m := CreateMiniEVMModel{
		inputs: make([]textinput.Model, 3),
	}

	var t textinput.Model
	t = textinput.New()
	t.Placeholder = "Chain ID (e.g. fhish-1)"
	t.Focus()
	m.inputs[0] = t

	t = textinput.New()
	t.Placeholder = "EVM Chain ID (e.g. 1234)"
	m.inputs[1] = t

	t = textinput.New()
	t.Placeholder = "Moniker (e.g. fhish-node)"
	m.inputs[2] = t

	return m
}

func (m CreateMiniEVMModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m CreateMiniEVMModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if m.focused == len(m.inputs)-1 {
				m.done = true
				return m, tea.Quit
			}
			m.focused++
			m.inputs[m.focused].Focus()
		}
	}

	var cmd tea.Cmd
	m.inputs[m.focused], cmd = m.inputs[m.focused].Update(msg)
	return m, cmd
}

func (m CreateMiniEVMModel) View() string {
	if m.done {
		return "Config collected! Starting build...\n"
	}

	s := "Fhish MiniEVM Configuration Wizard\n\n"
	for i := range m.inputs {
		s += fmt.Sprintf("%s\n", m.inputs[i].View())
	}
	s += "\n(press enter to continue)\n"

	return s
}

func (m CreateMiniEVMModel) GetValues() (string, string, string) {
	return m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value()
}
