package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//Retract
func Test_Retract_1(t *testing.T) {

	rs, _ := createRuleSession()

	//create a rule joining t1 and t3
	rule := ruleapi.NewRule("Retract_Test")
	err := rule.AddCondition("R7_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	if err != nil {
		t.Logf("%s", err)
		t.FailNow()
	}
	ruleActionCtx := make(map[string]string)
	rule.SetContext(ruleActionCtx)
	rule.SetAction(assertAction)
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	if err != nil {
		t.Logf("%s", err)
		t.FailNow()
	}
	t.Logf("Rule added: [%s]\n", rule.GetName())

	err = rs.Start(nil)
	if err != nil {
		t.Logf("%s", err)
		t.FailNow()
	}

	// Case1: assert a t1
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		tuple, _ := model.NewTupleWithKeyValues("t1", "t1")
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
	}

	// Case2: assert a t3 so that the rule fires for keys t1 and t3
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		tuple, _ := model.NewTupleWithKeyValues("t3", "t3")
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
		// make sure that rule action got fired by inspecting the rule context
		isActionFired, ok := ruleActionCtx["isActionFired"]
		if !ok || isActionFired != "Fired" {
			t.Log("Case2: rule action not fired")
			t.FailNow()
		}
		delete(ruleActionCtx, "isActionFired") // clear the context for next test case
	}

	// Case3: now retract t3
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		tuple, _ := model.NewTupleWithKeyValues("t3", "t3")
		rs.Retract(ctx, tuple)
	}

	/**
	Case4:
	now assert with same key again, see that the test does not fail, and rule fires
	for keys t1 and t3
	*/
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		tuple, _ := model.NewTupleWithKeyValues("t3", "t3")
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
		// make sure that rule action got fired by inspecting the rule context
		isActionFired, ok := ruleActionCtx["isActionFired"]
		if !ok || isActionFired != "Fired" {
			t.Log("Case4: rule action not fired")
			t.FailNow()
		}
		delete(ruleActionCtx, "isActionFired") // clear the context for next test case
	}

	/**
	Case5:
	now retract t3 again, just to check if a subsequent t1 does not fire the rule
	there by proving that t3 has been retracted
	*/
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		tuple, _ := model.NewTupleWithKeyValues("t3", "t3")
		rs.Retract(ctx, tuple)
	}

	/**
	Case6:
	now assert another t1 with a different key and observe that rule does not fire
	for keys t3 and t11 (since t3 has been retracted)
	*/
	{
		ctx := context.WithValue(context.TODO(), "key", t)
		tuple, _ := model.NewTupleWithKeyValues("t1", "t11")
		err := rs.Assert(ctx, tuple)
		if err != nil {
			t.Logf("%s", err)
			t.FailNow()
		}
		// make sure that rule action doesn't fire by inspecting the rule context
		_, ok := ruleActionCtx["isActionFired"]
		if ok {
			t.Log("Case6: rule action should not fire")
			t.FailNow()
		}
	}
	rs.Unregister()

}

func assertAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value("key").(*testing.T)
	t1 := tuples["t1"]
	t3 := tuples["t3"]
	t.Logf("Rule fired.. [%s], [%s]\n", t1.GetKey().String(), t3.GetKey().String())

	// add isActionFired to rule context
	usableCtx := ruleCtx.(map[string]string)
	usableCtx["isActionFired"] = "Fired"
}
