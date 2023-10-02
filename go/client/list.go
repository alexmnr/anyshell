package client

import (
  "types"
  "tui"

  "fmt"
  "github.com/charmbracelet/lipgloss"
)


func List(clientConfig types.ClientConfig) {
  tui.SelectHost(clientConfig)
}

func GetHostInfoConfig(hosts []types.HostInfo, verbose bool) types.HostInfoConfig {
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
    if len(str(host.ID)) > config.IDLength {
      config.IDLength = len(str(host.ID))
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

func GetHostInfoDescription(config types.HostInfoConfig) string {
  var string string
  if config.Verbose == true {
    des := " %-" + str(config.IDLength) + "s | %-" + str(config.NameLength) + "s | %-" + str(config.UserLength) + "s | %-" + str(config.PortLength) + "s | %-" + str(config.LastOnlineLength) + "s | %-" + str(config.LocalIPLength) + "s | %-" + str(config.PublicIPLength) + "s | %s "
    string = fmt.Sprintf(des,
    "ID", "Name", "User", "Port", "last-online", "local-IP", "public-IP", "version")
  } else {
    des := " %-" + str(config.IDLength) + "s | %-" + str(config.NameLength) + "s | %-" + str(config.UserLength) + "s | %-" + str(config.PortLength) + "s | %s"
    string = fmt.Sprintf(des,
    "ID", "Name", "User", "Port", "last-online")
  }
  return string
}

func GetHostInfoString(host types.HostInfo, config types.HostInfoConfig) string {
  var string string
  if config.Verbose == true {
    des := " %-" + str(config.IDLength) + "s | %-" + str(config.NameLength) + "s | %-" + str(config.UserLength) + "s | %-" + str(config.PortLength) + "s | %-" + str(config.LastOnlineLength) + "s | %-" + str(config.LocalIPLength) + "s | %-" + str(config.PublicIPLength) + "s | %s "
    string = fmt.Sprintf(des,
    fmt.Sprint(host.ID), host.Name, host.User, fmt.Sprint(host.Port), host.LastOnline, host.LocalIP, host.PublicIP, fmt.Sprint(host.Version))
  } else {
    des := " %-" + str(config.IDLength) + "s | %-" + str(config.NameLength) + "s | %-" + str(config.UserLength) + "s | %-" + str(config.PortLength) + "s | %s"
    string = fmt.Sprintf(des,
    fmt.Sprint(host.ID), host.Name, host.User, fmt.Sprint(host.Port), host.LastOnline)
  }
  string = color(string, host.Online)
  return string
}

func color(input string, online bool) string {
  onlineColor := "#6EFA72"
  offlineColor := "#FF0F80"
  if online == true {
    return lipgloss.NewStyle().Foreground(lipgloss.Color(onlineColor)).Render(input)
  } else {
    return lipgloss.NewStyle().Foreground(lipgloss.Color(offlineColor)).Render(input)
  }
}


func str(input interface{}) string {
  return fmt.Sprint(input)
}
