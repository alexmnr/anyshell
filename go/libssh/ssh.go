package libssh

import (
	"out"
	"strings"
	"types"

	"bytes"
	"fmt"

	"golang.org/x/crypto/ssh"
)

func endpointString(host string, port string) string {
  return fmt.Sprintf("%s:%s", host, port)
}

func GetFreeRemotePort(connectionInfo types.ConnectionInfo, start int) int {
  port := start

  err, output := CommandInSSH(connectionInfo, "netstat -tunlp")
  if err != nil {
    out.Error("Could not get free remote port")
    return 0
  }
  for {
    if strings.Contains(output, ":" + fmt.Sprint(port)) == false {
      return port
    } else {
      port++
    }
    if port >= 51000 {
      out.Error("Could not get free remote port")
      return 0
    }
  }
}

func CommandInSSH(connectionInfo types.ConnectionInfo, command string) (error, string) {
  sshConfig := &ssh.ClientConfig{
    User: connectionInfo.Name,
    Auth: []ssh.AuthMethod{
      ssh.Password(connectionInfo.Password),
    },
    HostKeyCallback: ssh.InsecureIgnoreHostKey(),
  }

  // Connect to SSH remote server using serverEndpoint
  serverConn, err := ssh.Dial("tcp", endpointString(connectionInfo.Host, connectionInfo.SshPort), sshConfig)
  if err != nil {
    out.Error("Dial INTO remote server error: " + fmt.Sprint(err))
    return err, ""
  }

  session, err := serverConn.NewSession()
  if err != nil {
    return err, ""
  }
  defer session.Close()
  var b bytes.Buffer  // import "bytes"
  session.Stdout = &b // get output

  // Finally, run the command
  err = session.Run(command)
  if err != nil {
    out.Error("Error running Command: " + command + " in ssh")
    return err, ""
  }
  return nil, b.String()
}
