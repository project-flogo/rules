package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

var actCnt int

//Forward chain-Data change in r3action and r32action triggers the r32action.
func Test_Three(t *testing.T) {

	rs, _ := createRuleSession()

	actionMap := make(map[string]string)

	rule := ruleapi.NewRule("R3")
	rule.AddCondition("R3c1", []string{"t1.id"}, trueCondition, nil)
	rule.SetAction(r3action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rule1 := ruleapi.NewRule("R32")
	rule1.AddCondition("R32c1", []string{"t1.p1"}, r3Condition, nil)
	rule1.SetAction(r32action)
	rule1.SetPriority(1)
	rule1.SetContext(actionMap)
	rs.AddRule(rule1)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	t1.SetInt(context.TODO(), "p1", 2000)
	rs.Assert(context.TODO(), t1)

	t2, _ := model.NewTupleWithKeyValues("t1", "t11")
	t2.SetInt(context.TODO(), "p1", 2000)
	rs.Assert(context.TODO(), t2)

	if actCnt != 2 {
		t.Errorf("Expecting [2] actions, got [%d]", actCnt)
		t.FailNow()
	}

	rs.Unregister()

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
	//fmt.Println("r13_action triggered")
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id, _ := t1.GetString("id")

	if id == "t11" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t10 := rs.GetAssertedTuple(tk).(model.MutableTuple)
		t10.SetInt(ctx, "p1", 100)
	}
}

func r32action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	//fmt.Println("r132_action triggered")
	actCnt++

	tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
	t10 := rs.GetAssertedTuple(tk).(model.MutableTuple)
	t10.SetInt(ctx, "p1", 500)
}
