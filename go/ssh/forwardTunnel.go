package ssh

import (
  "types"
  "out"

  "github.com/elliotchance/sshtunnel"
	libssh "golang.org/x/crypto/ssh"

  "fmt"
  "time"
  "os"
  "log"
)


func sshString(user string, host string, port string) string {
  string := user + "@" + host + ":" + port
  return string
}

func CreateTunnel(config types.ForwardTunnelConfig, ch chan error, quit chan bool) {
  // Setup the tunnel, but do not yet start it yet.
  out.Info(config.LocalPort)
  out.Info(config.RemotePort)
  tunnel, err := sshtunnel.NewSSHTunnel(
    sshString(config.ConnectionInfo.Name, config.ConnectionInfo.Host, config.ConnectionInfo.SshPort),
    libssh.Password(config.ConnectionInfo.Password),
    // The destination host and port of the actual server.
    "localhost:" + fmt.Sprint(config.RemotePort),
    fmt.Sprint(config.LocalPort),
  )
  // You can provide a logger for debugging, or remove this line to
  // make it silent.
  tunnel.Log = log.New(os.Stdout, "", log.Ldate | log.Lmicroseconds)
  ch <- nil
  if err != nil {
    out.Error("Could not create tunnel, error: " + fmt.Sprint(err))
    ch <- err
    return
  }
  // Start the server in the background. You will need to wait a
  // small amount of time for it to bind to the localhost port
  // before you can start sending connections.
  err = tunnel.Start()
  time.Sleep(100 * time.Millisecond)
  ch <- err
  for {
    select {
    case <- quit:
      tunnel.Close()
      return
    }
  }
}
