package out

import (
  "fmt"
)

func Help() {
  fmt.Println("Anyshell - A ssh device manager than can connect to devices behind firewalls using a self hosted server")
  fmt.Println("")
  fmt.Println("Usage:")
  fmt.Println("  anyshell | any               Open main menu")
  fmt.Println("    connect | con | c          Connect to client")
  fmt.Println("      <number>                 Connect to ID <number>")
  fmt.Println("      -l                       Force local connection")
  fmt.Println("      -t                       Create tunnel withouth connecting")
  fmt.Println("      -s <number>              Use server with ID <number>")
  fmt.Println("      -v                       Toggle verbose mode")
  fmt.Println("    list | ls | l              List clients")
  fmt.Println("      -v                       Toggle verbose mode")
  fmt.Println("    host                       Open host menu")
  fmt.Println("      setup                    Setup new host on local device")
  fmt.Println("      edit                     Edit host configuration")
  fmt.Println("      remove                   Remove host from server")
  fmt.Println("      daemon                   Run host daemon (needed for the host to be available)")
  fmt.Println("        service                Run daemon in service mode (no visuals)")
  fmt.Println("    client                     Open client menu")
  fmt.Println("    server                     Open server menu")
  fmt.Println("    help                       Open this help")
  fmt.Println("")
  fmt.Println("Thank you for using anyshell :)")
}
