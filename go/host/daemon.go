package host

import (
	"command"
	"db"
	"out"
	"libssh"
	"tui"
	"types"

  "strings"
	"database/sql"
	"fmt"
	"time"
  "sync"
  "strconv"
)

var sfunc func() error

func Daemon(config types.HostConfig, service bool, wg *sync.WaitGroup) {
  defer wg.Done()
  conn := db.Connect(config.Server)
  var hosting bool
  var remotePort int
  var connectionID int
  hosting = false
  errorCh := make(chan error, 100)
  quitCh := make(chan bool)
  // check if ssh start stop is activated
  for {
    // update online status
    // sfunc = func() error {
    //   UpdateHost(config, conn)
    //   return nil
    // }
    // if service == false {tui.RunAction(out.Style("Updating Database", 5, false), sfunc, false)} else {fmt.Println(out.Style("Updating Database", 5, false)); sfunc()}
    // if service == false {tui.RunAction(out.Style("Updating Database", 5, false), sfunc, false)} else {fmt.Println(out.Style("Updating Database", 5, false)); sfunc()}
    UpdateHost(config, conn)

    // check for requests
    query := fmt.Sprintf("SELECT 1 FROM requests WHERE HostID=%d;", config.ID)
    queryErr := conn.QueryRow(query).Scan()
    if queryErr != sql.ErrNoRows {
      if hosting == false {
        ///////// Starting /////////
        // start ssh server
        checkSSHServer(service)
        // create reverse tunnel
        sfunc = func() error {
          // localPort := GetSSHPort()
          remotePort = libssh.GetFreeRemotePort(config.Server, 50000)        
          sshPort, _ := strconv.Atoi(config.Server.SshPort)
          reverseTunnelConfig := types.ReverseTunnelConfig{
            User: config.Server.Name,
            Host: config.Server.Host,
            ServerPort: sshPort,
            Password: config.Server.Password,
            LocalPort: config.Port,
            RemotePort: remotePort,
          }
          go libssh.CreateReverseTunnel(reverseTunnelConfig, errorCh, quitCh)
          test := <- errorCh
          if test != nil {
            hosting = false
          } else {
            hosting = true
          }
          return test
        }
        if service == false {tui.RunAction(out.Style("Creating", 1, false) + " reverse tunnel", sfunc, false)} else {fmt.Println(out.Style("Creating", 1, false) + " reverse tunnel"); sfunc()}
        if hosting == false {
          out.Error("Could not create reverse tunnel!")
          continue
        }
        // create connection entry
        sfunc = func() error {
          connectionID = db.GetID(conn, "connections")
          query := fmt.Sprintf("INSERT INTO connections (`ID`, `HostID`, `ServerPort`) VALUES ('%d', '%d', '%d');", connectionID, config.ID, remotePort)
          _, err := conn.Exec(query)
          if err != nil {
            db.QueryError(query, fmt.Sprint(err))
          }
          return nil
        }
        if service == false {tui.RunAction(out.Style("Creating", 1, false) + " connection entry", sfunc, false)} else {fmt.Println(out.Style("Creating", 1, false) + " connection entry"); sfunc()}

        hosting = true
      } else {
        ///////// Updating /////////
        sfunc = func() error {
          UpdateConnection(config, conn)
          return nil
        }
        if service == false {tui.RunAction(out.Style("Keeping connection alive", 5, false), sfunc, false)} else {fmt.Println(out.Style("Keeping connection alive", 5, false)); sfunc()}
      }
      ///////// Stopping /////////
    } else if queryErr == sql.ErrNoRows {
      if hosting == true {
        // stop ssh server if necesary
        if config.SSHStartStop == true {
          sfunc = func() error {
            command.Cmd("sudo systemctl stop sshd.service", false)
            return nil
          }
          if service == false {tui.RunAction(out.Style("Stopping", 0, false) + " ssh server", sfunc, false)} else {fmt.Println(out.Style("Stopping", 0, false) + " ssh server"); sfunc()}
        }
        // Stop reverse tunnel
        sfunc = func() error {
          quitCh <- true
          query = fmt.Sprintf("DELETE FROM requests WHERE `ID`='%d';", connectionID)
          _, err := conn.Exec(query)
          if err != nil {
            db.QueryError(query, fmt.Sprint(err))
          }
          return nil
        }
        if service == false {tui.RunAction(out.Style("Deleting", 0, false) + " connection with ID: " + fmt.Sprint(connectionID), sfunc, false)} else {fmt.Println(out.Style("Deleting", 0, false) + " connection with ID: " + fmt.Sprint(connectionID)); sfunc()}
        hosting = false
      }
    }

    time.Sleep(2 * time.Second)
  }
}

func UpdateHost(hostConfig types.HostConfig, conn *sql.DB) {
  // get data from host
  localIP := GetLocalIP()
  publicIP := GetPublicIP()
  version := fmt.Sprint(GetVersion())

  query := fmt.Sprintf("UPDATE hosts SET `LastOnline`=CURRENT_TIMESTAMP, `Online`='1', `LocalIP`='%s', `PublicIP`='%s', `Version`='%s' WHERE Name='%s' AND User='%s' AND Port='%s';",
  localIP, publicIP, version, hostConfig.Name, hostConfig.User, fmt.Sprint(hostConfig.Port))
  _, err := conn.Exec(query)
  if err != nil {
    db.QueryError(query, fmt.Sprint(err))
  }
}

func UpdateConnection(hostConfig types.HostConfig, conn *sql.DB) {
  // updating timestamp
  query := fmt.Sprintf("UPDATE connections SET `LastUsed`=CURRENT_TIMESTAMP WHERE ID=%d;", hostConfig.ID)
  _, err := conn.Exec(query)
  if err != nil {
    db.QueryError(query, fmt.Sprint(err))
  }
}

func checkSSHServer(service bool) {
  _, output, _ := command.Cmd("sudo systemctl is-active sshd.service", false)
  if strings.Contains(output, "inactive") == true {
    sfunc = func() error {
      command.Cmd("sudo systemctl start sshd.service", false)
      return nil
    }
    if service == false {tui.RunAction(out.Style("Starting", 2, false) + " ssh server", sfunc, false)} else {fmt.Println(out.Style("Starting", 2, false) + " ssh server"); sfunc()}
  }
}
