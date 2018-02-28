#!/usr/bin/env bash
set -euo pipefail
export CUST_DSN="./demo.sqlite"
export CUST_USER="james"

echo "dsn=$CUST_DSN, user=$CUST_USER"

# start server and track pid to send shutdown signal
./custody serve &
SRVPID=$!


./custody create

function sign() {
    msg="$1"
    echo "$msg" | ./custody sign 
}

function list() {
    ./custody list --username "$CUST_USER"
}

# sign some messages
sign "Hello World"
sign "upload screenshot.png"
sign "enhance screenshot.png"
sign "run 'facedetections' on screenshot.png"
sign "print screenshot.png"
sign "submit screenshot.png to court"

# list the messages we just signed
list

#send shutdown signal
echo "Shutting down server..."
kill ${SRVPID}
wait
echo "Server stopped"

echo "To clean up the database remove $CUST_DSN"
exit 0