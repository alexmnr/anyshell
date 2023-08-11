package host

import (
	"db"
	"out"
	"tools"
	"tui"
	"types"

	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

func RemoveHostFromLocalConfig(server types.ConnectionInfo) {
  // read old config
  homeDir := tools.GetHomeDir()
  configDir := homeDir + "/.config/anyshell"
  yamlFile, _ := os.ReadFile(configDir + "/client-config.yml")
  clientConfig := types.ClientConfig{}
  if err := yaml.Unmarshal(yamlFile, &clientConfig); err != nil {
    fmt.Printf(out.Style("Error while reading config: ", 0, false) + "%v \n", err)
    os.Exit(1)
  }

  oldConfig := clientConfig.HostConfigs
  for n, k := range clientConfig.HostConfigs {
    if k.Server == server {
      clientConfig.HostConfigs = append(oldConfig[:n], oldConfig[n+1:]...)
    }
  }

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

func RemoveHost(server types.ConnectionInfo) {
  conn := db.Connect(server)
  // deleting from database
  sfunc = func() error {
    query := fmt.Sprintf("DELETE FROM hosts WHERE `Name`='%s' AND `User`='%s' AND `Port`='%s';", tools.GetHostName(), tools.GetUser(), fmt.Sprint(GetSSHPort()))
    _, err := conn.Query(query)
    if err != nil {
      db.QueryError(query, fmt.Sprint(err))
      os.Exit(0)
    }
    return nil
  }
  tui.RunAction("Deleting Host from Database", sfunc, false)

  // delete from local config
  sfunc = func() error {
    RemoveHostFromLocalConfig(server)
    return nil
  }
  tui.RunAction("Deleting Host from local config", sfunc, false)

  out.Info("Done!")
}
