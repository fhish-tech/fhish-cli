package models

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Step int

const (
	StepChainID Step = iota
	StepMoniker
	StepGasDenom
	StepDeployerKey
	StepRelayerSecret
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
	t1.Placeholder = "Chain ID (e.g. fhish-1)"
	t1.SetValue("fhish-1")
	t1.Focus()
	inputs[StepChainID] = &t1

	t2 := textinput.New()
	t2.Placeholder = "Node Moniker (e.g. fhish-node)"
	t2.SetValue("fhish-node")
	inputs[StepMoniker] = &t2

	t3 := textinput.New()
	t3.Placeholder = "Gas Denom (uinit)"
	t3.SetValue("uinit")
	inputs[StepGasDenom] = &t3

	t4 := textinput.New()
	t4.Placeholder = "Deployer Private Key (0x...)"
	t4.EchoMode = textinput.EchoPassword
	t4.EchoCharacter = '•'
	t4.SetValue("0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80")
	inputs[StepDeployerKey] = &t4

	t5 := textinput.New()
	t5.Placeholder = "FHE Relayer Secret"
	t5.EchoMode = textinput.EchoPassword
	t5.EchoCharacter = '•'
	t5.SetValue("fhish-test-secret")
	inputs[StepRelayerSecret] = &t5

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
			if m.step <= StepRelayerSecret {
				m.inputs[m.step].Focus()
			}
			return m, nil
		}
	}

	if m.step <= StepRelayerSecret {
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
		return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Bold(true).Render("\n✔ Configuration complete! Initializing Fhish Stack...")
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Render("Fhish Rollup Setup Wizard")
	
	var body string
	switch m.step {
	case StepChainID:
		body = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Chain ID") + "\n" + m.inputs[StepChainID].View()
	case StepMoniker:
		body = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Node Moniker") + "\n" + m.inputs[StepMoniker].View()
	case StepGasDenom:
		body = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Gas Denom") + "\n" + m.inputs[StepGasDenom].View()
	case StepDeployerKey:
		body = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("Deployer Private Key (EVM)") + "\n" + m.inputs[StepDeployerKey].View()
	case StepRelayerSecret:
		body = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("FHE Relayer Secret") + "\n" + m.inputs[StepRelayerSecret].View()
	case StepConfirm:
		body = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214")).Render("Ready to launch? (Press Enter to start deployment)")
	}

	help := lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render("\n\n(enter: next • ctrl+z: back • ctrl+c: quit)")
	
	return header + "\n\n" + body + help
}

func (m CreateMiniEVMModel) Config() map[string]string {
	res := make(map[string]string)
	res["chain_id"] = m.inputs[StepChainID].Value()
	res["moniker"] = m.inputs[StepMoniker].Value()
	res["gas_denom"] = m.inputs[StepGasDenom].Value()
	res["deployer_key"] = m.inputs[StepDeployerKey].Value()
	res["relayer_secret"] = m.inputs[StepRelayerSecret].Value()
	return res
}
