#!/usr/bin/env bash
set -euo pipefail
export CUST_DSN="./demo.sqlite"
export CUST_USER="james"

echo "dsn=$CUST_DSN, user=$CUST_USER"

#tear down the server after we are done
trap "kill 0" exit

./custody serve &
./custody create

function sign() {
    msg="$1"
    echo "$msg" | ./custody sign 
}

function list() {
    ./custody list --username "$CUST_USER"
}

sign "Hello World"
sign "upload screenshot.png"
sign "enhance screenshot.png"
sign "run 'facedetections' on screenshot.png"
sign "print screenshot.png"
sign "submit screenshot.png to court"

list

echo "Returning server to foreground Ctrl-C to stop"
wait