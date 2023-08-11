package host

import (
	"command"
	"db"
	"out"
	"tunnel"
	"tui"
	"types"

	"database/sql"
	"fmt"
	"time"
)

var sfunc func() error

func Daemon(config types.HostConfig) {
  out.Info("Started Daemon for Host: " + fmt.Sprint(config))
  // check if ssh start stop is activated
  for {
    conn := db.Connect(config.Server)
    // update online status
    sfunc = func() error {
      UpdateHostOnDB(config, conn)
      return nil
    }
    tui.RunAction("Updating Database", sfunc, false)

    // check for requests
    query := fmt.Sprintf("SELECT requests.ID, hosts.ID, hosts.Port FROM requests, hosts WHERE requests.`HostID`=hosts.ID AND hosts.Name='%s' AND hosts.User='%s' AND hosts.Port='%s';", config.Name, config.User, fmt.Sprint(config.Port))
    rows, err := conn.Query(query)
    if err != nil {
      db.QueryError(query, fmt.Sprint(err))
    }
    for rows.Next() {
      // found reqeust
      var requestID int
      var hostID int
      var hostPort int
      err := rows.Scan(&requestID, &hostID, &hostPort)
      if err == nil {
        // start ssh server
        command.Cmd("sudo systemctl start sshd.service", false)
        // remove old request from database
        query := fmt.Sprintf("DELETE FROM requests WHERE `ID`='%s';", fmt.Sprint(requestID))
        _, err := conn.Query(query)
        if err != nil {
          db.QueryError(query, fmt.Sprint(err))
        }
        // create reverse tunnel
        reverseTunnelConfig := types.ReverseTunnelConfig{
          ConnectionInfo: config.Server,
          LocalPort: GetSSHPort(),
          RemotePort: tunnel.GetFreeRemotePort(config.Server, 50000),
        }
        tunnel.CreateReverseTunnel(reverseTunnelConfig)
         
        out.Info("Requested with ID: " + fmt.Sprint(requestID))
      }
    }

    conn.Close()
    time.Sleep(5 * time.Second)
  }
}

func UpdateHostOnDB(hostConfig types.HostConfig, conn *sql.DB) {
  // get data from host
  localIP := GetLocalIP()
  publicIP := GetPublicIP()
  version := fmt.Sprint(GetVersion())

  query := fmt.Sprintf("UPDATE hosts SET `LastOnline`=CURRENT_TIMESTAMP, `Online`='1', `LocalIP`='%s', `PublicIP`='%s', `Version`='%s' WHERE Name='%s' AND User='%s' AND Port='%s';",
  localIP, publicIP, version, hostConfig.Name, hostConfig.User, fmt.Sprint(hostConfig.Port))
  _, err := conn.Query(query)
  if err != nil {
    db.QueryError(query, fmt.Sprint(err))
  }
    // string version = exec("git -C /opt/anyshell rev-list --count master");
    // sprintf(sql_query,
    //         "UPDATE hosts "
    //         "SET "
    //         "`last-online`=CURRENT_TIMESTAMP, "
    //         "`online`='1', "
    //         "`localIP`='%s', "
    //         "`publicIP`='%s', "
    //         "`version`='%i' "
    //         "WHERE Name='%s';",
    //         user_details->localIP, user_details->publicIP, atoi(version.c_str()), user_details->hostname);
    // res = mysql_run(conn, sql_query);
    // mysql_free_result(res);
  
}

