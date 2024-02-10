#!/bin/bash
# check if dependencies are met
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
version="1.21.1"
dep=true
missing=""
if ! command -v git &> /dev/null; then
  dep=false
  missing=$(echo "git $missing")
fi
if ! command -v wget &> /dev/null; then
  dep=false
  missing=$(echo "wget $missing")
fi
if ! command -v curl &> /dev/null; then
  dep=false
  missing=$(echo "curl $missing")
fi
if ! command -v sudo &> /dev/null; then
  dep=false
  missing=$(echo "sudo $missing")
fi

if [ "$dep" = "false" ]; then
  echo "INFO: Not all dependencies were met, please install following packages: $missing"
  exit
fi

# check if go is installed
if ! command -v go &> /dev/null; then
  go=false
  echo "INFO: go not found, trying to install..."
else
  go=true
fi

# install go if necessary
if [ "$go" = "false" ]; then
  # get info about cpu
  if [ ! -z "$(lscpu | grep 'aarch64')" ]; then
    arc="aarch64"
    #arc="armv7l"
    echo "INFO: Detected aarch64 architecture"
  elif [ ! -z "$(lscpu | grep 'armv6l')" ]; then
    arc="armv7l"
    echo "INFO: Detected armv6l architecture"
  elif [ ! -z "$(lscpu | grep 'armv7l')" ]; then
    arc="armv7l"
    echo "INFO: Detected armv7l architecture"
  elif [ ! -z "$(lscpu | grep 'x86_64')" ]; then
    arc="x86_64"
    echo "INFO: Detected x86_64 architecture"
  else
    echo "ERROR: architecture not detected"
    exit 1
  fi

  # install go from source
  if [ "$arc" = "x86_64" ]; then
    sudo rm -rf /usr/local/go /usr/bin/go /usr/bin/gofmt &> /dev/null
    cd /tmp
    link="https://go.dev/dl/go$version.linux-amd64.tar.gz"
    wget $link
    sudo tar -C /usr/local -xzf go$version.linux-amd64.tar.gz
    sudo rm -f /usr/bin/go
    sudo rm -f /usr/bin/gofmt
    sudo ln -s /usr/local/go/bin/go /usr/bin
    sudo ln -s /usr/local/go/bin/gofmt /usr/bin
  elif [ "$arc" = "armv7l" ]; then
    sudo rm -rf /usr/local/go /usr/bin/go /usr/bin/gofmt &> /dev/null
    cd /tmp
    link="https://go.dev/dl/go$version.linux-armv6l.tar.gz"
    wget $link
    sudo tar -C /usr/local -xzf go$version.linux-armv6l.tar.gz
    sudo rm -f /usr/bin/go
    sudo rm -f /usr/bin/gofmt
    sudo ln -s /usr/local/go/bin/go /usr/bin
    sudo ln -s /usr/local/go/bin/gofmt /usr/bin
  elif [ "$arc" = "aarch64" ]; then
    sudo rm -rf /usr/local/go /usr/bin/go /usr/bin/gofmt &> /dev/null
    cd /tmp
    link="https://go.dev/dl/go$version.linux-arm64.tar.gz"
    wget $link
    sudo tar -C /usr/local -xzf go$version.linux-arm64.tar.gz
    sudo rm -f /usr/bin/go
    sudo rm -f /usr/bin/gofmt
    sudo ln -s /usr/local/go/bin/go /usr/bin
    sudo ln -s /usr/local/go/bin/gofmt /usr/bin
  else
    echo "ERROR: Can't automatically install go on your system, you need to do it manually"
    exit 1
  fi

fi

echo "INFO: All dependencies were met!"

# move current dir to /opt/anyshell
echo "INFO: Installing anyshell in /opt/anyshell..."
if [ "$SCRIPT_DIR" = "/opt/anyshell" ]; then
  sudo chown $USER:$USER /opt/anyshell -R
  cd /opt/anyshell
else
  sudo rm -rf /opt/anyshell &>/dev/null
  sudo cp -r $SCRIPT_DIR /opt/anyshell
  sudo chown $USER:$USER /opt/anyshell -R
  cd /opt/anyshell
fi

echo "INFO: Building anyshell..."

cd /opt/anyshell/go
go build .

sudo rm -f /usr/bin/anyshell &>/dev/null
sudo rm -f /usr/bin/any &>/dev/null
sudo ln -s /opt/anyshell/go/anyshell /usr/bin/anyshell
sudo ln -s /opt/anyshell/go/anyshell /usr/bin/any

echo "INFO: Linking anyshell service..."
sudo rm -f /etc/systemd/system/anyshell.service &>/dev/null
sudo ln -s /opt/anyshell/etc/anyshell.service /etc/systemd/system

echo "INFO: done!"
