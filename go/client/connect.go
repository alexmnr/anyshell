package client

import (
	"types"
  "out"
  "db"

	"os"
	"strconv"
)


func CheckArgs(config types.ClientConfig) (bool, types.HostInfo, types.ConnectionInfo) {
  args := os.Args
  // hostIndex := 0
  serverIndex := -1
  hostIndex := -1
  var hostInfo types.HostInfo
  var connectionInfo types.ConnectionInfo

  skip := false
  for n, arg := range args {
    if arg == "-s" {
      if len(args) - 1 <= n {
        out.Error("You need to specify a server index")
        os.Exit(0)
      } else {
        var err error
        serverIndex, err = strconv.Atoi(args[n + 1])
        if err != nil {
          out.Error("Only numbers are allowed!")
          os.Exit(0)
        }
        skip = true
      }
    } else if _, err := strconv.Atoi(arg); err == nil {
      if skip == true {
        skip = false
      } else {
        hostIndex, _ = strconv.Atoi(arg)
      }
    }
  }

  if hostIndex == -1 {
    return false, hostInfo, connectionInfo
  } else {
    if serverIndex == -1 {
      serverIndex = 0
    }
  }
  if len(config.ConnectionConfigs) <= serverIndex {
    out.Error("Server index out of range!")
    os.Exit(0)
  }
  connectionInfo = config.ConnectionConfigs[serverIndex]
  conn := db.Connect(connectionInfo)
  hosts := db.GetHosts(conn)
  if len(hosts) <= hostIndex {
    out.Error("Host index out of range!")
    os.Exit(0)
  }
  hostInfo = hosts[hostIndex]
  return true, hostInfo, connectionInfo
}
