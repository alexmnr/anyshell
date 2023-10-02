package libssh

import (
  "types"
  "out"

  "fmt"
  "golang.org/x/crypto/ssh"
  "io"
  "net"
)

func CreateForwardTunnel(config types.ForwardTunnelConfig, ch chan error, quit chan bool) {
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

  // create listener on local port
  listener, err := net.Listen("tcp", endpointString("localhost", fmt.Sprint(config.LocalPort)))
  if err != nil {
    out.Error("Opening Listener ON local port error: " + fmt.Sprint(err))
    ch <- err
    return
  }

  go listenerQuitCheck(listener, serverConn, quit)
  ch <- nil

  for {
    localConn, err := listener.Accept()
    if err != nil {
      // out.Warning("Stopped forward tunnel")
      break
    }
    // out.Info("Accepted client")
    // tunnel.logf("accepted connection")
    remoteConn, err := serverConn.Dial("tcp", endpointString("localhost", fmt.Sprint(config.RemotePort)))
    if err != nil {
        out.Error("remote dial error: " + fmt.Sprint(err))
        return
    }
    copyConn := func(writer, reader net.Conn) {
        _, err := io.Copy(writer, reader)
        if err != nil {
            out.Error("io.Copy error: " + fmt.Sprint(err))
        }
    }
    go copyConn(localConn, remoteConn)
    go copyConn(remoteConn, localConn)
  }
}
