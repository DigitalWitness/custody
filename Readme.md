# Custodyctl

Have you ever wanted to track chain of custody for digital forensic evidence?
Custodyctl is the solution you have been looking for.
You can create identities, sign records, and display ledgers. 

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
export CUST_DSN="./demo.sqlite"
export CUST_USER="$USER"

echo "dsn=$CUST_DSN, user=$CUST_USER"

custody create

function sign() {
    msg="$1"
    echo "$msg" | custody sign
}

function list() {
    custody list --username "$CUST_USER"
}

sign "Hello World"
sign "upload screenshot.png"
sign "enhance screenshot.png"
sign "run 'facedetections' on screenshot.png"
sign "print screenshot.png"
sign "submit screenshot.png to court"

list
```

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

Please read [CONTRIBUTING.md](https://gist.github.com/PurpleBooth/b24679402957c63ec426) for details on the process for submitting pull requests to us.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.gatech.edu/NIJ-Grant/custody/tags). 

## Authors

* James Fairbanks <james.fairbanks@gtri.gatech.edu>

See also the list of [contributors](https://github.com/your/project/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* NIJ Grant Number: XXXXXXXX
* Dekalb County Police Department
* Blockchain
