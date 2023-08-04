package server

import (
	"out"
	"tools"
	"tui"
  "command"

	"os"
  "time"
	"strings"
  "fmt"
)

var message string
var options []string
var input string

// main entry point for server config
func Menu(){
  // get sudo rights
  command.Cmd("sudo true", true)

  // check for existing servers
  exists := CheckExistingServer()

  // ask for actions
  if exists == false {
    message = "No existing server found"
    options = append(options, out.Style("Create", 1, false) + " new server")
  } else {
    message = "Existing server found"
    options = append(options, out.Style("Replace", 2, false) + " existing server")
    options = append(options, out.Style("Add", 1, false) + " new database to existing server")
    options = append(options, out.Style("Remove", 0, false) + " existing server")
  }
  options = append(options, out.Style("Exit", 0, false))
  ret := tui.Survey(message, options)

  // handle input
  if strings.Contains(ret, "Replace") == true {
    // remove old server
    out.Info("Removing old server...")
    command.SmartCmd("sudo rm -rf /opt/anyshell-server")
    command.Cmd("docker stop anyshell-db && docker stop anyshell-ssh && docker stop anyshell-web", false)
    command.Cmd("docker container rm anyshell-db && docker container rm anyshell-ssh && docker container rm anyshell-web", false)
    ret = out.Style("Create", 1, false) + " new server"
  }
  if strings.Contains(ret, "Exit") == true {
    out.Info("Bye!")
    os.Exit(0)

  ///////// Create new server /////////
  } else if strings.Contains(ret, "Create") == true {
    CheckDependencies()
    // get Information
    serverInfo := tui.GetServerInfo()
    // create dirs
    CreateDirectory(serverInfo)
    // fill docker-compose.yml with data
    FillDockerCompose(serverInfo)
    // ask for edit
    tui.AskEdit("/opt/anyshell-server/docker/docker-compose.yml")
    // Starting docker-compose
    out.Info("Starting docker-compose...")
    command.Cmd("cd /opt/anyshell-server/docker && docker-compose up -d", true)
    // waiting for db to be ready
    fmt.Printf(out.Style("Info: ", 1, true) + "Waiting for db to be ready ")
    for {
      err, _, _ := command.Cmd("docker exec anyshell-db /bin/mariadb -uroot -p" + serverInfo.RootPassword, false)
      if err == nil {
        fmt.Println()
        out.Info("Starting db configuration")
        break
      } else {
        fmt.Printf(".")
      }
    }
    // configurating db
    time.Sleep(1 * time.Second)
    ConfigureDb(serverInfo)
    out.Info("done!")

    out.Warning("You need to forward these ports:\n  " + serverInfo.DbPort + "\n  " + serverInfo.SshPort)

  ///////// Add database to server /////////
  } else if strings.Contains(ret, "Add") == true {
    CheckDependencies()
    // get Information
    dbInfo := tui.GetDbInfo()
    // check if info correct
    CheckDbInfo(dbInfo)
    // create db
    nfo := tui.ServerInfo{
      Name: dbInfo.Name,
      RootPassword: dbInfo.RootPassword,
    }
    out.Info("Starting db configuration")
    time.Sleep(1 * time.Second)
    ConfigureDb(nfo)
    out.Info("done!")

  ///////// remove /////////
  } else if strings.Contains(ret, "Remove") == true {
    CheckDependencies()
    command.SmartCmd("sudo rm -rf /opt/anyshell-server")
    command.Cmd("docker stop anyshell-db && docker stop anyshell-ssh && docker stop anyshell-web", false)
    command.Cmd("docker container rm anyshell-db && docker container rm anyshell-ssh && docker container rm anyshell-web", false)
    command.Cmd("docker image rm docker-ssh", false)
    out.Info("done")
  } else {
    os.Exit(1)
  }
  os.Exit(0)

}

func CheckExistingServer() bool {
  exists := false
  dirs := tools.GetDirs("/opt")
  for _, dir := range dirs {
    if strings.Contains(dir, "anyshell-server") == true {
      exists = true
    }
  }

  return exists
}

func CheckDependencies() {
  if tools.CommandExists("docker") == false {
    out.Error("docker is not installed")
    os.Exit(0)
  }
  if tools.CommandExists("docker-compose") == false {
    out.Error("docker-compose is not installed")
    os.Exit(0)
  }
}

func CreateDirectory(serverInfo tui.ServerInfo) {
  // create dir
  command.SmartCmd("sudo mkdir /opt/anyshell-server")
  // copy configs
  command.SmartCmd("sudo cp -r /opt/anyshell/server/db-config /opt/anyshell-server/")
  command.SmartCmd("sudo cp -r /opt/anyshell/server/sql /opt/anyshell-server/")
  // copy docker files
  if serverInfo.WebInterface == true && tools.GetCPU() == "x86_64" {
    out.Info("Detected " + tools.GetCPU() + " system")
    command.SmartCmd("sudo cp -r /opt/anyshell/server/docker/docker-web /opt/anyshell-server/docker")
  } else if serverInfo.WebInterface == true && tools.GetCPU() != "x86_64" {
    out.Info("Detected " + tools.GetCPU() + " system")
    command.SmartCmd("sudo cp -r /opt/anyshell/server/docker/docker-web-arm /opt/anyshell-server/docker")
  } else if serverInfo.WebInterface == false && tools.GetCPU() == "x86_64" {
    out.Info("Detected " + tools.GetCPU() + " system")
    command.SmartCmd("sudo cp -r /opt/anyshell/server/docker/docker-simple /opt/anyshell-server/docker")
  } else if serverInfo.WebInterface == false && tools.GetCPU() != "x86_64" {
    out.Info("Detected " + tools.GetCPU() + " system")
    command.SmartCmd("sudo cp -r /opt/anyshell/server/docker/docker-simple-arm /opt/anyshell-server/docker")
  } else {
    out.Error("System is not supported yet")
    os.Exit(0)
  }
  // change ownership
  user := tools.GetUser()
  command.SmartCmd("sudo chown " + user + ":" + user + " /opt/anyshell-server -R")
}

func FillDockerCompose(serverInfo tui.ServerInfo) {
  // fill docker-compose.yml with data
  dir := "/opt/anyshell-server/"
  path := dir + "docker/docker-compose.yml"
  read, err := os.ReadFile(path)
  if err != nil {
    out.Error("Cannot read docker-compose.yml")
    os.Exit(1)
  }

  newContents := strings.Replace(string(read), "<dbPort>", serverInfo.DbPort, -1)
  newContents = strings.Replace(newContents, "<sshPort>", serverInfo.SshPort, -1)
  newContents = strings.Replace(newContents, "<webPort>", serverInfo.WebPort, -1)
  newContents = strings.Replace(newContents, "<rootPassword>", serverInfo.RootPassword, -1)
  newContents = strings.Replace(newContents, "<userPassword>", serverInfo.UserPassword, -1)

  // write
  err = os.WriteFile(path, []byte(newContents), 0)
  if err != nil {
    out.Error("Could not write docker-compose.yml")
    os.Exit(1)
  }
  out.Info("Generated docker-compose.yml!")
}
