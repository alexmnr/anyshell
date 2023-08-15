package out

import (
  "github.com/charmbracelet/lipgloss"
  "fmt"
)

var Color [6]string = [6] string {
  "#FA4453",
  "#6EFA73",
  "#FFAE03",
  "#FF0F80",
  "#04D1F1",
  "#59656F",
}

func Style(input string, style int, bold bool) string {
  if style < len(Color) {
    return lipgloss.NewStyle().Foreground(lipgloss.Color(Color[style])).Bold(bold).Render(input)
  } else {
    return input
  }
}

func Error(error interface{}) {
  string := fmt.Sprint(error)
  fmt.Println(Style("Error: ", 0, true) + string)
}
func Info(info interface{}) {
  string := fmt.Sprint(info)
  fmt.Println(Style("Info: ", 1, true) + string)
}
func Warning(info interface{}) {
  string := fmt.Sprint(info)
  fmt.Println(Style("Warning: ", 2, true) + string)
}

func CommandError(command string, err error, out string, error string) {
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(Color[0])).Bold(true).PaddingLeft(0).Render("Error running Command: ") + command) 
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(Color[2])).Bold(false).PaddingLeft(1).Render("Error Code: ") + fmt.Sprint(err)) 
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(Color[2])).Bold(false).PaddingLeft(1).Render("Command stdout: ") + out) 
  fmt.Println(lipgloss.NewStyle().Foreground(lipgloss.Color(Color[2])).Bold(false).PaddingLeft(1).Render("Command stderr: ") + error) 
}

func TestColors() {
  for i := 0; i < len(Color); i++ {
    fmt.Println(Style("This is Color: " + fmt.Sprint(i), i, false))
    fmt.Println(Style("This is Color: " + fmt.Sprint(i) + " in bold", i, true))
  }
}
