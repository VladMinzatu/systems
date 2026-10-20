package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type model struct {
	input  textinput.Model
	output viewport.Model
	lines  []string
	ready  bool
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type something and press Enter..."
	ti.Prompt = "> "
	ti.Focus()

	return model{
		input:  ti,
		output: viewport.New(),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		frameW := panelStyle.GetHorizontalFrameSize()
		frameH := panelStyle.GetVerticalFrameSize()
		// Input panel is one line tall plus its frame; help is one line.
		inputH := 1 + frameH
		helpH := 1

		m.input.SetWidth(msg.Width - frameW - len(m.input.Prompt) - 1)
		m.output.SetWidth(msg.Width - frameW)
		m.output.SetHeight(max(msg.Height-inputH-frameH-helpH, 1))
		m.ready = true

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		case "enter":
			if line := strings.TrimSpace(m.input.Value()); line != "" {
				m.lines = append(m.lines, fmt.Sprintf("%3d │ %s", len(m.lines)+1, line))
				m.output.SetContent(strings.Join(m.lines, "\n"))
				m.output.GotoBottom()
			}
			m.input.Reset()
			return m, nil
		case "pgup", "pgdown":
			// Let the output panel handle scrolling.
			var cmd tea.Cmd
			m.output, cmd = m.output.Update(msg)
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		panelStyle.Render(m.output.View()),
		panelStyle.Render(m.input.View()),
		helpStyle.Render(fmt.Sprintf(" %d lines • enter: append • pgup/pgdown: scroll • esc/ctrl+c: quit", len(m.lines))),
	)

	v := tea.NewView(content)
	v.AltScreen = true
	v.WindowTitle = "tui-test"
	return v
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
