#!/bin/sh
# Pull new config from remote git repo.
git --git-dir=/config/scripts/network/.git --work-tree=/config/scripts/network pull origin master

# Replace bootup config with current one.
cp /config/scripts/network/vyos/scripts/vyos-postconfig-bootup.script /config/scripts/vyos-postconfig-bootup.script

# TODO: GENERATE AND LOAD CONFIG HERE.
