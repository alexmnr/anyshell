package main

import (
	"client"
	"config"
	"host"
	"out"
	"server"
	"tui"
	"types"

	"os"
	"os/signal"
	"strings"
	"syscall"
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
    options = append(options, out.Style("Connect", 1, false))
    options = append(options, out.Style("List", 2, false))
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
    client.List(clientConfig)
  // connect
  } else if strings.Contains(ret, "Connect") {
    client.Connect(clientConfig)
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

// handle ctrl-c
func init() {
    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        os.Exit(1)
    }()
}

