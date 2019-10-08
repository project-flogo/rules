package tests

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

var cnt uint64

//Using 3 Identifiers, different Join conditions and triggering respective actions --->Verify order of actions and count.
func Test_I2(t *testing.T) {
	cnt = 0
	rs, err := createRuleSession()
	assert.Nil(t, err)

	//actionMap := make(map[string]string)

	rule := ruleapi.NewRule("I21")
	err = rule.AddCondition("I2_c1", []string{"t1.none", "t2.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, i21_action))
	rule.SetPriority(1)
	//rule.SetContext(actionMap)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rule1 := ruleapi.NewRule("I22")
	err = rule1.AddCondition("I2_c2", []string{"t1.none", "t3.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule1.SetActionService(createActionServiceFromFunction(t, i22_action))
	rule1.SetPriority(1)
	//rule.SetContext(actionMap)
	err = rs.AddRule(rule1)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	rule2 := ruleapi.NewRule("I23")
	err = rule2.AddCondition("I2_c3", []string{"t2.none", "t3.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule2.SetActionService(createActionServiceFromFunction(t, i23_action))
	rule2.SetPriority(1)
	//rule.SetContext(actionMap)
	err = rs.AddRule(rule2)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	rule3 := ruleapi.NewRule("I24")
	err = rule3.AddCondition("I2_c4", []string{"t1.none", "t2.none", "t3.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule3.SetActionService(createActionServiceFromFunction(t, i24_action))
	rule3.SetPriority(1)
	//rule.SetContext(actionMap)
	err = rs.AddRule(rule3)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule2.GetName())

	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)

	t2, err := model.NewTupleWithKeyValues("t2", "t11")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t2)
	assert.Nil(t, err)

	if count := atomic.LoadUint64(&cnt); count != 1 {
		t.Errorf("Expecting [1] actions, got [%d]", count)
		t.FailNow()
	}

	t3, err := model.NewTupleWithKeyValues("t3", "t12")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t3)
	assert.Nil(t, err)

	if count := atomic.LoadUint64(&cnt); count != 2 {
		t.Errorf("Expecting [2] actions, got [%d]", count)
		t.FailNow()
	}

	t4, err := model.NewTupleWithKeyValues("t2", "t13")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t4)
	assert.Nil(t, err)

	if count := atomic.LoadUint64(&cnt); count != 5 {
		t.Errorf("Expecting [5] actions, got [%d]", count)
		t.FailNow()
	}

	deleteRuleSession(t, rs, t1, t3)

}

func i21_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t2")].(model.MutableTuple)
	id2, _ := t2.GetString("id")

	if count := atomic.LoadUint64(&cnt); id1 == "t10" && id2 == "t11" && count == 0 {
		atomic.AddUint64(&cnt, 1)
	}

	if id1 == "t10" && id2 == "t13" {
		if count := atomic.LoadUint64(&cnt); count >= 2 && count <= 4 {
			atomic.AddUint64(&cnt, 1)
		}
	}
}

func i22_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id3, _ := t2.GetString("id")

	if count := atomic.LoadUint64(&cnt); id1 == "t10" && id3 == "t12" && count == 1 {
		atomic.AddUint64(&cnt, 1)
	}
}

func i23_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t2")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id3, _ := t2.GetString("id")

	if id1 == "t13" && id3 == "t12" {
		if count := atomic.LoadUint64(&cnt); count >= 2 && count <= 4 {
			atomic.AddUint64(&cnt, 1)
		}
	}
}

func i24_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t2")].(model.MutableTuple)
	id2, _ := t2.GetString("id")

	t3 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id3, _ := t3.GetString("id")

	if id1 == "t10" && id2 == "t13" && id3 == "t12" {
		if count := atomic.LoadUint64(&cnt); count >= 2 && count <= 4 {
			atomic.AddUint64(&cnt, 1)
		}
	}
}
