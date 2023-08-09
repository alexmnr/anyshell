package tui

import (
	"out"

	"fmt"
	"os"
	"regexp"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
)

type selectModel struct {
  options []string
  shown []string
  legend string
  filter string
  filtering bool
  index int
  error bool
  help help.Model
  keys keyMap
}

var (
  selector            = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0F80"))
  filter              = lipgloss.NewStyle().Foreground(lipgloss.Color("#04D1F1"))
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Search    key.Binding
	Escape   key.Binding
	Select key.Binding
	Quit  key.Binding
}
func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Search, k.Escape, k.Select, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{k.Up, k.Down, k.Search, k.Escape, k.Select, k.Quit},}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Search: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "search"),
	),
	Escape: key.NewBinding(
		key.WithKeys("Esc"),
		key.WithHelp("Esc", "stop searching"),
	),
	Select: key.NewBinding(
		key.WithKeys("Enter"),
		key.WithHelp("Enter", "select"),
	),
	Quit: key.NewBinding(
		key.WithKeys("ctrl-c", "q"),
		key.WithHelp("ctrl-c/q", "quit"),
	),
}

func (m selectModel) Init() tea.Cmd {
  m.shown = m.options
  m.filtering = false
  return nil
}

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
    // no filtering
    if m.filtering == false {
      switch msg.String() {
      case "/":
        m.filtering = true
      case "j":
        m.index += 1
        if m.index >= len(m.shown) {
          m.index = 0
        }
      case "k":
        m.index -= 1
        if m.index < 0 {
          m.index = len(m.shown) - 1
        }
      }
    // with filtering
    } else {
      if len(msg.String()) == 1 {
        m.filter += msg.String()
        m.filter = regexp.MustCompile(`[^a-z]+`).ReplaceAllString(m.filter, "")
      }
      switch msg.String() {
      case "esc":
        m.filtering = false
      case "/":
        m.filtering = false
      case "backspace":
        m.filter = m.filter[:len(m.filter)-1]
      }
    }
    // every time
    switch msg.String() {
    case "q", "ctrl+c":
      m.error = true
      return m, tea.Quit
    case "enter":
      return m, tea.Quit
    case "down":
      m.index += 1
      if m.index >= len(m.shown) {
        m.index = 0
      }
    case "up":
      m.index -= 1
      if m.index < 0 {
        m.index = len(m.shown) - 1
      }
    }
	}

  if m.filter == "" {
    m.shown = m.options
  } else {
    m.shown = nil
    for _, k := range m.options {
      if strings.Contains(k, m.filter) == true {
        m.shown = append(m.shown, k)
      }
    }
  }
	return m, cmd
}

func (m selectModel) View() string {
  var output string
  if m.legend != "" {
    output += selector.Render("   ") + m.legend + " "
  }
  if m.filtering == true {
    output += filter.Render("[" + m.filter + "]")
  }
  output += "\n"
  for n, k := range m.shown {
    if n == m.index {
      output += selector.Render(" > ") + lipgloss.NewStyle().Bold(true).Render(k) + "\n"
    } else {
      output += "   " + k + "\n"
    }
  }

	helpView := m.help.View(m.keys)

  return output + "\n" + helpView
}

func Select(options []string, legend string) int {
  var index int
	m := selectModel{
    options: options,
    legend: legend,
    help:       help.New(),
    keys:       keys,
  }
  tm, err := tea.NewProgram(m).Run();
	if err != nil {
		out.Error("Error running Selector: " + fmt.Sprint(err))
		os.Exit(1)
	}
  mm := tm.(selectModel)
  if mm.error == true {
		out.Error("You need to select something!")
    os.Exit(0)
  }
  index = mm.index
  return index
}
