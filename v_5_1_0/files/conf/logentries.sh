#!/bin/sh

LOGENTRIES_PREFIX=$1
LOGENTRIES_TOKEN=$2

journalctl -o json -f | jq -r --arg LOGENTRIES_PREFIX "${LOGENTRIES_PREFIX}" '$LOGENTRIES_PREFIX + " " + .__REALTIME_TIMESTAMP + " " + ._HOSTNAME + " " + ._SYSTEMD_UNIT + " " + .MESSAGE' | awk -v LOGENTRIES_TOKEN=${LOGENTRIES_TOKEN} '{ print LOGENTRIES_TOKEN, $0; fflush(); }' | ncat data.logentries.com 10000
