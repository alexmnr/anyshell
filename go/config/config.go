package config

import (
	"command"
	"tools"
  "out"
  "tui"
  "types"

  "os"
  "fmt"
	"gopkg.in/yaml.v2"
)



func ClientConfigCheck() bool {
  // check if necessary directory exists
  found := true
  homeDir := tools.GetHomeDir()
  configDir := homeDir + "/.config/anyshell"
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
  homeDir := tools.GetHomeDir()
  configDir := homeDir + "/.config/anyshell"
  command.Mkdir(configDir, false)

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
  homeDir := tools.GetHomeDir()
  configDir := homeDir + "/.config/anyshell"

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
  homeDir := tools.GetHomeDir()
  configDir := homeDir + "/.config/anyshell"
  yamlFile, _ := os.ReadFile(configDir + "/client-config.yml")
  clientConfig := types.ClientConfig{}
  if err := yaml.Unmarshal(yamlFile, &clientConfig); err != nil {
    fmt.Printf(out.Style("Error while reading config: ", 0, false) + "%v \n", err)
    os.Exit(1)
  }
  return clientConfig
}
