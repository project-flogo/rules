# Kafka trigger based simple rules action usage

This example demonstrates usage of kafka trigger with rules action. Consider system related to order identification and processing, kafka trigger will listen to a topic where order information is available, based on price and order category gift coupons can be issued to the customer.

## Setup
1. Get rules repo.
```sh
go get -u github.com/project-flogo/rules/... 
```
2. Install [kafka](https://kafka.apache.org/quickstart).

## Usage
Get the repo and in this example `main.go`, `functions.go` both are available. We can directly build and run the app or create flogo rule app and run it.

### Run kafka and zookeeper
Open a terminal and navigate to kafka home.
```sh
cd $KAFKA_HOME
#run zookeeper
bin/zookeeper-server-start.sh config/zookeeper.properties
``` 
Open another terminal.
```sh
cd $KAFKA_HOME
#run kafka server
bin/kafka-server-start.sh config/server.properties
```

### Create topic
Open terminal and create topic which is given in kafka trigger. Here `orderinfo` is topic name used in trigger.
```sh
cd $KAFKA_HOME
bin/kafka-topics.sh --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic orderinfo

```

### Direct build and run
```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/simple-kafka
go build
./simple-kafka
```
OR
### Create app using flogo cli
```sh
cd $GOPATH/src/github.com/project-flogo/rules/examples/flogo/simple-kafka
flogo create -f flogo.json simpleKafkaRulesApp
cp functions.go simpleKafkaRulesApp/src
cd simpleKafkaRulesApp
flogo build
./bin/simpleKafkaRulesApp
```

### Send the data over kafka producer
Open terminal send order information through kafka producer.
```sh
cd $KAFKA_HOME
bin/kafka-console-producer.sh --broker-list localhost:9092 --topic orderinfo
>{"type":"grocery","totalPrice":"2001.0"}
```
Expected output in rules terminal.
```
Rule fired: [groceryCheckRule]
Congratulations you are eligible for Rs. 500 gift coupon
```