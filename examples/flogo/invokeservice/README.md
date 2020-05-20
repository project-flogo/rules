# Invoke Rule Service

This example demonstrates how a rule can invoke a rule `service`. A rule `service` is a `go-function` or a `flogo-activity`.

## Setup and build
Once you have the `flogo.json` file and a `functions.go` file, you are ready to build your Flogo App

### Steps

```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/invokeservice
flogo create -f flogo.json invokeservice
cp functions.go ./invokeservice/src
cd invokeservice
flogo build
cd bin
./invokeservice
```
### Testing

#### #1 Invoke go-function based service

Send a curl request
`curl http://localhost:7777/test/n1?name=function`
You should see following output:
```
Rule[n1.name == function] fired. serviceFunctionAction() function got invoked.
```

#### #2 Invoke flogo-activity based service
Send a curl request
`curl http://localhost:7777/test/n1?name=activity`
You should see following output:
```
2019-08-19T20:11:11.068+0530	INFO	[flogo.test] -	activity
```