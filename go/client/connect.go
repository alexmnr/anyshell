package client

import (
	"command"
	"db"
	"host"
	"libssh"
	"out"

	// "ssh"
	"tui"
	"types"

	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Connect(clientConfig types.ClientConfig) {
  check, hostInfo, connectionInfo, onlyTunnel, forceLocal, verbose := CheckArgs(clientConfig)
  clientConfig.Verbose = verbose
  if check == false {
    hostInfo, connectionInfo, onlyTunnel, forceLocal = tui.SelectHost(clientConfig)
  }
  // if client is on same network
  if hostInfo.PublicIP == host.GetPublicIP() {
    out.Info("Host is on same network, connecting locally...")
    command.Cmd("ssh -o ConnectTimeout=5 " + hostInfo.User + "@" + hostInfo.LocalIP + " -p " + fmt.Sprint(hostInfo.Port), true)
    return
  // if connection is forced local
  } else if forceLocal == true {
    out.Info("Forcing local connect!")
    command.Cmd("ssh -o ConnectTimeout=5 " + hostInfo.User + "@" + hostInfo.LocalIP + " -p " + fmt.Sprint(hostInfo.Port), true)
    return
  }
  requestId := Request(hostInfo, connectionInfo)
  // starting keep alive
  errorCh := make(chan error)
  requestQuitCh := make(chan bool)
  tunnelQuitCh := make(chan bool)
  doneCh := make(chan bool)
  go RequestKeepAlive(requestId, connectionInfo, requestQuitCh)
  // show timer and wait for connection
  go WaitForConnection(hostInfo.ID, connectionInfo, doneCh)
  ret := tui.Timer(10, doneCh)
  if ret == true {
    out.Error("Could not connect to Host!")
    os.Exit(0)
  }
  out.Checkmark("Host " + out.Style("accepted", 1, false) + " request!")
  // if not get tunnel info
  remotePort := getRemotePort(hostInfo.ID, connectionInfo)
  localPort := GetFreeLocalPort(50000)
  serverPort, _ := strconv.Atoi(connectionInfo.SshPort)
  // creating tunnel
  sfunc := func() error {
    tunnelConfig := types.ForwardTunnelConfig{
      Host: connectionInfo.Host,
      User: connectionInfo.Name,
      Password: connectionInfo.Password,
      ServerPort: serverPort,
      LocalPort: localPort,
      RemotePort: remotePort,
    }
    go libssh.CreateForwardTunnel(tunnelConfig, errorCh, tunnelQuitCh)
    test := <- errorCh
    if test != nil {
      out.Warning("Could not create forward tunnel! ")
      out.Error(test)
      return test
    }
    return nil
  }
  tui.RunAction(out.Style("Creating", 4, false) + " ssh tunnel from local port " + fmt.Sprint(localPort) + " to remote port " + fmt.Sprint(serverPort), sfunc, false)
  if onlyTunnel == true {
    out.Info("Created Tunnel on LocalPort: " +  fmt.Sprint(localPort))
    for {

    }
  } else {
    // connect ssh locally
    command.Cmd("ssh " + hostInfo.User + "@localhost -p " + fmt.Sprint(localPort), true)
  }
  tunnelQuitCh <- true
  requestQuitCh <- true
  DeleteRequest(requestId, connectionInfo)
}

func getRemotePort(hostID int, server types.ConnectionInfo) int {
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

func DeleteRequest(id int, server types.ConnectionInfo) {
  conn := db.Connect(server)
  // Stop reverse tunnel
  var query string
  sfunc := func() error {
    query = fmt.Sprintf("DELETE FROM requests WHERE `ID`='%d';", id)
    _, err := conn.Exec(query)
    if err != nil {
      db.QueryError(query, fmt.Sprint(err))
    }
    return nil
  }
  tui.RunAction(out.Style("Deleting", 0, false) + " request with ID: " + fmt.Sprint(id), sfunc, false) 
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

func CheckArgs(config types.ClientConfig) (bool, types.HostInfo, types.ConnectionInfo, bool, bool, bool) {
  args := os.Args
  // hostIndex := 0
  serverIndex := -1
  hostIndex := -1
  var hostInfo types.HostInfo
  var connectionInfo types.ConnectionInfo
  onlyTunnel := false
  forceLocal := false
  verbose := false

  skip := false
  for n, arg := range args {
    if arg == "-v" {
      verbose = true
    }
    if arg == "-t" {
      onlyTunnel = true
    }
    if arg == "-l" {
      forceLocal = true
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
    return false, hostInfo, connectionInfo, onlyTunnel, forceLocal, verbose
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
  return true, hostInfo, connectionInfo, onlyTunnel, forceLocal, verbose
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
