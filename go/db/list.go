package db

import (
	"out"

	"database/sql"
  "time"
	"fmt"
  "os"
	_ "github.com/go-sql-driver/mysql"
)

func GetHosts(db *sql.DB) []HostInfo {
  query := "SELECT * FROM hosts ORDER BY `ID` ASC;"
  rows, err := db.Query(query)
  if err != nil {
    QueryError(query, fmt.Sprint(err))
  }
  defer rows.Close()

  var hosts []HostInfo

  for rows.Next() {
    var host HostInfo
    err := rows.Scan(&host.ID, &host.Name, &host.User, &host.Port, &host.Online, &host.PublicIP, &host.LocalIP, &host.LastOnline, &host.Version)
    t, err := time.Parse("2006-01-02 15:04:05", host.LastOnline)
    if err != nil {
      out.Error(err)
      os.Exit(0)
    }
    host.LastOnline = TimeDiffString(time.Now(), t)
    if err != nil {
      QueryError(query, fmt.Sprint(err))
      os.Exit(1)
    }
    hosts = append(hosts, host)
  }
  return hosts
}
func TimeDiffString(a, b time.Time) string {
    var output string
    year, month, day, hour, min, sec := TimeDiff(a, b)
    if year != "0" {
      output = year + "Y " + month + "M"
    } else if month != "0" {
      output = month + "M " + day + "D"
    } else if day != "0" {
      output = day + "D " + hour + "h"
    } else if hour != "0" {
      output = hour + "h " + min + "m"
    } else if min != "0" {
      output = min + "m " + sec + "s"
    } else {
      output = sec + "s "
    }
    return output
}

func TimeDiff(a, b time.Time) (year, month, day, hour, min, sec string) {
    if a.Location() != b.Location() {
        b = b.In(a.Location())
    }
    if a.After(b) {
        a, b = b, a
    }
    y1, M1, d1 := a.Date()
    y2, M2, d2 := b.Date()

    h1, m1, s1 := a.Clock()
    h2, m2, s2 := b.Clock()

    iyear  := int(y2 - y1)
    imonth := int(M2 - M1)
    iday   := int(d2 - d1)
    ihour  := int(h2 - h1)
    imin   := int(m2 - m1)
    isec   := int(s2 - s1)

    // Normalize negative values
    if isec < 0 {
        isec += 60
        imin--
    }
    if imin < 0 {
        imin += 60
        ihour--
    }
    if ihour < 0 {
        ihour += 24
        iday--
    }
    if iday < 0 {
        // days in month:
        t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
        iday += 32 - t.Day()
        imonth--
    }
    if imonth < 0 {
        imonth += 12
        iyear--
    }

    year  = fmt.Sprint(iyear)
    month = fmt.Sprint(imonth)
    day   = fmt.Sprint(iday)
    hour  = fmt.Sprint(ihour)
    min   = fmt.Sprint(imin)
    sec   = fmt.Sprint(isec)

    return
}

