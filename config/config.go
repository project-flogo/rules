package config

import (
	"encoding/json"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

type RuleSession struct {
	Rules []*Rule `json:"rules"`
}

type Rule struct {
	Name       string
	Conditions []*Condition
	ActionFunc model.ActionFunction
	Priority   int
}

type Condition struct {
	Name        string
	Identifiers []string
	Evaluator   model.ConditionEvaluator
}

func (c *Rule) UnmarshalJSON(d []byte) error {

	ser := &struct {
		Name         string       `json:"name"`
		Conditions   []*Condition `json:"conditions"`
		ActionFuncId string       `json:"actionFunction"`
		Priority     int          `json:"priority"`
	}{}

	if err := json.Unmarshal(d, ser); err != nil {
		return err
	}

	c.Name = ser.Name
	c.Conditions = ser.Conditions
	c.ActionFunc = GetActionFunction(ser.ActionFuncId)
	c.Priority = ser.Priority

	return nil
}

func (c *Condition) UnmarshalJSON(d []byte) error {

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

//todo this should probably move to ruleapi
func GetOrCreateRuleSessionFromConfig(name string, config *RuleSession) (model.RuleSession, error) {

	rs, err := ruleapi.GetOrCreateRuleSession(name)

	if err != nil {
		return nil, err
	}

	for _, ruleCfg := range config.Rules {

		rule := ruleapi.NewRule(ruleCfg.Name)
		rule.SetContext("This is a test of context")
		rule.SetAction(ruleCfg.ActionFunc)
		rule.SetPriority(ruleCfg.Priority)

		for _, condCfg := range ruleCfg.Conditions {
			rule.AddCondition(condCfg.Name, condCfg.Identifiers, condCfg.Evaluator, nil)
		}

		rs.AddRule(rule)
	}

	rs.SetStartupFunction(GetStartupRSFunction(name))

	return rs, nil
}
