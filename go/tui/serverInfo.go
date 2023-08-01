package tui

import (
  "out"

	"fmt"
	"strings"
  "regexp"
  "os"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ServerInfo struct {
  Name string
  DbPort string
  SshPort string
  WebPort string
  UserPassword string
  RootPassword string
  WebInterface bool
}

var (
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorStyle         = focusedStyle.Copy()
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type serverInfoModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initServerInfoModel() serverInfoModel {
	m := serverInfoModel{
		inputs: make([]textinput.Model, 5),
	}

	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Database Name [anyshell]"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "db-port [42998]"
			t.CharLimit = 5
		case 2:
			t.Placeholder = "ssh-port [42999]"
			t.CharLimit = 5
		case 3:
			t.Placeholder = "user-password (needed for login)"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 4:
			t.Placeholder = "db-root-password (only for maintance and web interface)"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m serverInfoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m serverInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				return m, tea.Quit
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			if m.focusIndex > len(m.inputs) {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = len(m.inputs)
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i <= len(m.inputs)-1; i++ {
				if i == m.focusIndex {
					// Set focused state
					cmds[i] = m.inputs[i].Focus()
					m.inputs[i].PromptStyle = focusedStyle
					m.inputs[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.inputs[i].Blur()
				m.inputs[i].PromptStyle = noStyle
				m.inputs[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)

	return m, cmd
}

func (m *serverInfoModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m serverInfoModel) View() string {
	var b strings.Builder

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if m.focusIndex == len(m.inputs) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n", *button)
	return b.String()
}

func GetServerInfo() ServerInfo {
  // create struct
  info := ServerInfo{}
  // get input
  m := initServerInfoModel()
  tm, _ := tea.NewProgram(&m).Run()
  mm := tm.(serverInfoModel)
  // var name, dbPort, sshPort, userPassword, rootPassword string
  info.Name = mm.inputs[0].Value()
  info.DbPort = mm.inputs[1].Value()
  info.SshPort = mm.inputs[2].Value()
  info.UserPassword = mm.inputs[3].Value()
  info.RootPassword = mm.inputs[4].Value()
  // clean and interpretate input
  if len(info.Name) != 0 {
    info.Name = cleanString(info.Name)
  } else {
    info.Name = "anyshell"
  }
  info.DbPort = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(info.DbPort, "")
  if len(info.DbPort) == 0 {
    info.DbPort = "42998"
  }
  info.SshPort = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(info.SshPort, "")
  if len(info.SshPort) == 0 {
    info.SshPort = "42999"
  }
  if len(info.UserPassword) == 0 {
    out.Error("You need to specify a password!")
    os.Exit(0)
  }
  if len(info.RootPassword) == 0 {
    out.Error("You need to specify a password!")
    os.Exit(0)
  }
  // ask if web interface is needed
  fmt.Println()
  message := "Add webInterface for managing database?"
  var options []string
  options = append(options, out.Style("Yes", 1, false))
  options = append(options, out.Style("No", 0, false))
  ret := Survey(message, options)
  if ret == options[0] {
    info.WebInterface = true
  } else {
    info.WebInterface = false
  }

  // ask for WebPort
  if info.WebInterface == true {
    m := initWebPortModel()
    tm, _ := tea.NewProgram(&m).Run()
    mm := tm.(webPortModel)
    info.WebPort = mm.inputs[0].Value()
    info.WebPort = mm.inputs[0].Value()
    info.WebPort = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(info.WebPort, "")
    if len(info.WebPort) == 0 {
      info.WebPort = "42997"
    }
  }

  return info
}

func cleanString(input string) string {
  clean := strings.Replace(input, " ", "", -1)
  clean = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(clean, "")
  return clean
} 
