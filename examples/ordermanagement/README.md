# FlogoRules based Order Management System

The example demonstrates the use of [FlogoRules]("https://github.com/project-flogo/rules") library to create a light weight rules based Order Management System with an ingest from a source. Here we are using [MQTT]("http://mqtt.org/") as the ingest source. MQTT is a very commonly used lightweight pub/sub protocol for edge devices. But pretty much any ingest source(via flogo/custom) can be used, i.e. rest/kafka/kinesis/lambda,etc. For this example we are using Flogo MQTT Trigger.
Additionally, an Audit Trail of the activities (events/rule executions) is published to a stream. Here we are publishing to AWS Kinesis Stream. Stream is created if one does not exists, provided valid IAM exists for the given user credentials. And finally, reading off the Kinesis stream and pushing it out to a web client as a continous stream.

## Deployment
Below is one sample illustration of an deployment approach when setting things up on AWS. For now to keep things simple you can set it up on your local box.

<p align="center">
  <img src ="https://raw.githubusercontent.com/project-flogo/rules/master/examples/ordermanagement/web/resources/awsDeployment.png" />
</p>

## Installation

### Prerequisites
To get started with the Flogo Rules you'll need to have a few things
* The Go programming language version 1.8 or later should be [installed](https://golang.org/doc/install).
* The **GOPATH** environment variable on your system must be set properly

### Try this out
Make sure you have valid AWS credentials with appropriate IAM for interacting with Kinesis streams. For MQTT, we are using the publicly available test broker. You can always setup your own broker and refer to that.

```
$ go get github.com/project-flogo/rules/examples/ordermanagement
$ cd $GOPATH/src/github.com/project-flogo/rules/examples/ordermanagement
$ go generate
$ go build
$ ./ordermanagement -h
```

Finally load '$GOPATH/src/github.com/project-flogo/rules/examples/ordermanagement/web/oms.html' in your preferred browser and send events to the server to trigger various rules. As the engine processes different events, a stream of Audit Trail items should be generated and loaded in the right pane as a continous stream.

