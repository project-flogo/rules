## Flogo Rules based Track and Trace

This example demonstrates the capability of rules to track and trace. In this example three tuples are used, tuples description is given below.

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
            "name": "targetstate",
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

`package` tuple is always stored in network, while the other tuples `moveevent` and `movetimeoutevent` are removed after usage as `ttl` is given as `0`. During startup `PACKAGE1` is asserted into network.

Consider system having incoming packages. In order to move packages from source to destination we have steps like sitting,moving,dropped and delayed. When the first event sitting comes, package is scheduled to 10s timeout. Within 10s if we receive  moving event then scheduled timeout event will get canceled. Otherwise package is marked as delayed and retracted from network.

### Package State Info

<p align="center">
  <img src ="./trackntrace.png" />
</p>

In detail, above image represents state change for a given package. Consider insert package event this will insert a package into network with state as `normal`. This package now accepts only `sitting` event. When `sitting` event is triggered, a 10 seconds timer is created to trigger `delayed` event, within 10s only `moving` event can cancel the timer and state is changed to `moving` otherwise `delayed` event gets triggered. If `dropped` event occurs on a package with `moving` state, then package state is changed to `dropped`. Any package with state as `dropped` or `delayed` is retracted from network. You have to insert the package again to use.

### Actions used here

`aJoinMoveEventAndPackage` : If targetstate value sitting then package is scheduled to movetimeoutevent by 10 seconds.<br>
`aPrintMoveEvent`: Prints received moveevents and store packages into network based on targetstate.<br>
`aMoveTimeoutEvent`: Prints received movetimeoutevent.<br>
`aJoinMoveTimeoutEventAndPackage`: Package is modified to moveevent target state.<br>
`aPackageInSitting`: Prints package as sitting.<br>
`aPackageInMoving`: Prints package as moving.<br>
`aPackageInDropped`: Prints package as dropped and retracts package tuple from network.<br>
`aPackageInDelayed`: Prints package as delayed and retracts package tuple from network.

## Run the example

The example contains `main.go` and `rulesapp.json`. To run the application:
```
go run main.go
``` 

In the given example 2 packages PACKAGE1 and PACKAGE2 are inserted.
The PACKAGE1 is asserted into the Rulesession with `normal` state. The package is further moved to `sitting`, `moving` and `dropped` states.When the package reaches the dropped state the package is removed from the rulesession.
The PACKAGE2 is asserted into the Rulesession with `normal` state. The package is then moved to sitting state. After 10 seconds a movetimeout event is asserted and the package state is changed to `delayed`.

Below is the final output for package PACKAGE2:

```
Received a 'moveevent' [PACKAGE2] target state [dropped]
Joining a 'moveevent' with packageid [PACKAGE2] to package [PACKAGE2], target state [dropped]
PACKAGE [PACKAGE2] is Dropped
Saving tuple. Type [package] Key [package:id:PACKAGE2], Val [&{package map[id:PACKAGE2 state:dropped] 0xc0000bad20 0xc0000c1720}]
```


Below is the final output for package PACKAGE1:

```
Received a 'movetimeoutevent' id [01E4JW2J1V9AEJJ1KTTGXAAW54], packageid [PACKAGE1], timeoutinmillis [10000]
Joining a 'movetimeoutevent' [PACKAGE1] to package [PACKAGE1], timeout [10000]
PACKAGE [PACKAGE1] is Delayed
Saving tuple. Type [package] Key [package:id:PACKAGE1], Val [&{package map[id:PACKAGE1 state:delayed] 0xc0000baa20 0xc0000c0c80}]
```