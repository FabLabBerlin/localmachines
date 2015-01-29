#!/bin/bash
#
# 1. Copy this file 
#    cp launch-example.sh ~/fabsmith.sh
# 
# 2. Replace the values with yours and save
#
# 3. Make it launchable
#    chmod a+x ~/fabsmith.sh
#
# 4. Use it!

sudo FABSMITH_RUNMODE="prod" \
FABSMITH_HTTP_PORT="80" \
FABSMITH_MYSQL_USER="youruser" \
FABSMITH_MYSQL_PASS="yourpass" \
FABSMITH_MYSQL_DB="fabsmith" \
../fabsmith
