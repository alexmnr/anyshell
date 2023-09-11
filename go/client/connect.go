package client

import (
	"command"
	"db"
	"out"
	// "ssh"
	"tui"
	"types"

	"time"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Connect(clientConfig types.ClientConfig) {
  check, hostInfo, connectionInfo, onlyTunnel := CheckArgs(clientConfig)
  if check == false {
    hostInfo, connectionInfo, onlyTunnel = tui.SelectHost(clientConfig)
  }
  requestId := Request(hostInfo, connectionInfo)
  // starting keep alive
  errCh := make(chan error)
  quitCh := make(chan bool)
  doneCh := make(chan bool)
  go RequestKeepAlive(requestId, connectionInfo, quitCh)
  // show timer and wait for connection
  go WaitForConnection(hostInfo.ID, connectionInfo, doneCh)
  ret := tui.Timer(10, doneCh)
  if ret == true {
    out.Error("Could not connect to Host!")
    os.Exit(0)
  }
  out.Checkmark("Host " + out.Style("accepted", 1, false) + " request!")
  // creating tunnel
  serverPort := getServerPort(hostInfo.ID, connectionInfo)
  localPort := GetFreeLocalPort(50000)
  tunnelConfig := types.ForwardTunnelConfig{
    ConnectionInfo: connectionInfo,
    LocalPort: localPort,
    RemotePort: serverPort,
  }
  _ = errCh
  _ = tunnelConfig
  for {
    out.Info("running")
    time.Sleep(1 * time.Second)
  }
  // go ssh.CreateTunnel(tunnelConfig, errCh, quitCh)
  // err := <-errCh
  // out.Info(err)
  // err = <-errCh
  // out.Info(err)
  // sfunc := func() error {
  //   tunnelConfig := types.ForwardTunnelConfig{
  //     ConnectionInfo: connectionInfo,
  //     LocalPort: localPort,
  //     RemotePort: serverPort,
  //   }
  //   go ssh.CreateTunnel(tunnelConfig, errCh, quitCh)
  //   err := <-errCh
  //   if err != nil {
  //     out.Error(err)
  //     return err
  //   }
  //   return nil
  // }
  // tui.RunAction(out.Style("Creating", 4, false) + " ssh tunnel from local port " + fmt.Sprint(localPort) + " to remote port " + fmt.Sprint(serverPort), sfunc, false)
  _ = onlyTunnel
}

func getServerPort(hostID int, server types.ConnectionInfo) int {
  conn := db.Connect(server)
  var serverPort int
  query := fmt.Sprintf("SELECT ServerPort FROM connections WHERE HostID=%d;", hostID)
  queryErr := conn.QueryRow(query).Scan(&serverPort)
  if queryErr != sql.ErrNoRows {
    return serverPort
  } else {
    out.Error("Could not get Server Port!")
    os.Exit(1)
    return -1
  }
}

func WaitForConnection(hostID int, server types.ConnectionInfo, quit chan bool)  {
  conn := db.Connect(server)
  for {
    query := fmt.Sprintf("SELECT 1 FROM connections WHERE HostID=%d;", hostID)
    queryErr := conn.QueryRow(query).Scan()
    if queryErr != sql.ErrNoRows {
      quit <- true
      quit <- true
      conn.Close()
      return
    }
  }
}

func Request(hostInfo types.HostInfo, server types.ConnectionInfo) int {
  conn := db.Connect(server)
  // requesting host
  var requestId int
  sfunc := func() error {
    // get unique request ID
    requestId = db.GetID(conn, "requests")
    // create request
    query := fmt.Sprintf("INSERT INTO requests (`ID`, `HostID`) VALUES ('%d', '%d');", requestId, hostInfo.ID)
    _, err := conn.Exec(query)
    if err != nil {
      db.QueryError(query, fmt.Sprint(err))
    }
    return nil
  }
  tui.RunAction(out.Style("Requesting", 1, false) + " host with ID: " + fmt.Sprint(hostInfo.ID), sfunc, false)

  conn.Close()
  return requestId
}

func RequestKeepAlive(id int, server types.ConnectionInfo, quit chan bool) {
  conn := db.Connect(server)
  for {
    select {
    default:
      query := fmt.Sprintf("UPDATE requests SET `LastUsed`=CURRENT_TIMESTAMP WHERE ID=%d;", id)
      _, err := conn.Exec(query)
      if err != nil {
        db.QueryError(query, fmt.Sprint(err))
        return
      }
    case <- quit:
      conn.Close()
      return
    }
  }
}

func CheckArgs(config types.ClientConfig) (bool, types.HostInfo, types.ConnectionInfo, bool) {
  args := os.Args
  // hostIndex := 0
  serverIndex := -1
  hostIndex := -1
  var hostInfo types.HostInfo
  var connectionInfo types.ConnectionInfo
  onlyTunnel := false

  skip := false
  for n, arg := range args {
    if arg == "-t" {
      onlyTunnel = true
    }
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
    return false, hostInfo, connectionInfo, false
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
  return true, hostInfo, connectionInfo, onlyTunnel
}


func GetFreeLocalPort(start int) int {
  port := start
  err, output, _ := command.Cmd("netstat -tunlp", false)
  if err != nil {
    out.Error("Could not get free local port")
    return 0
  } else {
    if strings.Contains(output, ":" + fmt.Sprint(port)) == false {
      return port
    } else {
      port++
    }
    if port >= start + 1000 {
      out.Error("Could not get free local port")
      return 0
    }
  }
  return port
}
