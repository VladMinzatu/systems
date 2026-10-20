package main

import (
	"fmt"
	"os"
	"strings"

	"charm.land/bubbles/v2/textinput"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/NimbleMarkets/ntcharts/v2/canvas/runes"
	"github.com/NimbleMarkets/ntcharts/v2/linechart/streamlinechart"
)

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1)
	helpStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	lineStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("42"))
)

// ntcharts panics if the chart is smaller than its axes, so keep a floor.
const (
	minChartW = 20
	minChartH = 5
)

type model struct {
	input  textinput.Model
	output viewport.Model
	lines  []string

	cpuChart   streamlinechart.Model
	cpuPrev    cpuStat
	cpuHasPrev bool
	cpuNow     float64
	cpuErr     error
	cpuWidth   int

	ready bool
}

func newModel() model {
	ti := textinput.New()
	ti.Placeholder = "Type something and press Enter..."
	ti.Prompt = "> "
	ti.Focus()

	chart := streamlinechart.New(minChartW, minChartH, streamlinechart.WithYRange(0, 100))
	chart.SetStyles(runes.ArcLineStyle, lineStyle)

	return model{
		input:    ti,
		output:   viewport.New(),
		cpuChart: chart,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, sampleCPU())
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
		bodyH := msg.Height - helpH

		m.cpuWidth = max(msg.Width*2/5, 30)
		leftW := msg.Width - m.cpuWidth

		m.input.SetWidth(leftW - frameW - len(m.input.Prompt) - 1)
		m.output.SetWidth(leftW - frameW)
		m.output.SetHeight(max(bodyH-inputH-frameH, 1))

		// The CPU panel spans the full body height; one line is the title.
		m.cpuChart.Resize(max(m.cpuWidth-frameW, minChartW), max(bodyH-frameH-1, minChartH))
		m.cpuChart.Draw()
		m.ready = true

	case cpuSampleMsg:
		m.cpuErr = msg.err
		if msg.err == nil {
			if m.cpuHasPrev {
				m.cpuNow = cpuPercent(m.cpuPrev, msg.stat)
				m.cpuChart.Push(m.cpuNow)
				m.cpuChart.Draw()
			}
			m.cpuPrev, m.cpuHasPrev = msg.stat, true
		}
		return m, sampleCPU()

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

func (m model) cpuView() string {
	var title, body string
	switch {
	case m.cpuErr != nil:
		title = titleStyle.Render("CPU")
		body = errStyle.Render(m.cpuErr.Error())
	case !m.cpuHasPrev:
		title = titleStyle.Render("CPU")
		body = helpStyle.Render("sampling...")
	default:
		title = titleStyle.Render(fmt.Sprintf("CPU %5.1f%%", m.cpuNow))
		body = m.cpuChart.View()
	}
	return lipgloss.JoinVertical(lipgloss.Left, title, body)
}

func (m model) View() tea.View {
	if !m.ready {
		return tea.NewView("Initializing...")
	}

	left := lipgloss.JoinVertical(lipgloss.Left,
		panelStyle.Render(m.output.View()),
		panelStyle.Render(m.input.View()),
	)
	right := panelStyle.
		Width(m.cpuWidth).
		Height(lipgloss.Height(left)).
		Render(m.cpuView())

	content := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Top, left, right),
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
