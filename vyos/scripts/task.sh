#!/bin/sh
# Pull new config from remote git repo.
git --git-dir=/config/scripts/network/.git --work-tree=/config/scripts/network pull origin master

# TODO: GENERATE AND LOAD CONFIG HERE.
