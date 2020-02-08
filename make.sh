#!/bin/bash

set -xe

CC=/opt/cross-pi-gcc/bin/arm-linux-gnueabihf-gcc CGO_ENABLED=1 GOARCH=arm GOARM=6 go build -o cnc_gamepad-raspi
scp -i raspi.key cnc_gamepad-raspi pi@10.1.0.123:cnc/cnc_gamepad

