package ssh

import (
  "types"
  "out"

	"fmt"
	"io"
	"log"
	"net"

	libssh "golang.org/x/crypto/ssh"
)

func HandleClient(client net.Conn, remote net.Conn) {
  defer client.Close()
  chDone := make(chan bool)

  // Start remote -> local data transfer
  go func() {
    _, err := io.Copy(client, remote)
    if err != nil {
      log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
    }
    chDone <- true
  }()

  // Start local -> remote data transfer
  go func() {
    _, err := io.Copy(remote, client)
    if err != nil {
      log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
    }
    chDone <- true
  }()

  <-chDone
}

func endpointString(host string, port string) string {
  return fmt.Sprintf("%s:%s", host, port)
}

func CreateReverseTunnel(config types.ReverseTunnelConfig, ch chan error, quit chan bool) {
  // quit := false
  sshConfig := &libssh.ClientConfig{
    User: config.ConnectionInfo.Name,
    Auth: []libssh.AuthMethod{
      libssh.Password(config.ConnectionInfo.Password),
    },
    HostKeyCallback: libssh.InsecureIgnoreHostKey(),
  }

  // Connect to SSH remote server using serverEndpoint
  serverConn, err := libssh.Dial("tcp", endpointString(config.ConnectionInfo.Host, config.ConnectionInfo.SshPort), sshConfig)
  if err != nil {
    out.Error("Dial INTO remote server error: " + fmt.Sprint(err))
    ch <- err
  }
  defer serverConn.Close()

  // Listen on remote server port
  listener, err := serverConn.Listen("tcp", endpointString("localhost", fmt.Sprint(config.RemotePort)))
  if err != nil {
    out.Error("Listen open port ON remote server error: " + fmt.Sprint(err))
    ch <- err
  }
  defer listener.Close()

  go listenerQuitCheck(listener, serverConn, quit)
  // handle incoming connections on reverse forwarded tunnel
  for {
    // Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
    local, err := net.Dial("tcp", endpointString("localhost", fmt.Sprint(config.LocalPort)))
    if err != nil {
      out.Error("Dial INTO local service error: " + fmt.Sprint(err))
      ch <- err
    }

    ch <- nil
    client, err := listener.Accept()
    if err != nil {
      break
    }

    HandleClient(client, local)
  }
}

func listenerQuitCheck(listener net.Listener, serverConn *libssh.Client, quit chan bool) {
  for {
    select {
    case <- quit:
      listener.Close()
      serverConn.Close()
      return
    }
  }
}
