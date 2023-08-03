package main

import (
	"config"
	"out"
	"tui"
  "server"
  "tools"
  "command"

	"strings"
	"os"
)
var message string
var options []string
var ret string
var clientConfig config.ClientConfig

func main() {
  ///////// Config /////////
  check := config.ClientConfigCheck()

  ///////// Menu /////////
  if check == false {
    options = append(options, out.Style("Client", 4, false) + " configuration")
  } else {
    options = append(options, out.Style("Client", 5, false) + " configuration")
    // load config
    clientConfig = config.GetClientConfig()
  }
  options = append(options, out.Style("Server", 5, false) + " configuration")
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
