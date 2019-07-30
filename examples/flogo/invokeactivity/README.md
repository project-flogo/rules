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

Then from another command line, send a curl request
`curl localhost:7777/test/n1?name=Bob`
You should see this o/p on the console

