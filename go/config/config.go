package config

import (
	"command"
	"out"
	"tools"
	"tui"
	"types"

	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

var message string
var options []string
var ret string

func Menu() {
  check := ClientConfigCheck()
  if check == false {
    CreateClientConfig()
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
      AddConnectionConfig()
      out.Info("Succesfully edited client config!")
    } else if strings.Contains(ret, "Remove") {
      configDir := "/etc/anyshell"
      command.Cmd("rm -f " + configDir + "/client-config.yml", false)
      out.Info("Succesfully removed client config!")
    } else if strings.Contains(ret, "Edit") {
      configDir := "/etc/anyshell"
      tui.Edit(configDir + "/client-config.yml")
      out.Info("Succesfully edited client config!")
    }
  }
}

func ClientConfigCheck() bool {
  // check if necessary directory exists
  found := true
  configDir := "/etc/anyshell"
  if tools.CheckExist(configDir) == false {
    found = false
  } else {
    if tools.CheckExist(configDir + "/client-config.yml") == false {
      found = false
    }
  }
  return found
}

func CreateClientConfig() {
  command.Cmd("sudo true", true)
  configDir := "/etc/anyshell"
  command.Cmd("sudo mkdir " + configDir, false)
  command.Cmd("sudo chmod a+rw " + configDir, false)

  var connectionConfigs []types.ConnectionInfo
  connectionConfig := tui.GetConnectionInfo()
  connectionConfigs = append(connectionConfigs, connectionConfig)

  clientConfig := types.ClientConfig{
    ConnectionConfigs: connectionConfigs,
  }
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

func AddConnectionConfig() {
  configDir := "/etc/anyshell"

  clientConfig := GetClientConfig()
  connectionConfig := tui.GetConnectionInfo()
  clientConfig.ConnectionConfigs = append(clientConfig.ConnectionConfigs, connectionConfig)

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

func GetClientConfig() types.ClientConfig {
  configDir := "/etc/anyshell"
  yamlFile, _ := os.ReadFile(configDir + "/client-config.yml")
  clientConfig := types.ClientConfig{}
  if err := yaml.Unmarshal(yamlFile, &clientConfig); err != nil {
    fmt.Printf(out.Style("Error while reading config: ", 0, false) + "%v \n", err)
    os.Exit(1)
  }
  return clientConfig
}
