package libssh

import (
  "types"
  "out"

	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func HandleReverseClient(client net.Conn, remote net.Conn) {
  // defer client.Close()

  // Start remote -> local data transfer
  go func() {
    _, err := io.Copy(client, remote)
    if err != nil {
      log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
    }
  }()

  // Start local -> remote data transfer
  go func() {
    _, err := io.Copy(remote, client)
    if err != nil {
      log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
    }
  }()
}

func CreateReverseTunnel(config types.ReverseTunnelConfig, ch chan error, quit chan bool) {
  // quit := false
  sshConfig := &ssh.ClientConfig{
    User: config.User,
    Auth: []ssh.AuthMethod{
      ssh.Password(config.Password),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  }

  // Connect to SSH remote server using serverEndpoint
  serverConn, err := ssh.Dial("tcp", endpointString(config.Host, fmt.Sprint(config.ServerPort)), sshConfig)
  if err != nil {
    out.Error("Dial INTO remote server error: " + fmt.Sprint(err))
    ch <- err
    return
  }

  // Listen on remote server port
  listener, err := serverConn.Listen("tcp", endpointString("localhost", fmt.Sprint(config.RemotePort)))
  if err != nil {
    out.Error("Listen open port ON remote server error: " + fmt.Sprint(err))
    ch <- err
    return
  }

  go listenerQuitCheck(listener, serverConn, quit)
  // handle incoming connections on reverse forwarded tunnel
  for {
    // Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
    local, err := net.Dial("tcp", endpointString("localhost", fmt.Sprint(config.LocalPort)))
    if err != nil {
      out.Error("Dial INTO local service error: " + fmt.Sprint(err))
      ch <- err
      continue
    }

    ch <- nil
    client, err := listener.Accept()
    if err != nil {
      out.Warning("Stopped reverse tunnel")
      break
    }
    out.Info("Accepted client")

    HandleReverseClient(client, local)
  }
}

func listenerQuitCheck(listener net.Listener, serverConn *ssh.Client, quit chan bool) {
  for {
    select {
    case <- quit:
      listener.Close()
      serverConn.Close()
      return
    }
  }
}
