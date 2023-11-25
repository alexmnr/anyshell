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

// styles
var (
  color = []string{
    "#5f5f5f",
    "#ffffff",
    "#AF82E8",
    "#59656F",
    "#FFAE03",
  }
  
  selectIcon = lipgloss.NewStyle().Foreground(lipgloss.Color(color[4]))

  selectStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color(color[1])).
    BorderStyle(lipgloss.RoundedBorder())

  titleStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color(color[0])).
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color(color[0])).
    BorderBottom(true).
    Align(lipgloss.Center).
    Bold(true)

  serverStyle = titleStyle.Copy().
    Foreground(lipgloss.Color(color[1])).
    Bold(true).
    BorderBottom(false).
    BorderTop(false).
    BorderForeground(lipgloss.Color(color[4])).
    Align(lipgloss.Left)

    // BorderForeground(lipgloss.Color(color[2]))

  inactiveStyle = lipgloss.NewStyle().
    Foreground(lipgloss.Color(color[0])).
    Align(lipgloss.Center).
    BorderStyle(lipgloss.NormalBorder()).
    BorderForeground(lipgloss.Color(color[0])).
    BorderRight(true)

  activeStyle = inactiveStyle.Copy().
    Foreground(lipgloss.Color(color[1])).
    BorderForeground(lipgloss.Color(color[1])).
    Background(lipgloss.Color(color[2]))
  filterStyle = inactiveStyle.Copy().
    Foreground(lipgloss.Color(color[1])).
    BorderForeground(lipgloss.Color(color[3])).
    Background(lipgloss.Color(color[3]))

  // menu entry
  filterText = "(/) filter"
  idFilterText = "(0-9) ID filter"
  tunnelText = "(t) Only tunnel"
  forceLocalText = "(l) Force local"
  verboseText = "(v) verbose Mode"
  onlyOnlineText = "(o) only online"
)

type selectHostModel struct {
  hosts [][]types.HostInfo
  servers []types.ConnectionInfo
  shownHosts [][]types.HostInfo
  shownServers []types.ConnectionInfo
  hostIndex int
  serverIndex int

  // legend string
  filter string
  filtering bool
  idFilter string

  verbose bool
  tunnel bool
  onlyOnline bool
  forceLocal bool

  error bool
  width int

  selectedHost types.HostInfo
  selectedConnection types.ConnectionInfo
}


func (m selectHostModel) Init() tea.Cmd {
  m.shownHosts = m.hosts
  m.shownServers = m.servers
  m.selectedHost = m.hosts[0][0]
  m.selectedConnection = m.servers[0]
  m.filtering = false
  m.verbose = false
  m.onlyOnline = false
  m.forceLocal = false
  m.tunnel = false
  m.hostIndex = 0
  m.serverIndex = 0
  return nil
}
func (m selectHostModel) Down() selectHostModel {
  m.hostIndex += 1
  if m.hostIndex >= len(m.shownHosts[m.serverIndex]) {
    m.hostIndex = 0
    if len(m.shownServers) > 0 {
      if m.serverIndex < len(m.shownServers) - 1 {
        m.serverIndex += 1
      } else {
        m.serverIndex = 0
      }
    } 
  }
  return m
}
func (m selectHostModel) Up() selectHostModel {
  m.hostIndex -= 1
  if m.hostIndex < 0 {
    if len(m.shownServers) > 0 {
      if m.serverIndex > 0 {
        m.serverIndex -= 1
      } else {
        m.serverIndex = len(m.shownServers) - 1
      }
      m.hostIndex = len(m.shownHosts[m.serverIndex]) - 1
    } else {
      m.hostIndex = len(m.shownHosts[0]) - 1
    }
  }
  return m
}

func (m selectHostModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
      m.hostIndex = 0
      m.filtering = false
      m.filter = ""
    }
    // no filtering
    if m.filtering == false {
      switch msg.String() {
      case "/":
        m.filtering = true
        m.idFilter = ""
        m.hostIndex = 0
      case "j":
        m = m.Down()
      case "k":
        m = m.Up()
      }
    // with filtering
    } else {
      if len(msg.String()) == 1 {
        m.filter += msg.String()
        m.filter = regexp.MustCompile(`[^a-z]+`).ReplaceAllString(m.filter, "")
        m.hostIndex = 0
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
    case "t":
      if m.tunnel == true {m.tunnel = false} else {m.tunnel=true}
    case "l":
      if m.forceLocal == true {m.forceLocal = false} else {m.forceLocal=true}
    case "o":
      if m.onlyOnline == true {m.onlyOnline = false} else {m.onlyOnline=true}
      m.hostIndex = 0
      m.serverIndex = 0
    case "q", "ctrl+c":
      m.error = true
      return m, tea.Quit
    case "enter":
      return m, tea.Quit
    case "down":
      m = m.Down()
    case "up":
      m = m.Up()
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
    m.selectedHost = m.shownHosts[m.serverIndex][m.hostIndex]
    m.error = false
  } else {
    m.error = true
  }
  if len(m.shownServers) != 0 {
    m.selectedConnection = m.shownServers[m.serverIndex]
    m.error = false
  } else {
    m.error = true
  }
  return m, cmd
}

