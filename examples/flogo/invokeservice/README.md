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
2019-08-20T09:57:46.780+0530	INFO	[flogo.activity.log] -	asyncaction
2019-08-20T09:57:46.781+0530	INFO	[flogo.rules] -	service[ServiceFlowAction] outputs: map[] 

2019-08-20T09:57:46.781+0530	INFO	[flogo.flow] -	Instance [39470b3be53593aa827043a05086504f] Done
2019-08-20T09:57:46.781+0530	INFO	[flogo.rules] -	service[ServiceFlowAction] executed successfully asynchronously
```

#### #3 Invoke flogo sync action based service
Send a curl request
`curl http://localhost:7777/test/n1?name=syncaction`
You should see following output:
```
2019-08-20T09:58:21.090+0530	INFO	[flogo] -	Input: syncaction
2019-08-20T09:58:21.090+0530	INFO	[flogo.rules] -	service[ServiceCoreAction] executed successfully. Service outputs: map[anOutput:syncaction] 
```

#### #4 Invoke activity based service
Send a curl request
`curl http://localhost:7777/test/n1?name=activity`
You should see following output:
```
2019-08-19T20:11:11.068+0530	INFO	[flogo.test] -	activity
```