package tui

import (
  "tools"
  "command"

	"os"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Action struct{
  Name string
  Cmd func() error
  Interactive bool
}
type installedPkgMsg string
var (
	currentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6CD0D4"))
	doneStyle           = lipgloss.NewStyle().Margin(1, 2)
	checkMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("#4EF465")).SetString("✓")
)
type single_action_model struct {
	action   Action
	width    int
	height   int
	spinner  spinner.Model
	done     bool
  debug    bool
}

func new_single_model(action Action, debug bool) single_action_model {
	s := spinner.New()
	s.Spinner = spinner.MiniDot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))
	return single_action_model{
	  action: action,
		spinner:  s,
		debug: debug,
	}
}

func (m single_action_model) Init() tea.Cmd {
	return tea.Batch(
    m.spinner.Tick,
  )
}

func (m single_action_model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			return m, tea.Quit
		}
	case installedPkgMsg:
    m.done = true
    return m, tea.Sequence(
      tea.Printf("%s %s", checkMark, m.action.Name), // print success message above our program
      tea.Quit,
    )
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m single_action_model) View() string {
	spin := m.spinner.View() + " "
	cellsAvail := max(0, m.width-lipgloss.Width(spin))

	info := lipgloss.NewStyle().MaxWidth(cellsAvail).Render(m.action.Name)

  if m.debug == true {
    return ""
  }
  if m.action.Interactive == true {
    return ""
  }

	return spin + info 
}

func RunAction(name string, cmd func() error, debug bool) {
  action := Action{
    Name: name,
    Cmd: cmd,
  }
  // get sudo rights
  if tools.GetUser() != "root"{
    command.Cmd("sudo true", true)
  }
  // create model
  model := new_single_model(action, debug)
  p := tea.NewProgram(model)
  // run actions
  go func(){
    err := action.Cmd()
    if err != nil {
        p.Kill()
        p.Quit()
    }
    p.Send(installedPkgMsg(action.Name))
  }()
  // start manager
	if _, err := p.Run(); err != nil {
		os.Exit(0)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
