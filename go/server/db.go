package server

import (
  "tui"
  "command"
  "out"

  "os"
  "strings"
)

func CheckDbInfo(dbInfo tui.DbInfo) {
  // check if db is accessible
  err, _, _ := command.Cmd("docker exec anyshell-db true", false)
  if err != nil {
    out.Error("Database is not accessible!")
    os.Exit(1)
  }
  // check if root password is correct
  err, _, _ = command.Cmd("docker exec anyshell-db /bin/mariadb -uroot -p" + dbInfo.RootPassword, false)
  if err != nil {
    out.Error("Root password is wrong!")
    os.Exit(1)
  }
  // check if db already exists
  query := "/bin/mariadb -u root -p" + dbInfo.RootPassword + " -e \"SHOW DATABASES;\""
  err, output, _ := command.Cmd("docker exec anyshell-db " + query, false)
  if strings.Contains(output, dbInfo.Name) == true {
    out.Error("A Database with this name already exists!")
    os.Exit(1)
  }
}

