#!/usr/bin/env bash
#dsn="--dsn ./demo.sqlite"
#PATH="$PATH:."
export DSN="./demo.sqlite"
export CUSER="james"


custody create --dsn "$DSN" --username "$CUSER"

function sign() {
    msg="$1"
    echo "msg" | ./custody sign --dsn "$DSN" --username "$CUSER"
}

function list() {
    ./custody list --dsn "$DSN" --username "$CUSER"
}

sign "Hello World"
sign "upload screenshot.png"
sign "enhance screenshot.png"
sign "run 'facedetections' on screenshot.png"
sign "print screenshot.png"
sign "submit screenshot.png to court"

list
