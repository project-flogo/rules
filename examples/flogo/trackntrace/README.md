This recipe demonstrates the capability of rules to track and trace for a flogo app. Here three tuples used, tuples description is shown below.

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

`package` tuple is always stored in network, while the others `moveevent` and `movetimeoutevent` are removed immediate after its usage as `ttl` given as 0. By default during startup `PACKAGE1` is stored. 

## Move event test
Perform below mentioned steps to check the result.
```sh
curl http://localhost:7777/test/moveevent?id=PACKAGE1\&packageid=pkgid1\&sitting=0.4
```
Above command results in executing two actions one is `aJoinMoveEventAndPackage` and the other is `aPrintMoveEvent`.<br>
Expected output:
```
Joining a 'moveevent' with packageid [pkgid1] to package [PACKAGE1], sitting [0.400000], moving [0.000000], dropped [0.000000]
Received a 'moveevent' [PACKAGE1] sitting [0.400000], moving [0.000000], dropped [0.000000]
```

If sitting value is given more than 0.5 then package is moved with timeout.
```sh
curl http://localhost:7777/test/moveevent?id=PACKAGE1\&packageid=pkgid1\&sitting=0.6
```
Above commands results in executing action `aJoinMoveTimeoutEventAndPackage` from action `aJoinMoveEventAndPackage`. The action `aJoinMoveTimeoutEventAndPackage` results in changing state of PACKAGE1 to sitting this will trigger another action `aPackageInSitting`. So chain of actions getting executed.<br>
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