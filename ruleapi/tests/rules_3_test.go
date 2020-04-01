package tests

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

var actCnt uint64

//Forward chain-Data change in r3action and r32action triggers the r32action.
func Test_Three(t *testing.T) {
	actCnt = 0
	rs, err := createRuleSession(t)
	assert.Nil(t, err)

	actionMap := make(map[string]string)

	rule := ruleapi.NewRule("R3")
	err = rule.AddCondition("R3c1", []string{"t1.id"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, r3action))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rule1 := ruleapi.NewRule("R32")
	err = rule1.AddCondition("R32c1", []string{"t1.p1"}, r3Condition, nil)
	assert.Nil(t, err)
	rule1.SetActionService(createActionServiceFromFunction(t, r32action))
	rule1.SetPriority(1)
	rule1.SetContext(actionMap)
	err = rs.AddRule(rule1)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	err = t1.SetInt(context.TODO(), "p1", 2000)
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)

	t2, err := model.NewTupleWithKeyValues("t1", "t11")
	assert.Nil(t, err)
	err = t2.SetInt(context.TODO(), "p1", 2000)
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t2)
	assert.Nil(t, err)

	if count := atomic.LoadUint64(&actCnt); count != 2 {
		t.Errorf("Expecting [2] actions, got [%d]", count)
		t.FailNow()
	}

	deleteRuleSession(t, rs, t1, t2)

}

func r3Condition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	p1, _ := t1.GetInt("p1")
	if p1 < 1000 {
		return true
	} else {
		return false
	}

}

func r3action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id, _ := t1.GetString("id")

	if id == "t11" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t10 := rs.GetAssertedTuple(ctx, tk).(model.MutableTuple)
		t10.SetInt(ctx, "p1", 100)
	}
}

func r32action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	atomic.AddUint64(&actCnt, 1)

	tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
	t10 := rs.GetAssertedTuple(ctx, tk).(model.MutableTuple)
	t10.SetInt(ctx, "p1", 500)
}
