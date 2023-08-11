package host

import (
	"command"
	"db"
	"out"
	"strconv"
	"tools"
	"tui"
	"types"
  "server"

	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"
)
var (
  message string
  options []string
  ret string
)

func Menu (clientConfig types.ClientConfig) {
  args := os.Args
  for _, arg := range args {
    if arg == "setup" {
      ret = "Setup"
    } else if arg == "edit" {
      ret = "Edit"
    } else if arg == "remove" {
      ret = "Remove"
    } else if arg == "daemon" {
      ret = "Daemon"
    }
  }

  if ret == "" {
    options = append(options, out.Style("Setup", 4, false) + " new host")
    options = append(options, out.Style("Edit", 2, false) + " configuration")
    options = append(options, out.Style("Remove", 3, false) + " host from server")
    options = append(options, out.Style("Daemon", 1, false) + " start manually")
    options = append(options, out.Style("Exit", 0, false))
    message = "Host Configuration"
    ret = tui.Survey(message, options)
  }
  // host setup
  if strings.Contains(ret, "Setup") {
    connectionInfo := server.SelectConnection(clientConfig)
    Setup(connectionInfo)
  // Edit
  } else if strings.Contains(ret, "Edit") {
    homeDir := tools.GetHomeDir()
    configDir := homeDir + "/.config/anyshell"
    tui.Edit(configDir + "/client-config.yml")
    out.Info("Succesfully edited host config!")
  // Remove
  } else if strings.Contains(ret, "Remove") {
    connectionInfo := server.SelectConnection(clientConfig)
    RemoveHost(connectionInfo)
  // Daemon
  } else if strings.Contains(ret, "Daemon") {
    if len(clientConfig.HostConfigs) == 0 {
      out.Warning("No hosts configured, run 'anyshell host setup")
    }
    for _, hostConfig := range clientConfig.HostConfigs {
      Daemon(hostConfig)
    }
  // Exit
  } else {
    out.Info("Bye!")
    os.Exit(0)
  }

}


func GetID(conn *sql.DB) int {
  id := 0
  query := "SELECT ID FROM hosts ORDER BY `ID` ASC;"
  rows, err := conn.Query(query)
  if err != nil {
    db.QueryError(query, fmt.Sprint(err))
  }
  defer rows.Close()

  for rows.Next() {
    var check int
    err := rows.Scan(&check)
    if err != nil {
      out.Error(err)
      os.Exit(0)
    }
    if check == id {
      id++
    } else {
      return id
    }
  }
  return id
}

func GetSSHPort() int {
  err, output, _ := command.Cmd("cat /etc/ssh/sshd_config | grep 'Port '", false)
  if err != nil {
    return 22
  }
  str := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(output, "")
  port, _ := strconv.Atoi(str)
  return port
}

func GetLocalIP() string {
  err, output, _ := command.Cmd("ip -o -4  address show | awk ' NR==2 { gsub(/\\/.*/, \"\", $4); print $4 } '", false)
  if err != nil {
    return "0.0.0.0"
  }
  clean := regexp.MustCompile(`[^0-9.]+`).ReplaceAllString(output, "")
  clean = strings.Replace(clean, " ", "", -1)
  return clean
}

func GetPublicIP() string {
  err, output, _ := command.Cmd("(curl -s ifconfig.me | grep -o -Eo '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}')", false)
  if err != nil {
    return "0.0.0.0"
  }
  clean := regexp.MustCompile(`[^0-9.]+`).ReplaceAllString(output, "")
  return clean

}

func GetVersion() int {
  err, output, _ := command.Cmd("git -C /opt/anyshell rev-list --count main", false);
  if err != nil {
    return 0
  }
  clean := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(output, "")
  clean = strings.Replace(clean, " ", "", -1)
  version, _ := strconv.Atoi(clean)
  return version
}



