package db

import (
	"out"
	"tui"

	"database/sql"
	"fmt"
  "os"
	_ "github.com/go-sql-driver/mysql"
)

type HostInfo struct {
  ID int
  Name string
  User string
  Port int
  Online bool
  PublicIP string
  LocalIP string
  LastOnline string
  Version int
}

func Connect(config tui.ConnectionInfo) *sql.DB {
  db, err := sql.Open("mysql", "anyshell:user@tcp(localhost:42998)/anyshell")
	// defer db.Close()
  if err != nil {
    ConnectError(fmt.Sprint(err), config)
    os.Exit(0)
  }
  if Check(db) == false {
    ConnectError("database not responding", config)
    os.Exit(0)
  }
  return db
}

func ConnectError(error string, config tui.ConnectionInfo) {
  fmt.Println(out.Style("Can't connect to Database:", 0, false))
  fmt.Println(out.Style("  host: ", 0, false) + config.Host)
  fmt.Println(out.Style("  user: ", 0, false) + config.Name)
  fmt.Println(out.Style("  port: ", 0, false) + config.DbPort)
  fmt.Println(out.Style("  Error: ", 0, false) + error)
}

func QueryError(query string, error string) {
  fmt.Println(out.Style("SQL Error:", 0, false))
  fmt.Println(out.Style("  Query: ", 0, false) + query)
  fmt.Println(out.Style("  Error: ", 0, false) + error)
}

func Check(db *sql.DB) bool {
	var version string
  db.QueryRow("SELECT VERSION();").Scan(&version)
  if len(version) == 0 {
    return false
  } else {
    return true
  }
}
