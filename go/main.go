package main

import (
	"config"
	"db"
	"host"
	"out"
	"server"
	"tui"
	"types"
  "client"

	"fmt"
	"os"
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

  //////// Arguments ///////
  args := os.Args
  if check == true {
    for _, arg := range args {
      if arg == "-v" {
        verbose = true
      } else if arg == "-vv" {
        verbose = true
      } else if arg == "list" || arg == "ls" || arg == "l" {
        ret = "List"
      } else if arg == "connect" || arg == "con" || arg == "c" {
        ret = "Connect"
      } else if arg == "host" {
        ret = "Host"
      } else if arg == "setup" {
        ret += " setup"
      } else if arg == "setup" {
        ret += " setup"
      } else if arg == "daemon" {
        ret += " daemon"
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
    options = append(options, out.Style("Connect", 1, false) + " to host")
    options = append(options, out.Style("List", 2, false) + " hosts")
    options = append(options, out.Style("Host", 4, false) + " configuration")
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
    for n, connection := range clientConfig.ConnectionConfigs {
      conn := db.Connect(connection)
      hosts := db.GetHosts(conn)
      hostInfoConfig := client.GetHostInfoConfig(hosts, verbose)
      if n > 0 {
        fmt.Println()
      } else {
        fmt.Println(client.GetHostInfoDescription(hostInfoConfig))
      }
      for _, host := range hosts {
        fmt.Println(client.GetHostInfoString(host, hostInfoConfig))
      }
    }
  // connect
  } else if strings.Contains(ret, "Connect") {
    check, hostInfo, connectionInfo := client.CheckArgs(clientConfig)
    if check == false {
      hostInfo, connectionInfo = tui.SelectHost(clientConfig)
    }
    out.Info(hostInfo)
    out.Info(connectionInfo)
  // host
  } else if strings.Contains(ret, "Host") {
    host.Menu(clientConfig)
    // host menu
  // Server Config 
  } else if strings.Contains(ret, "Server") {
    server.Menu()
  // Client Config 
  } else if strings.Contains(ret, "Client") {
    config.Menu()
  }
}