func (m selectHostModel) View() string {
  // title
  var header string
  title := "Select host: "
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
    for hn, host := range hosts {
      string := getString(host, hostInfoConfig)
      if hn == m.hostIndex && sn == m.serverIndex {
        if host.Online == true {
          list += selectStyle.Width(m.width - 2).Bold(true).BorderForeground(lipgloss.Color(out.Color[1])).Render("  " + string)
        } else {
          list += selectStyle.Width(m.width - 2).Bold(true).BorderForeground(lipgloss.Color(out.Color[0])).Render("  " + string)
        }
      } else {
        list += "   " + string
      }
      list += "\n"
    }
  }

  // menu
  columnWidth := (m.width - 1) / 6
  var menu string

  if m.tunnel == false {
    menu += inactiveStyle.Width(columnWidth).BorderRight(true).Render(tunnelText)
  } else {
    menu += activeStyle.Width(columnWidth).BorderRight(true).Render(tunnelText)
  }
  if m.forceLocal == false {
    menu += inactiveStyle.Width(columnWidth).BorderRight(true).Render(forceLocalText)
  } else {
    menu += activeStyle.Width(columnWidth).BorderRight(true).Render(forceLocalText)
  }
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

func SelectHost(clientConfig types.ClientConfig) (types.HostInfo, types.ConnectionInfo, bool, bool) {
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
    out.Warning("No hosts to select!")
    os.Exit(0)
  }
	m := selectHostModel{
    hosts: hosts,
    servers: servers,
  }
  tm, err := tea.NewProgram(m).Run();
	if err != nil {
		out.Error("Error running Selector: " + fmt.Sprint(err))
		os.Exit(1)
	}
  mm := tm.(selectHostModel)
  if mm.error == true {
		out.Error("You need to select something!")
    os.Exit(0)
  }
  return mm.selectedHost, mm.selectedConnection, mm.tunnel, mm.forceLocal
}

func getString(host types.HostInfo, config types.HostInfoConfig) string {
  var string string
  if config.Verbose == true {
    des := " %-" + fmt.Sprint(config.IDLength) + "s | %-" + fmt.Sprint(config.NameLength) + "s | %-" + fmt.Sprint(config.UserLength) + "s | %-" + fmt.Sprint(config.PortLength) + "s | %-" + fmt.Sprint(config.LastOnlineLength) + "s | %-" + fmt.Sprint(config.LocalIPLength) + "s | %-" + fmt.Sprint(config.PublicIPLength) + "s | %s "
    string = fmt.Sprintf(des,
    fmt.Sprint(host.ID), host.Name, host.User, fmt.Sprint(host.Port), host.LastOnline, host.LocalIP, host.PublicIP, fmt.Sprint(host.Version))
  } else {
    des := " %-" + fmt.Sprint(config.IDLength) + "s | %-" + fmt.Sprint(config.NameLength) + "s | %-" + fmt.Sprint(config.UserLength) + "s | %-" + fmt.Sprint(config.PortLength) + "s | %s"
    string = fmt.Sprintf(des,
    fmt.Sprint(host.ID), host.Name, host.User, fmt.Sprint(host.Port), host.LastOnline)
  }
  string = colorString(string, host.Online)
  return string
}

func colorString(input string, online bool) string {
  onlineColor := "#6EFA72"
  offlineColor := "#FF0F80"
  if online == true {
    return lipgloss.NewStyle().Foreground(lipgloss.Color(onlineColor)).Render(input)
  } else {
    return lipgloss.NewStyle().Foreground(lipgloss.Color(offlineColor)).Render(input)
  }
}

func getHostInfoConfig(hosts []types.HostInfo, verbose bool) types.HostInfoConfig {
  config := types.HostInfoConfig{
    Verbose: verbose,
    IDLength: 2,
    NameLength: 4,
    UserLength: 4,
    PortLength: 4,
    PublicIPLength: 9,
    LocalIPLength: 8,
    LastOnlineLength: 11,
  }
  for _, host := range hosts {
    if len(fmt.Sprint(host.ID)) > config.IDLength {
      config.IDLength = len(fmt.Sprint(host.ID))
    }
    if len(host.Name) > config.NameLength {
      config.NameLength = len(host.Name)
    }
    if len(host.Name) > config.NameLength {
      config.NameLength = len(host.Name)
    }
    if len(host.User) > config.UserLength {
      config.UserLength = len(host.User)
    }
    if len(fmt.Sprint(host.Port)) > config.PortLength {
      config.PortLength = len(fmt.Sprint(host.Port))
    }
    if len(host.PublicIP) > config.PublicIPLength {
      config.PublicIPLength = len(host.PublicIP)
    }
    if len(host.LocalIP) > config.LocalIPLength {
      config.LocalIPLength = len(host.LocalIP)
    }
    if len(host.LastOnline) > config.LastOnlineLength {
      config.LastOnlineLength = len(host.LastOnline)
    }
  }
  return config
}
