# Custodyctl

Have you ever wanted to track chain of custody for digital forensic evidence?
Custodyctl is the solution you have been looking for.
You can create identities, sign records, and display ledgers. 

Chain of custody has two purposes in police work, first to ensure that at every moment in time,
there is a person responsible for the integrity of the evidence, and also to enable "fruit of the
poison tree" analysis to identify any evidence that is inadmissible due to a procedural flaw.

Custodyctl is a service that enables the construction of various custody tracking applications.

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

### Prerequisites

Make sure you have a working installation of golang [https://golang.org/doc/install](official docs).
And a go workspace [https://golang.org/doc/code.html](official docs).
If you want to change the database schema you will need [github.com/knq/xo](xo).


### Installing

To install the package, make sure you have a working golang environment established.
```bash
go get github.gatech.edu/NIJ-Grant/custody
go install github.gatech.edu/NIJ-Grant/custody 
```


You can run a little demo in ./demo.sh.
Where custody is on your $PATH.
```bash
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
```

You can get machine readable output with the `--json flag`.

## Running the tests

The tests are developed using go tests. You can run `make test` or `go test ./...`


### And coding style tests

This project uses `go fmt` to enforce a uniform coding style.

```
go fmt ./...
```

## Deployment

Because this application is written in go using sqlite for the database, deployment is trivial.

1. Get cross compiled binaries using`make all`
2. SCP them to your host
3. run `export DSN="path/todatabase.sqlite"; ./custody serve`

## Built With

* mattn/sqlite3
* knq/xo
* spf13/cobra
* spf13/viper


## Contributing

Please read [CONTRIBUTING.md](https://github.gatech.edu/NIJ-Grant/custody/Contributing.md) for details on the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.gatech.edu/NIJ-Grant/custody/tags). 

## Authors

* James Fairbanks <james.fairbanks@gtri.gatech.edu>

See also the list of [contributors](https://github.gatech.edu/NIJ-Grant/custody/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Future Work

- Event structure, the current application uses only identities and messages, but you could attach
  an operation and object to every operation.
- Dependency Graph, how can we store the dependencies between files.
- Grouping files by cases, how do you group files into cases.
- Integration with AffLib which is designed to store and track whole disk images


## Acknowledgments

* NIJ Grant Number: XXXXXXXX
* Dekalb County Police Department
* Blockchain
