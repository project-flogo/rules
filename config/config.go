package config

import (
	"bytes"
	"encoding/json"
	"strconv"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/rules/common/model"
)

// RuleSessionDescriptor is a collection of rules to be loaded

type RuleActionDescriptor struct {
	Name       string               `json:"name"`
	IOMetadata *metadata.IOMetadata `json:"metadata"`
	Rules      []*RuleDescriptor    `json:"rules"`
	Services   []*ServiceDescriptor `json:"services,omitempty"`
}

type RuleSessionDescriptor struct {
	Rules []*RuleDescriptor `json:"rules"`
}

// RuleDescriptor defines a rule
type RuleDescriptor struct {
	Name          string
	Conditions    []*ConditionDescriptor
	ActionFunc    model.ActionFunction
	ActionService *ActionServiceDescriptor
	Priority      int
}

// ConditionDescriptor defines a condition in a rule
type ConditionDescriptor struct {
	Name        string
	Identifiers []string
	Evaluator   model.ConditionEvaluator
}

// ActionServiceDescriptor defines rule action service
type ActionServiceDescriptor struct {
	Service string                 `json:"service"`
	Input   map[string]interface{} `json:"input,omitempty"`
}

// ServiceDescriptor defines a functional target that may be invoked by a rule
type ServiceDescriptor struct {
	Name        string
	Description string
	Function    model.ActionFunction
	Ref         string
	Settings    map[string]interface{}
}

// UnmarshalJSON unmarshals JSON into struct RuleDescriptor
func (c *RuleDescriptor) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name          string                   `json:"name"`
		Conditions    []*ConditionDescriptor   `json:"conditions"`
		ActionFuncId  string                   `json:"actionFunction"`
		ActionService *ActionServiceDescriptor `json:"actionService,omitempty"`
		Priority      int                      `json:"priority"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	c.Name = ser.Name
	c.Conditions = ser.Conditions
	c.ActionFunc = GetActionFunction(ser.ActionFuncId)
	c.ActionService = ser.ActionService
	c.Priority = ser.Priority

	return nil
}

// MarshalJSON returns JSON encoding of RuleDescriptor
func (c *RuleDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + c.Name + "\",")

	buffer.WriteString("\"conditions\":[")
	for _, condition := range c.Conditions {
		jsonCondition, err := condition.MarshalJSON()
		if err == nil {
			buffer.WriteString(string(jsonCondition) + ",")
		}
	}
	buffer.Truncate(buffer.Len() - 1)
	buffer.WriteString("],")

	actionFunctionID := GetActionFunctionID(c.ActionFunc)
	buffer.WriteString("\"actionFunction\":\"" + actionFunctionID + "\",")
	jsonActionService, err := json.Marshal(c.ActionService)
	if err == nil {
		buffer.WriteString("\"actionService\":" + string(jsonActionService) + ",")
	}
	buffer.WriteString("\"priority\":" + strconv.Itoa(c.Priority) + "}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals JSON into struct ConditionDescriptor
func (c *ConditionDescriptor) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name        string   `json:"name"`
		Identifiers []string `json:"identifiers"`
		EvaluatorId string   `json:"evaluator"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	c.Name = ser.Name
	c.Identifiers = ser.Identifiers
	c.Evaluator = GetConditionEvaluator(ser.EvaluatorId)

	return nil
}

// MarshalJSON returns JSON encoding of ConditionDescriptor
func (c *ConditionDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + c.Name + "\",")
	buffer.WriteString("\"identifiers\":[")
	for _, id := range c.Identifiers {
		buffer.WriteString("\"" + id + "\",")
	}
	buffer.Truncate(buffer.Len() - 1)
	buffer.WriteString("],")

	conditionEvaluatorID := GetConditionEvaluatorID(c.Evaluator)
	buffer.WriteString("\"evaluator\":\"" + conditionEvaluatorID + "\"}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals JSON into struct ServiceDescriptor
func (sd *ServiceDescriptor) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description,omitempty"`
		FunctionID  string                 `json:"function,omitempty"`
		Ref         string                 `json:"ref"`
		Settings    map[string]interface{} `json:"settings,omitempty"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	sd.Name = ser.Name
	sd.Description = ser.Description
	sd.Function = GetActionFunction(ser.FunctionID)
	sd.Ref = ser.Ref
	sd.Settings = ser.Settings

	return nil
}

// MarshalJSON returns JSON encoding of ServiceDescriptor
func (sd *ServiceDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + sd.Name + "\",")
	buffer.WriteString("\"description\":" + "\"" + sd.Description + "\",")
	functionID := GetActionFunctionID(sd.Function)
	buffer.WriteString("\"function\":" + "\"" + functionID + "\",")
	buffer.WriteString("\"ref\":" + "\"" + sd.Ref + "\",")
	jsonSettings, err := json.Marshal(sd.Settings)
	if err == nil {
		buffer.WriteString("\"settings\":" + string(jsonSettings) + "}")
	}

	return buffer.Bytes(), nil
}

//metadata support
type DefinitionConfig struct {
	Name     string               `json:"name"`
	Metadata *metadata.IOMetadata `json:"metadata"`
	Rules    []*RuleDescriptor    `json:"rules"`
}
