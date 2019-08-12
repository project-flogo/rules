# Invoke activity when rule fires

## Setup and build
Once you have the `flogo.json` file and a `functions.go` file, you are ready to build your Flogo App

###Pre-requisites
* Go 1.11
* Download and build the Flogo CLI 'flogo' and add it to your system PATH

### Steps

```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/invokeactivity
flogo create -f flogo.json
cp functions.go ./invokeactivity/src
cd invokeactivity
flogo build
cd bin
./invokeactivity
```
### Testing

#### Invoke activity based service
Snd a curl request
`curl localhost:7777/test/n1?name=Tom`
You should see following output:
```sh
2019-08-12T16:06:03.828+0530	INFO	[flogo.test] -	Tom
```

#### Invoke function based service

Snd a curl request
`curl localhost:7777/test/n1?name=Bob`
You should see following output:
```sh
Rule[n1.name == Bob] fired. checkForBobAction() function got invoked.
```