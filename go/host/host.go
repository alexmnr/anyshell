package host

import (
	"command"
	"out"
	"server"
	"strconv"
	// "time"
	"tui"
	"types"

	"os"
	"regexp"
	"strings"
	"sync"
)

var (
	message string
	options []string
	ret     string
	wg      sync.WaitGroup
)

func Menu(clientConfig types.ClientConfig) {
	service := false
	args := os.Args
	for _, arg := range args {
		if arg == "setup" {
			ret = "Setup"
		} else if arg == "edit" {
			ret = "Edit"
		} else if arg == "remove" {
			ret = "Remove"
		} else if arg == "daemon" {
			ret = "Daemon"
		} else if arg == "service" {
			service = true
		}
	}

	if ret == "" {
		options = append(options, out.Style("Setup", 4, false)+" new host")
		options = append(options, out.Style("Edit", 2, false)+" configuration")
		options = append(options, out.Style("Remove", 3, false)+" host from server")
		options = append(options, out.Style("Daemon", 1, false)+" start manually")
		options = append(options, out.Style("Exit", 0, false))
		message = "Host Configuration"
		ret = tui.Survey(message, options)
	}
	// host setup
	if strings.Contains(ret, "Setup") {
		connectionInfo := server.SelectConnection(clientConfig)
		Setup(connectionInfo)
		// Edit
	} else if strings.Contains(ret, "Edit") {
		configDir := "/etc/anyshell"
		tui.Edit(configDir + "/client-config.yml")
		out.Info("Succesfully edited host config!")
		// Remove
	} else if strings.Contains(ret, "Remove") {
		connectionInfo := server.SelectConnection(clientConfig)
		RemoveHost(connectionInfo)
		// Daemon
	} else if strings.Contains(ret, "Daemon") {
		if len(clientConfig.HostConfigs) == 0 {
			out.Warning("No hosts configured, run 'anyshell host setup")
		}
		for _, hostConfig := range clientConfig.HostConfigs {
			// TODO full multi host support
			wg.Add(1)
			go Daemon(hostConfig, service, &wg)
		}
		// for {time.Sleep(5 * time.Second)}
		wg.Wait()
		// Exit
	} else {
		out.Info("Bye!")
		os.Exit(0)
	}

}

func GetSSHPort() int {
	err, output, _ := command.Cmd("cat /etc/ssh/sshd_config | grep 'Port '", false)
	if err != nil {
		return 22
	}
	str := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(output, "")
	port, _ := strconv.Atoi(str)
	return port
}

func GetLocalIP() string {
	err, output, _ := command.Cmd("ip -o -4  address show | awk ' NR==2 { gsub(/\\/.*/, \"\", $4); print $4 } '", false)
	if err != nil {
		return "0.0.0.0"
	}
	clean := regexp.MustCompile(`[^0-9.]+`).ReplaceAllString(output, "")
	clean = strings.Replace(clean, " ", "", -1)
	return clean
}

func GetPublicIP() string {
	err, output, _ := command.Cmd("(curl -s ipinfo.io/ip | grep -o -Eo '[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}.[0-9]{1,3}')", false)
	if err != nil {
		return "0.0.0.0"
	}
	clean := regexp.MustCompile(`[^0-9.]+`).ReplaceAllString(output, "")
	return clean

}

func GetVersion() int {
	// err, output, _ := command.Cmd("git -C /opt/anyshell rev-list --count main", false);
	err, output, _ := command.Cmd("curl -s -I -k \"https://api.github.com/repos/alexmnr/anyshell/commits?per_page=1\" | sed -n '/^[Ll]ink:/ s/.*\"next\".*page=\\([0-9]*\\).*\"last\".*/\\1/p'", false)
	out.Info(output)
	if err != nil {
		return 0
	}
	clean := regexp.MustCompile(`[^0-9]+`).ReplaceAllString(output, "")
	clean = strings.Replace(clean, " ", "", -1)
	version, _ := strconv.Atoi(clean)
	return version
}
