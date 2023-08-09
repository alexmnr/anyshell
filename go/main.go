package main

import (
	"config"
	"db"
	"out"
	"server"
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
    options = append(options, out.Style("List", 2, false) + "")
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
      hostInfoConfig := db.GetHostInfoConfig(hosts, verbose)
      if n > 0 {
        fmt.Println()
      } else {
        fmt.Println(db.GetHostInfoDescription(hostInfoConfig))
      }
      for _, host := range hosts {
        fmt.Println(db.GetHostInfoString(host, hostInfoConfig))
      }
    }
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
