#!/bin/vbash
source /opt/vyatta/etc/functions/script-template
configure
delete system task-scheduler task pull-latest-config
set system task-scheduler task pull-latest-config interval 3m
set system task-scheduler task pull-latest-config executable path /config/scripts/network/vyos/scripts/task.sh
commit
save
exit
