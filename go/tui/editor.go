package tui

import (
  "out"
  "command"

	"os"
)

func AskEdit(file string) {
  // Ask for edit
  message := "Do you want to edit the file? (" + file + ")"
  var options []string
  options = append(options, out.Style("No", 0, false))
  options = append(options, out.Style("Yes", 1, false))
  ret := Survey(message, options)
  if ret == options[1] {
    Edit(file)
  } 
}

func Edit(file string) {
  editor := os.Getenv("EDITOR")
  if editor == "" {
    editor = "vim"
  }
  command.Cmd(editor + " " + file, true)
}
