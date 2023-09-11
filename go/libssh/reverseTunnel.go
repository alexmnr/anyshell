package libssh

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
  sshConfig := &libssh.ClientConfig{
    User: config.User,
    Auth: []libssh.AuthMethod{
      libssh.Password(config.Password),
    },
    HostKeyCallback: libssh.InsecureIgnoreHostKey(),
  }

  // Connect to SSH remote server using serverEndpoint
  serverConn, err := libssh.Dial("tcp", endpointString(config.Host, fmt.Sprint(config.ServerPort)), sshConfig)
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




// package libssh

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"net"
// 	"out"
// 	"types"

// 	"golang.org/x/crypto/ssh"
// )

// type Endpoint struct {
//   Host string
//   Port int
// }

// func (endpoint *Endpoint) String() string {
//   return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
// }

// func handleClient(client net.Conn, remote net.Conn) {
//   chDone := make(chan bool)

//   // Start remote -> local data transfer
//   go func() {
//     _, err := io.Copy(client, remote)
//     if err != nil {
//       log.Println(fmt.Sprintf("error while copy remote->local: %s", err))
//     }
//     chDone <- true
//   }()

//   // Start local -> remote data transfer
//   go func() {
//     _, err := io.Copy(remote, client)
//     if err != nil {
//       log.Println(fmt.Sprintf("error while copy local->remote: %s", err))
//     }
//     chDone <- true
//   }()

//   <-chDone
// }


// // local service to be forwarded
// var localEndpoint = Endpoint{
//   Host: "localhost",
//   Port: 22,
// }

// // remote SSH server
// var serverEndpoint = Endpoint{
//   Host: "10.8.77.91",
//   Port: 41999,
// }

// // remote forwarding port (on remote SSH server network)
// var remoteEndpoint = Endpoint{
//   Host: "localhost",
//   Port: 1111,
// }

// func CreateReverseTunnel(config types.ReverseTunnelConfig) {
//   // refer to https://godoc.org/golang.org/x/crypto/ssh for other authentication types
//   sshConfig := &ssh.ClientConfig{
//     // SSH connection username
//     User: "senaex",
//     Auth: []ssh.AuthMethod{
//       ssh.Password("Iafneo-4523"),
//     },
//     HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//   }

//   // Connect to SSH remote server using serverEndpoint
//   serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
//   if err != nil {
//     log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
//   }
//   // defer serverConn.Close()

//   // Listen on remote server port
//   listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
//   if err != nil {
//     log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
//   }
//   // defer listener.Close()

//   // handle incoming connections on reverse forwarded tunnel
//   go func() {
//     for {
//       // Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
//       local, err := net.Dial("tcp", localEndpoint.String())
//       if err != nil {
//         log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
//       }

//       client, err := listener.Accept()
//       if err != nil {
//         out.Warning("yeeeeppppp")
//         continue
//       } else {
//         out.Info("yeeeeppppp")
//       }

//       handleClient(client, local)
//     }
//   }()
// }
