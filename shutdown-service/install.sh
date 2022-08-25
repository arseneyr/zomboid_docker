#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

install -g root -o root -m 0544 "$SCRIPT_DIR/shutdownfifo-listener" /usr/local/libexec
install -g root -o root -m 0444 "$SCRIPT_DIR/shutdownfifo.service" "$SCRIPT_DIR/shutdownfifo.socket" /etc/systemd/system

systemctl daemon-reload