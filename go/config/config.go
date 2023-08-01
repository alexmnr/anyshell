package config

import (
	"tools"
)

type ClientConfig struct {
}

func ClientConfigCheck() bool {
  // check if necessary directory exists
  found := true
  homeDir := tools.GetHomeDir()
  configDir := homeDir + "/.config/anyshell"
  if tools.CheckExist(configDir) == false {
    found = false
  } else {
    if tools.CheckExist(configDir + "/config.yml") == false {
      found = false
    }
  }
  return found
}

func CreateClientConfig() {

}
