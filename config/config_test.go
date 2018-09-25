package config

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/stretchr/testify/assert"
)

var testRuleSessionDescriptorJson = `{
        "rules": [
          {
            "name": "n1.name == Bob",
            "conditions": [
              {
                "name": "c1",
                "identifiers": [ "n1" ],
                "evaluator": "checkForBob"
              }
            ],
            "actionFunction": "checkForBobAction"
          },
          {
            "name": "n1.name == Bob && n1.name == n2.name",
            "conditions": [
              {
                "name": "c1",
                "identifiers": [
                  "n1"
                ],
                "evaluator": "checkForBob"
              },
              {
                "name": "c2",
                "identifiers": [ "n1", "n2" ],
                "evaluator": "checkSameNamesCondition"
              }
            ],
            "actionFunction": "checkSameNamesAction"
          }
        ]
      }
`

func TestDeserialize(t *testing.T) {

	RegisterActionFunction("checkForBobAction", checkForBobAction)
	RegisterActionFunction("checkSameNamesAction", checkSameNamesAction)

	RegisterConditionEvaluator("checkForBob", checkForBob)
	RegisterConditionEvaluator("checkSameNamesCondition", checkSameNamesCondition)

	ruleSessionDescriptor := &RuleSessionDescriptor{}
	err := json.Unmarshal([]byte(testRuleSessionDescriptorJson), ruleSessionDescriptor)

	assert.Nil(t, err)
	assert.NotNil(t, ruleSessionDescriptor.Rules)
	assert.Equal(t, 2, len(ruleSessionDescriptor.Rules))

	r1Cfg := ruleSessionDescriptor.Rules[0]

	assert.Equal(t, "n1.name == Bob", r1Cfg.Name)
	assert.NotNil(t, r1Cfg.Conditions)
	assert.Equal(t, 1, len(r1Cfg.Conditions))

	sf1 := reflect.ValueOf(checkForBobAction)
	sf2 := reflect.ValueOf(r1Cfg.ActionFunc)
	assert.Equal(t, sf1.Pointer(), sf2.Pointer())

	r1c1Cfg := r1Cfg.Conditions[0]
	assert.Equal(t, "c1", r1c1Cfg.Name)
	assert.NotNil(t, r1c1Cfg.Identifiers)
	assert.Equal(t, 1, len(r1c1Cfg.Identifiers))

	sf1 = reflect.ValueOf(checkForBob)
	sf2 = reflect.ValueOf(r1c1Cfg.Evaluator)
	assert.Equal(t, sf1.Pointer(), sf2.Pointer())

	r2Cfg := ruleSessionDescriptor.Rules[1]

	assert.Equal(t, "n1.name == Bob && n1.name == n2.name", r2Cfg.Name)
	assert.NotNil(t, r2Cfg.Conditions)
	assert.Equal(t, 2, len(r2Cfg.Conditions))

	sf1 = reflect.ValueOf(checkSameNamesAction)
	sf2 = reflect.ValueOf(r2Cfg.ActionFunc)
	assert.Equal(t, sf1.Pointer(), sf2.Pointer())

	r2c1Cfg := r2Cfg.Conditions[0]
	assert.Equal(t, "c1", r2c1Cfg.Name)
	assert.NotNil(t, r2c1Cfg.Identifiers)
	assert.Equal(t, 1, len(r2c1Cfg.Identifiers))

	sf1 = reflect.ValueOf(checkForBob)
	sf2 = reflect.ValueOf(r2c1Cfg.Evaluator)
	assert.Equal(t, sf1.Pointer(), sf2.Pointer())

	r2c2Cfg := r2Cfg.Conditions[1]
	assert.Equal(t, "c2", r2c2Cfg.Name)
	assert.NotNil(t, r2c2Cfg.Identifiers)
	assert.Equal(t, 2, len(r2c2Cfg.Identifiers))

	sf1 = reflect.ValueOf(checkSameNamesCondition)
	sf2 = reflect.ValueOf(r2c2Cfg.Evaluator)
	assert.Equal(t, sf1.Pointer(), sf2.Pointer())
}

// TEST FUNCTIONS

func checkForBob(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}

func checkForBobAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
}

func checkSameNamesCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}

func checkSameNamesAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
}
