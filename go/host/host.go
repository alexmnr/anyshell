package host

import (
	"command"
	"config"
	"db"
	"out"
	"strconv"
	"tools"
  "tui"

  "strings"
  "regexp"
	"database/sql"
	"fmt"
	"os"
)

func Setup(conn *sql.DB, clientConfig config.ClientConfig) {
  command.Cmd("sudo true", true)
  // get data from host
  var info db.HostInfo
  sfunc := func() error {
    info = GetHostParameters(conn)
    return nil
  }
  tui.RunAction("Gathering information about host...", sfunc, false)
  out.Info(info)
}

func GetHostParameters(conn *sql.DB) db.HostInfo {
  info := db.HostInfo{}
  info.ID = GetID(conn)
  info.Name = tools.GetHostName()
  info.User = tools.GetUser()
  info.Port = GetSSHPort()
  info.Online = false
  info.LocalIP = GetLocalIP()
  info.PublicIP = GetPublicIP()
  info.Version = GetVersion()

  return info
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
  err, output, _ := command.Cmd("git -C /opt/anyshell rev-list --count master", false);
  if err != nil {
    return 0
  }
  clean := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(output, "")
  clean = strings.Replace(clean, " ", "", -1)
  version, _ := strconv.Atoi(clean)
  return version
}
