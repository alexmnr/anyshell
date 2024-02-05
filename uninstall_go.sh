#!/bin/bash
echo "INFO: Removing binary links..."
sudo rm -rf /usr/bin/go &>/dev/null
sudo rm -rf /usr/bin/gofmt &>/dev/null

echo "INFO: Removing go folder..."
sudo rm -rf /usr/local/go/

echo "INFO: done!"
