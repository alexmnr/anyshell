package tools

import (
	"out"
	"strings"

	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
)

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func CheckExist(file string) bool {
  if _, err := os.Stat(file); os.IsNotExist(err) {
    return false
  } else {
    return true
  }
}

func GetOS() string {
  check := CommandExists("apt")
  if check == true{
    return "debian"
  }
  check = CommandExists("pacman")
  if check == true{
    return "arch"
  }
  return "?"
}
func GetInstaller() string {
  check := CommandExists("apt")
  installer := ""
  if check == true{
    installer = "sudo apt update && sudo apt install -y "
  }
  check = CommandExists("pacman")
  if check == true{
    installer = "sudo pacman --needed --noconfirm -Sy "
  }
  return installer
}
func GetUser() string {
  user, _ := user.Current()
  name := user.Username
  return name
}
func GetHomeDir() string {
  dir, err := os.UserHomeDir()
  if err != nil {
    out.Error("Could not get HomeDir: " + fmt.Sprint(err))
    os.Exit(1)
  }
  return dir
}
func GetDirs(location string) []string {
  var dirs []string

  items, _ := ioutil.ReadDir(location)
  for _, item := range items {
    if item.IsDir() {
      dirs = append(dirs, item.Name())
    }
  }
  return dirs
}

func GetCPU() string {
  cpu := ""
  if CommandExists("lscpu") == false {
    out.Error("Could not detect cpu architecture")
    os.Exit(1)
  }
  cmd := exec.Command("lscpu")
  out, _ := cmd.Output()
  if strings.Contains(string(out), "x86_64") == true {
    cpu = "x86_64"
  } else if strings.Contains(string(out), "aarch64") == true {
    cpu = "aarch64"
  }
  return cpu
}
