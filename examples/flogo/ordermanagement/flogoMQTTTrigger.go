package main

import (
	"context"
	"encoding/json"

	"github.com/TIBCOSoftware/flogo-contrib/trigger/mqtt"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/engine"
	"github.com/TIBCOSoftware/flogo-lib/flogo"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/project-flogo/rules/common/model"
)

// MQTT constants
const (
	baseTopic              = "oms/#"
	broker                 = "tcp://test.mosquitto.org:1883"
	orderEventTopic        = "orderevent"
	itemEventTopic         = "itemevent"
	orderShippedEventTopic = "ordershippedevent"
)

func setupFlogoMQTTTriggers() {
	app := flogo.NewApp()

	mqttSettings := map[string]interface{}{
		"topic":     baseTopic,
		"broker":    broker,
		"id":        "oms",
		"qos":       "0",
		"cleansess": "false",
	}
	trg := app.NewTrigger(&mqtt.MqttTrigger{}, mqttSettings)
	trg.NewFuncHandler(map[string]interface{}{"topic": orderEventTopic}, HandleOrderEvent)
	trg.NewFuncHandler(map[string]interface{}{"topic": itemEventTopic}, HandleItemEvent)
	trg.NewFuncHandler(map[string]interface{}{"topic": orderShippedEventTopic}, HandleOrderShippedEvent)

	e, err := flogo.NewEngine(app)

	if err != nil {
		logger.Error(err)
		return
	}

	engine.RunEngine(e)
}

// HandleOrderEvent accepts and processes orderevent
func HandleOrderEvent(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	processMQTTEvent("orderevent", getValues(inputs))
	return nil, nil
}

// HandleItemEvent accepts and processes itemevent
func HandleItemEvent(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	processMQTTEvent("itemevent", getValues(inputs))
	return nil, nil
}

// HandleOrderShippedEvent accepts and processes ordershippedevent
func HandleOrderShippedEvent(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {
	processMQTTEvent("ordershippedevent", getValues(inputs))
	return nil, nil
}

// retrieves values off the inputs
func getValues(inputs map[string]*data.Attribute) string {
	valAttr, exists := inputs[msgValueField]
	if !exists {
		logger.Debugf("No values recieved")
		return ""
	}
	return valAttr.Value().(string)
}

// processes the data received from the event and asserts to the rule session
func processMQTTEvent(topic string, payload string) error {
	// logger.Infof("Topic [%s] & Payload[%s]", topic, payload)

	tupleType := model.TupleType(topic)
	tupleData := make(map[string]interface{})
	json.Unmarshal([]byte(payload), &tupleData)

	newTuple, err := model.NewTuple(tupleType, tupleData)
	if err != nil {
		logger.Errorf("Failed creating tuple of type [%s] with payload [%s] - %s", tupleType, payload, err.Error())
		panic(err)
	}

	ruleSession.Assert(nil, newTuple)
	// logger.Debug("Tuple Data - " + newTuple.GetKey().String())

	return nil
}
