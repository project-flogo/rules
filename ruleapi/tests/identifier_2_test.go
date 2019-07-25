package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

var cnt int

//Using 3 Identifiers, different Join conditions and triggering respective actions --->Verify order of actions and count.
func Test_I2(t *testing.T) {

	rs, _ := createRuleSession()

	//actionMap := make(map[string]string)

	rule := ruleapi.NewRule("I21")
	rule.AddCondition("I2_c1", []string{"t1.none", "t2.none"}, trueCondition, nil)
	rule.SetAction(i21_action)
	rule.SetPriority(1)
	//rule.SetContext(actionMap)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rule1 := ruleapi.NewRule("I22")
	rule1.AddCondition("I2_c2", []string{"t1.none", "t3.none"}, trueCondition, nil)
	rule1.SetAction(i22_action)
	rule1.SetPriority(1)
	//rule.SetContext(actionMap)
	rs.AddRule(rule1)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	rule2 := ruleapi.NewRule("I23")
	rule2.AddCondition("I2_c3", []string{"t2.none", "t3.none"}, trueCondition, nil)
	rule2.SetAction(i23_action)
	rule2.SetPriority(1)
	//rule.SetContext(actionMap)
	rs.AddRule(rule2)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	rule3 := ruleapi.NewRule("I24")
	rule3.AddCondition("I2_c4", []string{"t1.none", "t2.none", "t3.none"}, trueCondition, nil)
	rule3.SetAction(i24_action)
	rule3.SetPriority(1)
	//rule.SetContext(actionMap)
	rs.AddRule(rule3)
	t.Logf("Rule added: [%s]\n", rule2.GetName())

	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.Assert(context.TODO(), t1)

	t2, _ := model.NewTupleWithKeyValues("t2", "t11")
	rs.Assert(context.TODO(), t2)

	if cnt != 1 {
		t.Errorf("Expecting [1] actions, got [%d]", cnt)
		t.FailNow()
	}

	t3, _ := model.NewTupleWithKeyValues("t3", "t12")
	rs.Assert(context.TODO(), t3)

	if cnt != 2 {
		t.Errorf("Expecting [2] actions, got [%d]", cnt)
		t.FailNow()
	}

	t4, _ := model.NewTupleWithKeyValues("t2", "t13")
	rs.Assert(context.TODO(), t4)

	if cnt != 5 {
		t.Errorf("Expecting [5] actions, got [%d]", cnt)
		t.FailNow()
	}

	rs.Unregister()

}

func i21_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t2")].(model.MutableTuple)
	id2, _ := t2.GetString("id")

	if id1 == "t10" && id2 == "t11" && cnt == 0 {
		cnt++
	}

	if id1 == "t10" && id2 == "t13" {
		if cnt >= 2 && cnt <= 4 {
			cnt++
		}
	}
}

func i22_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id3, _ := t2.GetString("id")

	if id1 == "t10" && id3 == "t12" && cnt == 1 {
		cnt++
	}
}

func i23_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t2")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id3, _ := t2.GetString("id")

	if id1 == "t13" && id3 == "t12" {
		if cnt >= 2 && cnt <= 4 {
			cnt++
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
		if cnt >= 2 && cnt <= 4 {
			cnt++
		}
	}
}
