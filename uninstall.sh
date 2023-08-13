#!/bin/bash
# remove /opt/anyshell
echo "INFO: Removing anyshell from /opt/anyshell..."
sudo rm -rf /opt/anyshell &>/dev/null

echo "INFO: Removing anyshell binary..."
sudo rm -f /usr/bin/anyshell

echo "INFO: Removing anyshell service..."
sudo rm -f /etc/systemd/system/anyshell.service

echo "INFO: Linking anyshell config..."
sudo rm -f /etc/anyshell

echo "INFO: done!"
