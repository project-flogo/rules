# Invoke action when rule fires

This rules example demonstrates usage of flogo rule functions, flogo activity, flogo flow action and core action based services.

## Setup and build
Once you have the `flogo.json` file and a `functions.go` file, you are ready to build your Flogo App

### Pre-requisites
* Go 1.11
* Download and build the Flogo CLI 'flogo' and add it to your system PATH

### Steps

```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/invokeservice
flogo create -f flogo.json
cp functions.go ./invokeservice/src
cd invokeservice
flogo build
cd bin
./invokeservice
```
### Testing

#### Invoke activity based service
Send a curl request
`curl http://localhost:7777/test/n1?name=Tom`
You should see following output:
```
2019-08-12T16:06:03.828+0530	INFO	[flogo.test] -	Tom
```

#### Invoke function based service

Send a curl request
`curl http://localhost:7777/test/n1?name=Bob`
You should see following output:
```
Rule[n1.name == Bob] fired. checkForBobAction() function got invoked.
```


#### Invoke flogo core action based service
Send a curl request
`curl http://localhost:7777/test/n1?name=Michael`
You should see following output:
```
service[CoreActionService] executed successfully. Service outputs: map[anOutput:Michael]
```

#### Invoke flogo flow action  based service

Send a curl request
`curl http://localhost:7777/test/n1?name=Robert`
You should see following output:
```
service[FlowActionService] executed successfully asynchronously.
```