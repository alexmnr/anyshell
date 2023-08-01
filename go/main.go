package main

import (
	"config"
	"out"
	"tui"
  "server"

	"strings"
	"os"
)
var message string
var options []string
var ret string

func main() {
  ///////// Config /////////
  check := config.ClientConfigCheck()

  ///////// Menu /////////
  if check == false {
    options = append(options, out.Style("Client", 4, false) + " config")
  } else {
    //TODO get config
  }
  options = append(options, out.Style("Server", 3, false) + " config")
  options = append(options, out.Style("Exit", 0, false))
  message = "Welcome to anyshell!"

  ret = tui.Survey(message, options)
  
  // Exit 
  if strings.Contains(ret, "Exit") {
    out.Info("Bye!")
    os.Exit(0)
  // Server Config 
  } else if strings.Contains(ret, "Server") {
    server.Menu()
  // Client Config 
  } else if strings.Contains(ret, "Client") {
    config.CreateClientConfig()
  }
}
