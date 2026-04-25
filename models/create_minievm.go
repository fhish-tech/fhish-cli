package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Step int

const (
	StepChainID Step = iota
	StepEVMChainID
	StepMoniker
	StepGasDenom
	StepL1RPC
	StepConfirm
)

type CreateMiniEVMModel struct {
	step    Step
	history []Step
	inputs  map[Step]*textinput.Model
	err     error
	Done    bool
}

func NewCreateMiniEVMModel() CreateMiniEVMModel {
	inputs := make(map[Step]*textinput.Model)
	
	t1 := textinput.New()
	t1.Placeholder = "Chain ID (fhish-1)"
	t1.SetValue("fhish-1")
	t1.Focus()
	inputs[StepChainID] = &t1

	t2 := textinput.New()
	t2.Placeholder = "EVM Chain ID (1234)"
	t2.SetValue("1234")
	inputs[StepEVMChainID] = &t2

	t3 := textinput.New()
	t3.Placeholder = "Moniker (fhish-node)"
	t3.SetValue("fhish-node")
	inputs[StepMoniker] = &t3

	t4 := textinput.New()
	t4.Placeholder = "Gas Denom (uinit)"
	t4.SetValue("uinit")
	inputs[StepGasDenom] = &t4

	t5 := textinput.New()
	t5.Placeholder = "L1 RPC (https://rpc.testnet.initia.xyz)"
	t5.SetValue("https://rpc.testnet.initia.xyz")
	inputs[StepL1RPC] = &t5

	return CreateMiniEVMModel{
		step:    StepChainID,
		inputs:  inputs,
		history: []Step{},
	}
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
		case "ctrl+z":
			if len(m.history) > 0 {
				m.step = m.history[len(m.history)-1]
				m.history = m.history[:len(m.history)-1]
				m.inputs[m.step].Focus()
			}
			return m, nil
		case "enter":
			if m.step == StepConfirm {
				m.Done = true
				return m, tea.Quit
			}
			m.history = append(m.history, m.step)
			m.step++
			if m.step <= StepL1RPC {
				m.inputs[m.step].Focus()
			}
			return m, nil
		}
	}

	if m.step <= StepL1RPC {
		var cmd tea.Cmd
		input := *m.inputs[m.step]
		newModel, cmd := input.Update(msg)
		*m.inputs[m.step] = newModel
		return m, cmd
	}

	return m, nil
}

func (m CreateMiniEVMModel) View() string {
	if m.Done {
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render("✔ Configuration complete! Building...")
	}

	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("5")).Render("Fhish MiniEVM Setup Wizard")
	
	var body string
	switch m.step {
	case StepChainID:
		body = "Enter Chain ID:\n" + m.inputs[StepChainID].View()
	case StepEVMChainID:
		body = "Enter EVM Chain ID:\n" + m.inputs[StepEVMChainID].View()
	case StepMoniker:
		body = "Enter Node Moniker:\n" + m.inputs[StepMoniker].View()
	case StepGasDenom:
		body = "Enter Gas Denom:\n" + m.inputs[StepGasDenom].View()
	case StepL1RPC:
		body = "Enter L1 RPC URL:\n" + m.inputs[StepL1RPC].View()
	case StepConfirm:
		body = "Confirm configuration and start deployment? (Press Enter)"
	}

	help := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render("\n\n(enter: next • ctrl+z: back • ctrl+c: quit)")
	
	return header + "\n\n" + body + help
}

func (m CreateMiniEVMModel) Config() map[string]string {
	res := make(map[string]string)
	res["chain_id"] = m.inputs[StepChainID].Value()
	res["evm_chain_id"] = m.inputs[StepEVMChainID].Value()
	res["moniker"] = m.inputs[StepMoniker].Value()
	res["gas_denom"] = m.inputs[StepGasDenom].Value()
	res["l1_rpc"] = m.inputs[StepL1RPC].Value()
	return res
}
