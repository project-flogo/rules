
## Installation

### Prerequisites
* The Go programming language version 1.8 or later should be [installed](https://golang.org/doc/install).
* The **GOPATH** environment variable on your system must be set properly
* Docker 

## Setup and Usage

The following conditions are used in the example:

* `checkForBob`: Checks if the `n1` tuple name is Bob
* `checkSameNamesCondition`: Checks if the name in `n1` tuple is same as name in `n2` tuple
* `checkSameEnvName`: Checks if the name in `n1` tuple matches the value stored in environment variable `name`

The following rules are used in the example:

* `checkForBobAction`: Gets fired when `checkForBob` is true
* `checkSameNamesAction`: Gets fired when `checkForBob` and `checkSameNamesCondition` conditions are true
* `checkSameEnvNameAction`: Gets fired when `checkSameEnvName` condition is true

Run the example:

```
docker run -p 6383:6379 -d redis
go get -u github.com/project-flogo/rules/...
cd $GOPATH/src/github.com/project-flogo/rules/examples/rulesapp
export name=Smith
go run main.go
```