package main

import (
	"command"
	"config"
	"db"
	"out"
	"server"
	"tools"
	"tui"
  "host"
  "types"

	"os"
  "fmt"
	"strings"
)
var message string
var options []string
var ret string
var verbose bool
var clientConfig types.ClientConfig

func main() {
  ///////// Config /////////
  check := config.ClientConfigCheck()
  if check == true {
    // load config
    clientConfig = config.GetClientConfig()
  }
  verbose = false

  host.Daemon(clientConfig.HostConfigs[0])
  os.Exit(0)
  //////// Arguments ///////
  args := os.Args
  if check == true {
    for _, arg := range args {
      if arg == "-v" {
        verbose = true
      } else if arg == "-vv" {
        verbose = true
      } else if arg == "list" {
        ret = "List"
      } else if arg == "host" {
        ret = "Host"
      } else if arg == "setup" {
        ret += " setup"
      }
    }
  } else {
    if len(args) > 0 {
      out.Warning("You need to configure client first!")
    }
  }
  ///////// Menu /////////
  if check == false {
    options = append(options, out.Style("Client", 4, false) + " configuration")
  } else {
    options = append(options, out.Style("List", 2, false) + " hosts")
    options = append(options, out.Style("Host", 4, false) + " setup")
    // if host.HostConfigCheck() == true {
    //   options = append(options, out.Style("Host", 5, false) + " configuration")
    // }
    options = append(options, out.Style("Client", 5, false) + " configuration")
  }
  options = append(options, out.Style("Server", 5, false) + " configuration")
  options = append(options, out.Style("Exit", 0, false))
  message = "Welcome to anyshell!"

  if len(ret) == 0 {
    ret = tui.Survey(message, options)
  }
  
  // Exit 
  if strings.Contains(ret, "Exit") {
    out.Info("Bye!")
    os.Exit(0)
  // list
  } else if strings.Contains(ret, "List") {
    conn := db.Connect(clientConfig.ConnectionConfigs[0])
    hosts := db.GetHosts(conn)
    hostInfoConfig := db.GetHostInfoConfig(hosts, verbose)
    fmt.Println(db.GetHostInfoDescription(hostInfoConfig))
    for _, host := range hosts {
      fmt.Println(db.GetHostInfoString(host, hostInfoConfig))
    }
  // host
  } else if strings.Contains(ret, "Host") {
    // host setup
    if strings.Contains(ret, "setup") {
      // TODO select server
      host.Setup(clientConfig.ConnectionConfigs[0])
    }
  // Server Config 
  } else if strings.Contains(ret, "Server") {
    server.Menu()
  // Client Config 
  } else if strings.Contains(ret, "Client") {
    if check == false {
      config.CreateClientConfig()
      out.Info("Succesfully created client config!")
    } else {
      options = nil
      options = append(options, out.Style("Add", 4, false) + " another connection")
      options = append(options, out.Style("Edit", 2, false) + " configuration")
      options = append(options, out.Style("Remove", 3, false) + " configuration")
      options = append(options, out.Style("Exit", 0, false))
      message = "Client Configuration"

      ret = tui.Survey(message, options)
      if strings.Contains(ret, "Exit") {
        out.Info("Bye!")
        os.Exit(0)
      } else if strings.Contains(ret, "Add") {
        config.AddConnectionConfig()
        out.Info("Succesfully edited client config!")
      } else if strings.Contains(ret, "Remove") {
        homeDir := tools.GetHomeDir()
        configDir := homeDir + "/.config/anyshell"
        command.Cmd("rm -f " + configDir + "/client-config.yml", false)
        out.Info("Succesfully removed client config!")
      } else if strings.Contains(ret, "Edit") {
        homeDir := tools.GetHomeDir()
        configDir := homeDir + "/.config/anyshell"
        tui.Edit(configDir + "/client-config.yml")
        out.Info("Succesfully edited client config!")
      }
    }
  }
}
