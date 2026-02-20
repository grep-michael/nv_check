#!/bin/bash

if [ $# -ne 1 ]; then
  echo "Usage: $0 [driver_version]"
  exit 1
fi

DRIVER=$1
VER=$(uname -r)
sudo apt install -y linux-modules-nvidia-$DRIVER-$VER \
	linux-objects-nvidia-$DRIVER-$VER \
	nvidia-utils-$DRIVER
sudo apt install -y nvidia-cuda-toolkit nvidia-cuda-toolkit-gcc
sudo apt install -y build-essential gcc
