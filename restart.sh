#!/bin/bash

MYSQL_HOST="localhost"
MYSQL_USERNAME="chill"
MYSQL_PASSWORD=""
CHILL_USERNAME=""
CHILL_PASSWORD=""

# pidfile="chill.pid"
# if [ -f $pidfile ]; then
#     pid=$(cat $pidfile)
#     kill -15 $pid
# fi
# > $pidfile

cd "/home/alphaboi/services/chill/build"

MYSQL_HOST="$MYSQL_HOST" \
    MYSQL_USERNAME="$MYSQL_USERNAME"  \
    MYSQL_PASSWORD="$MYSQL_PASSWORD"  \
    CHILL_USERNAME="$CHILL_USERNAME"  \
    CHILL_PASSWORD="$CHILL_PASSWORD"  \
    ./app_chill &

pid="$!"
echo "Chill running with pid $pid"
# echo "$pid" > "$pidfile"
