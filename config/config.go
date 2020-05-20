package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/rules/common/model"
)

const (
	// TypeServiceFunction represents go function based rule service
	TypeServiceFunction = "function"
	// TypeServiceActivity represents flgo-activity based rule service
	TypeServiceActivity = "activity"
)

// RuleSessionDescriptor is a collection of rules to be loaded

type RuleActionDescriptor struct {
	Name       string               `json:"name"`
	IOMetadata *metadata.IOMetadata `json:"metadata"`
	Rules      []*RuleDescriptor    `json:"rules"`
	Services   []*ServiceDescriptor `json:"services,omitempty"`
}

type RuleSessionDescriptor struct {
	Rules    []*RuleDescriptor    `json:"rules"`
	Services []*ServiceDescriptor `json:"services,omitempty"`
}

// RuleDescriptor defines a rule
type RuleDescriptor struct {
	Name          string
	Conditions    []*ConditionDescriptor
	ActionService *ActionServiceDescriptor
	Priority      int
	Identifiers   []string
}

// ConditionDescriptor defines a condition in a rule
type ConditionDescriptor struct {
	Name        string
	Identifiers []string
	Evaluator   model.ConditionEvaluator
	Expression  string
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
	Type        string
	Function    model.ActionFunction
	Ref         string
	Settings    map[string]interface{}
}

// UnmarshalJSON unmarshals JSON into struct RuleDescriptor
func (c *RuleDescriptor) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name          string                   `json:"name"`
		Conditions    []*ConditionDescriptor   `json:"conditions"`
		ActionService *ActionServiceDescriptor `json:"actionService,omitempty"`
		Priority      int                      `json:"priority"`
		Identifiers   []string                 `json:"identifiers"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	c.Name = ser.Name
	c.Conditions = ser.Conditions
	c.ActionService = ser.ActionService
	c.Priority = ser.Priority
	c.Identifiers = ser.Identifiers

	return nil
}

// MarshalJSON returns JSON encoding of RuleDescriptor
func (c *RuleDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + c.Name + "\",")
	if c.Identifiers != nil {
		buffer.WriteString("\"identifiers\":[")
		for _, id := range c.Identifiers {
			buffer.WriteString("\"" + id + "\",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString("],")
	}

	buffer.WriteString("\"conditions\":[")
	for _, condition := range c.Conditions {
		jsonCondition, err := condition.MarshalJSON()
		if err == nil {
			buffer.WriteString(string(jsonCondition) + ",")
		}
	}
	buffer.Truncate(buffer.Len() - 1)
	buffer.WriteString("],")

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
		Expression  string   `json:"expression"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	c.Name = ser.Name
	c.Identifiers = ser.Identifiers
	c.Evaluator = GetConditionEvaluator(ser.EvaluatorId)
	c.Expression = ser.Expression

	return nil
}

// MarshalJSON returns JSON encoding of ConditionDescriptor
func (c *ConditionDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + c.Name + "\",")
	if c.Identifiers != nil {
		buffer.WriteString("\"identifiers\":[")
		for _, id := range c.Identifiers {
			buffer.WriteString("\"" + id + "\",")
		}
		buffer.Truncate(buffer.Len() - 1)
		buffer.WriteString("],")
	}

	conditionEvaluatorID := GetConditionEvaluatorID(c.Evaluator)
	buffer.WriteString("\"evaluator\":\"" + conditionEvaluatorID + "\",")
	buffer.WriteString("\"expression\":\"" + c.Expression + "\"}")

	return buffer.Bytes(), nil
}

// UnmarshalJSON unmarshals JSON into struct ServiceDescriptor
func (sd *ServiceDescriptor) UnmarshalJSON(d []byte) error {
	ser := &struct {
		Name        string                 `json:"name"`
		Description string                 `json:"description,omitempty"`
		Type        string                 `json:"type"`
		FunctionID  string                 `json:"function,omitempty"`
		Ref         string                 `json:"ref"`
		Settings    map[string]interface{} `json:"settings,omitempty"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	sd.Name = ser.Name
	sd.Description = ser.Description
	if ser.Type == TypeServiceFunction || ser.Type == TypeServiceActivity {
		sd.Type = ser.Type
	} else {
		return fmt.Errorf("unsupported type - '%s' is referenced in the service '%s'", ser.Type, ser.Name)
	}
	if ser.FunctionID != "" {
		fn := GetActionFunction(ser.FunctionID)
		if fn == nil {
			return fmt.Errorf("function - '%s' not found", ser.FunctionID)
		}
		sd.Function = fn
	}
	sd.Ref = ser.Ref
	sd.Settings = ser.Settings

	return nil
}

// MarshalJSON returns JSON encoding of ServiceDescriptor
func (sd *ServiceDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"name\":" + "\"" + sd.Name + "\",")
	buffer.WriteString("\"description\":" + "\"" + sd.Description + "\",")
	buffer.WriteString("\"type\":" + "\"" + sd.Type + "\",")
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
