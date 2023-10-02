package tui

import (
  "out"

	"time"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type tickMsg bool
type quitMsg bool

type timerKeyMap struct {
	Reset  key.Binding
	Quit  key.Binding
}

func (k timerKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Reset, k.Quit}
}
func (k timerKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Reset, k.Quit},}
}

var timerKeys = timerKeyMap{
	Reset: key.NewBinding(
		key.WithKeys("space"),
		key.WithHelp("space/r", "Reset Timer"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl-c", "q", "esc"),
		key.WithHelp("ctrl-c/q/esc", "Quit"),
	),
}

type timerModel struct {
  duration int
  count float64
	width    int
	progress progress.Model
	done     bool
	exit     bool
  help help.Model
  keys timerKeyMap
}

func newTimerModel(duration int) timerModel {
	p := progress.New(
		// progress.WithScaledGradient("#6EFA73", "#FF0F80"),
		progress.WithScaledGradient("#2fcdbb", "#ff00e4"),
		progress.WithWidth(80),
		progress.WithoutPercentage(),
	)
	return timerModel{
    duration: duration,
    count: 0.0,
		progress: p,
    help:       help.New(),
    keys:       timerKeys,
	}
}

func (m timerModel) Init() tea.Cmd {
	return nil
}

func (m timerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
    m.progress.Width = msg.Width - 30
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
      m.done = true
			return m, tea.Quit
		case " ", "r":
      m.count = -0.0
			return m, nil
		}
	case quitMsg:
    m.done = false
    m.exit = true
    return m, tea.Quit
	case tickMsg:
		// Update progress bar
    m.count += 0.05
		progressCmd := m.progress.SetPercent(float64(m.count) / float64(m.duration))
		if m.count >= float64(m.duration) {
			m.done = true
      m.exit = true
			return m, tea.Batch(
        tea.Quit,
        // tea.Printf(""),
      )
		}
		return m, progressCmd
	case progress.FrameMsg:
		newModel, cmd := m.progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.progress = newModel
		}
		return m, cmd
	}
	return m, nil
}

func (m timerModel) View() string {
  counter := fmt.Sprintf("%.2f/%d", m.count, m.duration)
// fmt.Sprintf("%.2f", 12.3456)

	prog := m.progress.View()

	cellsRemaining := max(0, m.width-lipgloss.Width(prog + "  " + counter))
	gap := strings.Repeat(" ", cellsRemaining / 2)

	helpView := m.help.View(m.keys)
	cellsRemaining = max(0, m.width-lipgloss.Width(helpView))
	helpGap := strings.Repeat(" ", cellsRemaining / 2)

  if m.exit == false {
    return "\n" + gap + prog + "  " + counter + "\n" + helpGap + helpView + "\n"
  } else {
    return ""
  }
}


func Timer(duration int, quit chan bool) bool {
  // create model
  model := newTimerModel(duration)
  p := tea.NewProgram(model)
  go func(){
    for {
      select {
      default:
        p.Send(tickMsg(true))
      case <- quit:
        p.Send(quitMsg(true))
        return
      }
      time.Sleep(50 * time.Millisecond)
    }
  }()
  // start manager
	tm, err := p.Run()
	if err != nil {
		out.Error("Error running Timer: " + fmt.Sprint(err))
		os.Exit(1)
	}
  mm := tm.(timerModel)
  return mm.done
}
