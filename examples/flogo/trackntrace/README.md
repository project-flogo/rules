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
            "name": "changeStateTo",
            "type": "string"
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

Consider system having incoming packages. In order to move packages from source to destination we have steps like sitting,moving,dropped and delayed. When the first event sitting comes, package is scheduled to 10s timeout. With in 10s if we receive any other event like moving/dropped then scheduled timeout event will get canceled. Otherwise package is marked as delayed and removed from cluster.

### Actions used here

`aJoinMoveEventAndPackage` : If changeStateTo value is more sitting then package is scheduled to movetimeoutevent by 10 seconds.<br>
`aPrintMoveEvent`: Prints received moveevents.<br>
`aMoveTimeoutEvent`: Prints received movetimeoutevent.<br>
`aJoinMoveTimeoutEventAndPackage`: All packages with state as sitting are modified to delayed state.<br>
`aPackageInSitting`: Prints package as sitting.<br>
`aPackageInMoving`: Prints package as moving.<br>
`aPackageInDropped`: Prints package as dropped.<br>
`aPackageInDelayed`: Prints package as delayed.

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
curl http://localhost:7777/moveevent?packageid=PACKAGE1\&changeStateTo=sitting
```

Above commands results in executing action `aJoinMoveTimeoutEventAndPackage` from `aJoinMoveEventAndPackage`. Action `aJoinMoveTimeoutEventAndPackage` results in changing state of PACKAGE1 to sitting this will trigger another action `aPackageInSitting`. So chain of actions getting executed.<br><br>
Expected output:
```
Received a 'moveevent' [01DC6XBXHMBGH043ZJQVSTS60A] change state to [sitting]
Joining a 'moveevent' with packageid [PACKAGE1] to package [PACKAGE1], change state to [sitting]
Starting a 10s timer.. [PACKAGE1]
PACKAGE [PACKAGE1] is Sitting
Received package [PACKAGE1]
Joining a 'moveevent' with packageid [PACKAGE1] to package [PACKAGE1], change state to [sitting]
Received a 'movetimeoutevent' id [01DC6XBXHMY4E5JNEFZMPS1PZG], packageid [PACKAGE1], timeoutinmillis [10000]
Joining a 'movetimeoutevent' [PACKAGE1] to package [PACKAGE1], timeout [10000]
PACKAGE [PACKAGE1] is Delayed
```
Above we can see `PACKAGE1` is went into delayed as no operation is done in scheduled 10s interval. Restart the rules app and run below command.

```sh
curl http://localhost:7777/moveevent?packageid=PACKAGE1\&changeStateTo=sitting
curl http://localhost:7777/moveevent?packageid=PACKAGE1\&changeStateTo=moving
```
Expected output:
```
Received a 'moveevent' [01DC6YYQBCHPP0VDQNA99ZDD86] change state to [sitting]
Joining a 'moveevent' with packageid [PACKAGE1] to package [PACKAGE1], change state to [sitting]
Starting a 10s timer.. [PACKAGE1]
PACKAGE [PACKAGE1] is Sitting
Received package [PACKAGE1]
Joining a 'moveevent' with packageid [PACKAGE1] to package [PACKAGE1], change state to [sitting]
Received a 'moveevent' [01DC6YYV4VZGHWWBS0YKZCXPVB] change state to [moving]
Joining a 'moveevent' with packageid [PACKAGE1] to package [PACKAGE1], change state to [moving]
Cancelling timer attached to key [PACKAGE1]
PACKAGE [PACKAGE1] is Moving
Received package [PACKAGE1]
Joining a 'moveevent' with packageid [PACKAGE1] to package [PACKAGE1], change state to [moving]
```
## Package test
Execute below command to store `package2`
```sh
curl http://localhost:7777/package?id=PACKAGE2\&state=normal
```
Expected output:
```
Received package [package2]
```