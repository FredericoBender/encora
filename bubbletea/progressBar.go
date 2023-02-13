package bubbletea

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyleLoad = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type tickMsg float64

type LoadModel struct {
	Progress          progress.Model
	CurrentPercentage float64
}

func (m *LoadModel) Init() tea.Cmd {
	return m.TickCmd()
}

func (m *LoadModel) CreateProgressBar() {
	if err := tea.NewProgram(m).Start(); err != nil {
		log.Fatalf("progress bar: %s", err)
		os.Exit(1)
	}
}

func (m *LoadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.Progress.Width = msg.Width - padding*2 - 4
		if m.Progress.Width > maxWidth {
			m.Progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.Progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		cmd := m.Progress.SetPercent(float64(msg))
		return m, tea.Batch(m.TickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.Progress.Update(msg)
		m.Progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func (e *LoadModel) View() string {
	pad := strings.Repeat(" ", padding)
	return "\n" +
		pad + e.Progress.View() + "\n\n" +
		pad + helpStyleLoad("Press any key to quit")
}

func (m *LoadModel) TickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(m.CurrentPercentage)
	})
}
