# Invoke service when rule fires

This example demonstrates how a rule can invoke rule `service`. A rule `service` is a `go function` or a `flogo action` or a `flogo activity`.

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

#### #1 Invoke function based service

Send a curl request
`curl http://localhost:7777/test/n1?name=function`
You should see following output:
```
Rule[n1.name == function] fired. serviceFunctionAction() function got invoked.
```

#### #2 Invoke flogo async action (Ex: flow action) based service

Send a curl request
`curl http://localhost:7777/test/n1?name=asyncaction`
You should see following output:
```
2019-08-19T20:10:07.267+0530	INFO	[flogo.activity.log] -	asyncaction
service[ServiceFlowAction] outputs: map[] 
2019-08-19T20:10:07.267+0530	INFO	[flogo.flow] -	Instance [e7f8a513674f52ce0104fc4e8e8392d5] Done
service[ServiceFlowAction] executed successfully asynchronously
```

#### #3 Invoke flogo sync action based service
Send a curl request
`curl http://localhost:7777/test/n1?name=syncaction`
You should see following output:
```
2019-08-19T20:10:43.173+0530	INFO	[flogo] -	Input: syncaction
service[ServiceCoreAction] executed successfully. Service outputs: map[anOutput:syncaction] 
```

#### #4 Invoke activity based service
Send a curl request
`curl http://localhost:7777/test/n1?name=activity`
You should see following output:
```
2019-08-19T20:11:11.068+0530	INFO	[flogo.test] -	activity
```