package model

//StreamSource represents a source of streaming data such as Kafka or MQTT etc.
//For now, it is a simple typedef. Later on, it can expand to keep meta information of the stream
//such as its datatypes and formats for possible validation

type StreamSource string
