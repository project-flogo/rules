This example demonstrates the capability of rules to track and trace for a flogo app. In this example three tuples used, tuples description is shown below.

```json
{
    "name": "package",
    "properties": [
        {
            "name": "id",
            "type": "string",
            "pk-index": 0
        },
        {
            "name": "state",
            "type": "string"
        }
    ]  
},
{
    "name": "moveevent",
    "ttl": 0,
    "properties": [
        {
            "name": "id",
            "type": "string",
            "pk-index": 0
        },
        {
            "name": "packageid",
            "type": "string"
        },
        {
            "name": "sitting",
            "type": "double"
        },
        {
            "name": "moving",
            "type": "double"
        },
        {
            "name": "dropped",
            "type": "double"
        }
    ]
},
{
    "name": "movetimeoutevent",
    "ttl": 0,
    "properties": [
        {
            "name": "id",
            "type": "string",
            "pk-index": 0
        },
        {
            "name": "packageid",
            "type": "string"
        },
        {
            "name": "timeoutinmillis",
            "type": "integer"
        }
    ]  
}
```

`package` tuple is always stored in network, while the other tuples `moveevent` and `movetimeoutevent` are removed immediate after its usage as `ttl` given `0`. During startup `PACKAGE1` is asserted into network.

### Actions used here

`aJoinMoveEventAndPackage` : Performs check on sitting value. If value is more than 0.5 then package is scheduled to movetimeoutevent only once.<br>
`aPrintMoveEvent`: Prints received moveevents.<br>
`aMoveTimeoutEvent`: Prints received movetimeoutevent.<br>
`aJoinMoveTimeoutEventAndPackage`: All packages with state as normal are modified to sitting state.<br>
`aPackageInSitting`: Prints package as sitting.

## Usage
Get the repo and in this example `main.go`, `functions.go` both are available we can directly build and run the app or create flogo rule app and run it.

### Direct build and run
```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/trackntrace
go build
./trackntrace
```

### Create app using flogo cli
```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/trackntrace
flogo create -f flogo.json trackNTraceApp
cp functions.go trackNTraceApp/src
cd trackNTraceApp
flogo build
./bin/trackNTraceApp
```

## Move event test
Run below command to check moveevent action on PACKAGE1.
```sh
curl http://localhost:7777/test/moveevent?id=PACKAGE1\&packageid=pkgid1\&sitting=0.4
```
Above command results in executing two actions one is `aJoinMoveEventAndPackage` and the other is `aPrintMoveEvent`.<br><br>
Expected output:
```
Joining a 'moveevent' with packageid [pkgid1] to package [PACKAGE1], sitting [0.400000], moving [0.000000], dropped [0.000000]
Received a 'moveevent' [PACKAGE1] sitting [0.400000], moving [0.000000], dropped [0.000000]
```

If sitting value is given more than 0.5 then movetimeoutevent gets invoked.
```sh
curl http://localhost:7777/test/moveevent?id=PACKAGE1\&packageid=pkgid1\&sitting=0.6
```
Above commands results in executing action `aJoinMoveTimeoutEventAndPackage` from `aJoinMoveEventAndPackage`. Action `aJoinMoveTimeoutEventAndPackage` results in changing state of PACKAGE1 to sitting this will trigger another action `aPackageInSitting`. So chain of actions getting executed.<br><br>
Expected output:
```
Joining a 'moveevent' with packageid [pkgid1] to package [PACKAGE1], sitting [0.600000], moving [0.000000], dropped [0.000000]
Starting a 15s timer.. [PACKAGE1]
Received a 'moveevent' [PACKAGE1] sitting [0.600000], moving [0.000000], dropped [0.000000]
Joining a 'movetimeoutevent' [PACKAGE1] to package [PACKAGE1], timeout [15000]
PACKAGE [PACKAGE1] is STTTING
'movetimeoutevent' id [01D9YMAAT45PABMKR4K297S2G5], packageid [PACKAGE1], timeoutinmillis [15000]
```

## Package test
Execute below command to store `package2`
```sh
curl http://localhost:7777/test/package?id=package2\&state=normal
```
Expected output:
```
Received package [package2]
```
Store `package3` as in sitting position.
```sh
curl http://localhost:7777/test/package?id=package3\&state=sitting
```
Expected output:
```
PACKAGE [package3] is STTTING
Received package [package3]
```

## Movetimeoutevent test
Execute below command to see `package2` is getting changed to sitting.
```sh
curl http://localhost:7777/test/movetimeoutevent?id=package2\&packageid=pkgid1
```
Expected result:
```
Joining a 'movetimeoutevent' [pkgid1] to package [package3], timeout [0]
Joining a 'movetimeoutevent' [pkgid1] to package [package2], timeout [0]
PACKAGE [package2] is STTTING
'movetimeoutevent' id [package2], packageid [pkgid1], timeoutinmillis [0]
```