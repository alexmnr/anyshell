package config

import (
	"command"
	"os/exec"
	"out"
	"strings"
	"tools"
)

func CheckAnyshellUpdate() {
	if tools.CheckExist("/tmp/anyshell_update") == true {
		return
	}
	update_needed := false
	// fetch remote
	command_string := "git -C /opt/anyshell remote update"
	cmd := exec.Command(command_string)
	cmd.Run()

	// check if remote is ahead
	command_string = "cd /opt/anyshell && git status"
	err, output, _ := command.Cmd(command_string, false)
	if err == nil {
		if strings.Contains(output, "behind") == true {
			update_needed = true
		}
	}

	// create file to signal that a update is available
	if update_needed == true {
		command_string = "touch /tmp/anyshell_update"
		err, _, _ := command.Cmd(command_string, false)
		if err != nil {
			out.Error("Could not create file")
			return
		}
	}
}

func Update() {
	command.Cmd("cd /opt/anyshell && git pull", true)
	command.Cmd("cd /opt/anyshell && ./install.sh", true)
	command.Cmd("rm -rf /tmp/anyshell_update", false)
}
