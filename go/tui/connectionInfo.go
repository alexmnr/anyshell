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
)

type ConnectionInfo struct {
  Name string
  Host string
  SshPort string
  DbPort string
  Password string
}

type ConnectionInfoModel struct {
	focusIndex int
	inputs     []textinput.Model
	cursorMode cursor.Mode
}

func initConnectionInfoModel() ConnectionInfoModel {
	m := ConnectionInfoModel{
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
			t.Placeholder = "Host (domain or IP)"
		case 2:
			t.Placeholder = "ssh-port [42999]"
			t.CharLimit = 5
		case 3:
			t.Placeholder = "db-port [42998]"
			t.CharLimit = 5
		case 4:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m ConnectionInfoModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m ConnectionInfoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m *ConnectionInfoModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m ConnectionInfoModel) View() string {
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

func GetConnectionInfo() ConnectionInfo {
  // create struct
  info := ConnectionInfo{}
  // get input
  m := initConnectionInfoModel()
  tm, _ := tea.NewProgram(&m).Run()
  mm := tm.(ConnectionInfoModel)
  // var name, dbPort, sshPort, userPassword, rootPassword string
  info.Name = mm.inputs[0].Value()
  info.Host = mm.inputs[1].Value()
  info.SshPort = mm.inputs[2].Value()
  info.DbPort = mm.inputs[3].Value()
  info.Password = mm.inputs[4].Value()
  // clean and interpretate input
  if len(info.Name) != 0 {
    info.Name = cleanString(info.Name)
  } else {
    info.Name = "anyshell"
  }
  info.DbPort = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(info.DbPort, "")
  if len(info.DbPort) == 0 {
    info.DbPort = fmt.Sprint(42998)
  }
  info.SshPort = regexp.MustCompile(`[^0-9]+`).ReplaceAllString(info.SshPort, "")
  if len(info.SshPort) == 0 {
    info.SshPort = fmt.Sprint(42999)
  }
  if len(info.Host) == 0 {
    out.Error("You need to specify a host!")
    os.Exit(1)
  }
  if len(info.Password) == 0 {
    out.Error("You need to specify a password!")
    os.Exit(1)
  }

  return info
}
