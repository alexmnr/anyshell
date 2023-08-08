package host

import (
	"command"
	"config"
	"db"
	"out"
	"strconv"
	"tools"
	"tui"
	"types"

	"database/sql"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

var sfunc func() error

func Daemon(config types.HostConfig) {
  // check if ssh start stop is activated
  for {
    // check for requests
    conn := db.Connect(config.Server)
    query := fmt.Sprintf("SELECT requests.ID, hosts.ID, hosts.Port FROM requests, hosts WHERE requests.`HostID`=hosts.ID AND hosts.Name='%s';", config.Name)
    rows, err := conn.Query(query)
    if err != nil {
      db.QueryError(query, fmt.Sprint(err))
    }
    found := false
    for rows.Next() {
      found = true
    }
    conn.Close()

    out.Info(found)
    time.Sleep(5 * time.Second)
  }
}

func Setup(server types.ConnectionInfo) {
  conn := db.Connect(server)
  command.Cmd("sudo true", true)
  // check dependencies 
  sfunc = func() error {
    var missing string
    // TODO Tmux session integration
    // if tools.CommandExists("tmux") == false {
    //   missing += "tmux"
    // }
    if tools.CommandExists("/usr/bin/sshd") == false {
      missing += "openssh-server"
    }

    if missing != "" {
      out.Error("Missing dependencies: " + missing)
      return errors.New(":(")
    }
    return nil
  }
  tui.RunAction("Checking dependencies", sfunc, false)
  // get data from host
  var info types.HostInfo
  sfunc = func() error {
    info = GetHostParameters(conn)
    return nil
  }
  tui.RunAction("Gathering information about host", sfunc, false)

  // check if host already exists locally
  sfunc = func() error {
    check := CheckLocalConcfig(info, server)
    if check == true {
      out.Error("This Host already exists in local Config!")
      return errors.New("This Host already exists in local Config!")
    }
    return nil
  }
  tui.RunAction("Checking if host already exists", sfunc, false)

  // allow edit of data
  command.Cmd("rm -f /tmp/hostSetup.yml", false)
  yamlData, err := yaml.Marshal(&info)
  if err != nil {
    out.Error(err)
    os.Exit(0)
  }
  fileName := "/tmp/hostSetup.yml"
  err = os.WriteFile(fileName, yamlData, 0644)
  if err != nil {
    fmt.Printf(out.Style("Error while writing Config: ", 0, false) + "%v \n", err)
    os.Exit(1)
  }
  // ask for edit
  tui.AskEdit("/tmp/hostSetup.yml")
  // read edited config
  yamlFile, _ := os.ReadFile("/tmp/hostSetup.yml")
  hostInfo := types.HostInfo{}
  if err := yaml.Unmarshal(yamlFile, &hostInfo); err != nil {
    fmt.Printf(out.Style("Error while reading config: ", 0, false) + "%v \n", err)
    os.Exit(1)
  }

  // write info to database
  sfunc = func() error {
    err := AddHostTODB(hostInfo, conn)
    return err
  }
  tui.RunAction("Adding host to database", sfunc, false)

  // add host to config
  sfunc = func() error {
    AddHostToConfig(hostInfo, server)
    return nil
  }
  tui.RunAction("Adding host to local config", sfunc, false)
  
  out.Info("done!")
}

func GetHostParameters(conn *sql.DB) types.HostInfo {
  info := types.HostInfo{}
  info.ID = GetID(conn)
  info.Name = tools.GetHostName()
  info.User = tools.GetUser()
  info.Port = GetSSHPort()
  info.Online = false
  info.LocalIP = GetLocalIP()
  info.PublicIP = GetPublicIP()
  info.LastOnline = fmt.Sprint(time.Now().UTC())
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
  err, output, _ := command.Cmd("git -C /opt/anyshell rev-list --count main", false);
  if err != nil {
    return 0
  }
  clean := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(output, "")
  clean = strings.Replace(clean, " ", "", -1)
  version, _ := strconv.Atoi(clean)
  return version
}

func AddHostToConfig(info types.HostInfo, server types.ConnectionInfo) {
  // read old config
  homeDir := tools.GetHomeDir()
  configDir := homeDir + "/.config/anyshell"
  yamlFile, _ := os.ReadFile(configDir + "/client-config.yml")
  clientConfig := types.ClientConfig{}
  if err := yaml.Unmarshal(yamlFile, &clientConfig); err != nil {
    fmt.Printf(out.Style("Error while reading config: ", 0, false) + "%v \n", err)
    os.Exit(1)
  }

  // create config
  hostConfig := types.HostConfig{
    Server: server,
    Name: info.Name,
    User: info.User,
    Port: info.Port,
  }
  // append to new one
  clientConfig.HostConfigs = append(clientConfig.HostConfigs, hostConfig)

  // write file
  yamlData, err := yaml.Marshal(&clientConfig)
  if err != nil {
    out.Error(err)
    os.Exit(0)
  }
  fileName := configDir + "/client-config.yml"
  err = os.WriteFile(fileName, yamlData, 0644)
  if err != nil {
    fmt.Printf(out.Style("Error while writing Config: ", 0, false) + "%v \n", err)
    os.Exit(1)
  }
}

func AddHostTODB(hostInfo types.HostInfo, conn *sql.DB) error {
  query := fmt.Sprintf("INSERT INTO hosts (`ID`, `Name`, `User`, `Port`, `PublicIP`, `LocalIP`, `Online`, `Version`) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '1', '%s');",
  fmt.Sprint(hostInfo.ID), hostInfo.Name, hostInfo.User, fmt.Sprint(hostInfo.Port), hostInfo.PublicIP, hostInfo.LocalIP, fmt.Sprint(hostInfo.Version))
  _, err := conn.Query(query)
  if err != nil {
    db.QueryError(query, fmt.Sprint(err))
    return err
  }
  return nil
}

func CheckLocalConcfig(info types.HostInfo, server types.ConnectionInfo) bool {
  conf := config.GetClientConfig()
  if len(conf.HostConfigs) == 0 {
    return false
  }
  for _, host := range conf.HostConfigs {
    if host.Name != info.Name {return false}
    if host.User != info.User {return false}
    if host.Port != info.Port {return false}
    if host.Server != server {return false}
  }
  return true
}

func DeleteHost(hostInfo types.HostInfo) {
  return
}
