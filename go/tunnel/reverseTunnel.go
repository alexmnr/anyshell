package tunnel

import (
  "types"
  "out"

	"fmt"
	"io"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func handleClient(client net.Conn, remote net.Conn) {
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

func CreateReverseTunnel(config types.ReverseTunnelConfig) error {
  sshConfig := &ssh.ClientConfig{
    User: config.ConnectionInfo.Name,
    Auth: []ssh.AuthMethod{
      ssh.Password(config.ConnectionInfo.Password),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  }

  // Connect to SSH remote server using serverEndpoint
  serverConn, err := ssh.Dial("tcp", endpointString(config.ConnectionInfo.Host, config.ConnectionInfo.SshPort), sshConfig)
  if err != nil {
    out.Error("Dial INTO remote server error: " + fmt.Sprint(err))
    return err
  }

  // Listen on remote server port
  listener, err := serverConn.Listen("tcp", endpointString("localhost", fmt.Sprint(config.RemotePort)))
  if err != nil {
    out.Error("Listen open port ON remote server error: " + fmt.Sprint(err))
    return err
  }
  defer listener.Close()

  // handle incoming connections on reverse forwarded tunnel
  go func() {
    for {
      // Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
      local, err := net.Dial("tcp", endpointString("localhost", fmt.Sprint(config.LocalPort)))
      if err != nil {
        out.Error("Dial INTO local service error: " + fmt.Sprint(err))
        return 
      }

      client, err := listener.Accept()
      if err != nil {
        out.Error(err)
      }

      handleClient(client, local)
    }
  }()

  return nil
}

// func main() {
//         sshConfig := &ssh.ClientConfig{
//                 // SSH connection username
//                 User: "senaex",
//                 Auth: []ssh.AuthMethod{
//                         // put here your private key path
//                         publicKeyFile("/home/senaex/.ssh/id_rsa"),
//                 },
//                 HostKeyCallback: ssh.InsecureIgnoreHostKey(),
//         }

//         // Connect to SSH remote server using serverEndpoint
//         serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)
//         if err != nil {
//                 log.Fatalln(fmt.Printf("Dial INTO remote server error: %s", err))
//         }

//         // Listen on remote server port
//         listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
//         if err != nil {
//                 log.Fatalln(fmt.Printf("Listen open port ON remote server error: %s", err))
//         }
//         defer listener.Close()

//         // handle incoming connections on reverse forwarded tunnel
//         for {
//                 // Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
//                 local, err := net.Dial("tcp", localEndpoint.String())
//                 if err != nil {
//                         log.Fatalln(fmt.Printf("Dial INTO local service error: %s", err))
//                 }

//                 client, err := listener.Accept()
//                 if err != nil {
//                         log.Fatalln(err)
//                 }

//                 handleClient(client, local)
//         }

// }
