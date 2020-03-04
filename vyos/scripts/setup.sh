#!/bin/sh
# Setup auto deploy on VyOS.
cp /config/scripts/network/vyos/scripts/vyos-postconfig-bootup.script /config/scripts/vyos-postconfig-bootup.script
/config/scripts/vyos-postconfig-bootup.script

# Make sure the task-scheduler is set.
sg vyattacfg -c /config/scripts/network/vyos/scripts/set-task.sh
