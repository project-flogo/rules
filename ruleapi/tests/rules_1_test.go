package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

/**
In this test, three rules are configured, each with conditions that do not have
any dependency on any tuple property. Thus when an action for one rule changes a property,
the expected outcome is that all three rules should fire
*/

func Test_One(t *testing.T) {

	rs, err := createRuleSession()
	assert.Nil(t, err)

	actionMap := make(map[string]string)

	//// rule 1
	r1 := ruleapi.NewRule("R1")
	err = r1.AddCondition("C1", []string{"t1"}, checkC1, nil)
	assert.Nil(t, err)
	r1.SetActionService(createActionServiceFromFunction(t, actionA1))
	r1.SetPriority(1)
	r1.SetContext(actionMap)

	err = rs.AddRule(r1)
	assert.Nil(t, err)

	// rule 2
	r2 := ruleapi.NewRule("R2")
	err = r2.AddCondition("C2", []string{"t1"}, checkC2, nil)
	assert.Nil(t, err)
	r2.SetActionService(createActionServiceFromFunction(t, actionA2))
	r2.SetPriority(2)
	r2.SetContext(actionMap)

	err = rs.AddRule(r2)
	assert.Nil(t, err)

	// rule 3
	r3 := ruleapi.NewRule("R3")
	err = r3.AddCondition("C3", []string{"t1"}, checkC3, nil)
	assert.Nil(t, err)
	r3.SetActionService(createActionServiceFromFunction(t, actionA3))
	r3.SetPriority(3)
	r3.SetContext(actionMap)

	err = rs.AddRule(r3)
	assert.Nil(t, err)

	//Start the rule session
	err = rs.Start(nil)
	assert.Nil(t, err)

	//Now assert a "t1" tuple
	t1, err := model.NewTupleWithKeyValues("t1", "Tom")
	assert.Nil(t, err)
	t1.SetString(nil, "p3", "test")
	err = rs.Assert(nil, t1)
	assert.Nil(t, err)

	//unregister the session, i.e; cleanup
	deleteRuleSession(t, rs, t1)

	if len(actionMap) != 3 {
		t.Errorf("Expecting [3] actions, got [%d]", len(actionMap))
		t.FailNow()
	}
}

func checkC1(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	fmt.Println("In Condition C1")
	return true
}

func actionA1(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("In Action A1 Start")

	// change t1 field
	t1 := tuples["t1"].(model.MutableTuple)
	t1.SetString(ctx, "p3", "somethingnew")

	fmt.Println("In Action A1 End")
	firedMap := ruleCtx.(map[string]string)
	firedMap["A1"] = "Fired"
}

func checkC2(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	fmt.Println("In Condition C2")
	return true
}

func actionA2(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("In Action A2 Start")
	t1 := tuples["t1"]
	val, _ := t1.GetString("p3")
	fmt.Println("In Action A2 End ", val)
	firedMap := ruleCtx.(map[string]string)
	firedMap["A2"] = "Fired"
}

func checkC3(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	fmt.Println("In Condition C3")
	return true
}

func actionA3(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("In Action A3 Start")
	t1 := tuples["t1"]
	val, _ := t1.GetString("p3")
	fmt.Println("In Action A3 End ", val)
	firedMap := ruleCtx.(map[string]string)
	firedMap["A3"] = "Fired"
}
