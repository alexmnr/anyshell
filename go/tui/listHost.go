package tui

import (
	"db"
	"out"
	"strconv"
	"strings"
	"types"

	"fmt"
	"os"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type listHostModel struct {
  hosts [][]types.HostInfo
  servers []types.ConnectionInfo
  shownHosts [][]types.HostInfo
  shownServers []types.ConnectionInfo

  // legend string
  filter string
  filtering bool
  idFilter string

  verbose bool
  onlyOnline bool

  error bool
  width int
}


func (m listHostModel) Init() tea.Cmd {
  m.shownHosts = m.hosts
  m.shownServers = m.servers
  m.filtering = false
  m.verbose = false
  m.onlyOnline = false
  return nil
}

func (m listHostModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
  length := 0
  for _, k := range m.hosts {
    length += len(k)
  }
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width= msg.Width
	case tea.KeyMsg:
    // id filtering
    _, err := strconv.Atoi(msg.String())
    if err == nil {
      m.idFilter += msg.String() 
      m.filtering = false
      m.filter = ""
    }
    // no filtering
    if m.filtering == false {
      switch msg.String() {
      case "/":
        m.filtering = true
        m.idFilter = ""
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
      }
    }
    // every time
    switch msg.String() {
    case "backspace":
      if len(m.filter) > 0 {
        m.filter = m.filter[:len(m.filter)-1]
      }
      if len(m.idFilter) > 0 {
        m.idFilter = m.idFilter[:len(m.idFilter)-1]
      }
    case "v":
      if m.verbose == true {m.verbose = false} else {m.verbose=true}
    case "o":
      if m.onlyOnline == true {m.onlyOnline = false} else {m.onlyOnline=true}
    case "q", "ctrl+c":
      m.error = true
      return m, tea.Quit
    case "enter":
      return m, tea.Quit
    }
	}

  if m.filter == "" {
    m.shownHosts = m.hosts
    m.shownServers = m.servers
  } else {
    m.shownHosts = nil
    m.shownServers = nil
    for n, hosts := range m.hosts {
      var buffer []types.HostInfo
      for _, host := range hosts {
        string := host.Name + " " + host.User
        if strings.Contains(string, m.filter) == true {
          buffer = append(buffer, host)
        }
      }
      if len(buffer) != 0 {
        m.shownHosts = append(m.shownHosts, buffer)
        m.shownServers = append(m.shownServers, m.servers[n])
      }
    }
  }

  if m.filtering == false {
    if m.idFilter == "" {
      m.shownHosts = m.hosts
      m.shownServers = m.servers
    } else {
      m.shownHosts = nil
      m.shownServers = nil
      for n, hosts := range m.hosts {
        var buffer []types.HostInfo
        for _, host := range hosts {
          if strings.Contains(fmt.Sprint(host.ID), m.idFilter) == true {
            buffer = append(buffer, host)
          }
        }
        if len(buffer) != 0 {
          m.shownHosts = append(m.shownHosts, buffer)
          m.shownServers = append(m.shownServers, m.servers[n])
        }
      }
    }
  }
  if m.filtering == false {
    if m.idFilter == "" {
      m.shownHosts = m.hosts
      m.shownServers = m.servers
    } else {
      m.shownHosts = nil
      m.shownServers = nil
      for n, hosts := range m.hosts {
        var buffer []types.HostInfo
        for _, host := range hosts {
          if strings.Contains(fmt.Sprint(host.ID), m.idFilter) == true {
            buffer = append(buffer, host)
          }
        }
        if len(buffer) != 0 {
          m.shownHosts = append(m.shownHosts, buffer)
          m.shownServers = append(m.shownServers, m.servers[n])
        }
      }
    }
  }

  if m.onlyOnline == true {
    hostbuffer := m.shownHosts
    serverbuffer := m.shownServers
    m.shownHosts = nil
    m.shownServers = nil
    for n, hosts := range hostbuffer {
      var buffer []types.HostInfo
      for _, host := range hosts {
        if host.Online == true {
          buffer = append(buffer, host)
        }
      }
      if len(buffer) != 0 {
        m.shownHosts = append(m.shownHosts, buffer)
        m.shownServers = append(m.shownServers, serverbuffer[n])
      }
    }
  }
  if len(m.shownHosts) != 0 {
    m.error = false
  } else {
    m.error = true
  }
  if len(m.shownServers) != 0 {
    m.error = false
  } else {
    m.error = true
  }
  return m, cmd
}

func (m listHostModel) View() string {
  // title
  var header string
  title := "Listing hosts: "
  if m.filtering == true {
    title = out.Style("[" + m.filter + "]", 4, false)
  } else if m.idFilter != "" {
    title = out.Style("[" + m.idFilter + "]", 3, false)
  }
  header += titleStyle.Width(m.width).Render(title)

  // list
  var list string
  var temp []types.HostInfo
  for _, k := range m.hosts {
    temp = append(temp, k...)
  }
  hostInfoConfig := getHostInfoConfig(temp, m.verbose)
  for sn, hosts := range m.shownHosts {
    server := m.shownServers[sn]
    list += serverStyle.Width(m.width).Render("Server " + fmt.Sprint(sn) + ": " + out.Style(server.Name, 2, true) + "@" + out.Style(server.Host, 4, true) + ":" + out.Style(server.DbPort, 5, true)) + "\n"
    for _, host := range hosts {
      string := getString(host, hostInfoConfig)
      list += "   " + string
      list += "\n"
    }
  }

  // menu
  columnWidth := (m.width - 1) / 4
  var menu string

  if m.verbose == false {
    menu += inactiveStyle.Width(columnWidth).BorderRight(true).Render(verboseText)
  } else {
    menu += activeStyle.Copy().Width(columnWidth).Background(lipgloss.Color(out.Color[2])).Foreground(lipgloss.Color("#000000")).BorderRight(true).Render(verboseText)
  }
  if m.onlyOnline == false {
    menu += inactiveStyle.Width(columnWidth).BorderRight(true).Render(onlyOnlineText)
  } else {
    menu += activeStyle.Copy().Width(columnWidth).Background(lipgloss.Color(out.Color[1])).Foreground(lipgloss.Color("#000000")).BorderRight(true).Render(onlyOnlineText)
  }
  if m.idFilter == "" {
    menu += inactiveStyle.Width(columnWidth).Render(idFilterText)
  } else {
    menu += filterStyle.Width(columnWidth).Render(idFilterText)
  }
  if m.filtering == false {
    menu += inactiveStyle.Width(columnWidth).BorderRight(true).Render(filterText)
  } else {
    menu += filterStyle.Width(columnWidth).BorderRight(true).Render(filterText)
  }
  return header + "\n" + list + lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderTop(true).BorderForeground(lipgloss.Color(color[0])).Render(menu)
}

func ListHost(clientConfig types.ClientConfig) {
  var hosts [][]types.HostInfo
  var servers []types.ConnectionInfo
  for _, config := range clientConfig.ConnectionConfigs {
    conn := db.Connect(config)
    buffer := db.GetHosts(conn)
    if len(buffer) != 0 {
      hosts = append(hosts, buffer)
      servers = append(servers, config)
    }
    conn.Close()
  }
  if len(hosts) == 0 {
    out.Warning("No hosts available")
    os.Exit(0)
  }
	m := listHostModel{
    hosts: hosts,
    servers: servers,
    verbose: clientConfig.Verbose,
  }
  tea.NewProgram(m).Run();
  return 
}
