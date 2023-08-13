package db

import (
	"out"
  "types"

	"database/sql"
	"fmt"
  "os"
	_ "github.com/go-sql-driver/mysql"
)


func Connect(config types.ConnectionInfo) *sql.DB {
  connString := config.Name  + ":" + config.Password + "@tcp(" + config.Host + ":" + config.DbPort + ")/" + config.Name 
  db, err := sql.Open("mysql", connString)
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

func ConnectError(error string, config types.ConnectionInfo) {
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

func GetID(conn *sql.DB, database string) int {
  id := 0
  query := "SELECT ID FROM " + database + " ORDER BY `ID` ASC;"
  rows, err := conn.Query(query)
  if err != nil {
    QueryError(query, fmt.Sprint(err))
  }
  defer rows.Close()

  for rows.Next() {
    var check int
    err := rows.Scan(&check)
    if err != nil {
      out.Error(err)
      os.Exit(0)
    }
    if check == id {
      id++
    } else {
      return id
    }
  }
  return id
}

