package tests

import (
	"context"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

// Add previously removed rule
func Test_Four(t *testing.T) {
	rs, _ := createRuleSession()

	// create rule
	r1 := ruleapi.NewRule("Rule1")
	r1.AddCondition("r1c1", []string{"t1.none", "t2.none"}, trueCondition, nil)
	r1.SetAction(r1Action)
	// create tuples
	t1, _ := model.NewTupleWithKeyValues("t1", "one") // No TTL
	t2, _ := model.NewTupleWithKeyValues("t2", "two") // TTL is 0

	// start rule session
	rs.Start(nil)

	assertCtxValues := make(map[string]interface{})
	assertCtxValues["test"] = t
	assertCtxValues["actionFired"] = "NOTFIRED"
	assertCtx := context.WithValue(context.TODO(), "values", assertCtxValues)

	// Test1: add r1 and assert t1 & t2; rule action SHOULD be fired
	addRule(t, rs, r1)
	assert(assertCtx, rs, t1)
	assert(assertCtx, rs, t2)
	actionFired, _ := assertCtxValues["actionFired"].(string)
	if actionFired != "FIRED" {
		t.Log("rule action SHOULD be fired")
		t.FailNow()
	}

	// Test2: remove r1 and assert t2; rule action SHOULD NOT be fired
	deleteRule(t, rs, r1)
	assertCtxValues["actionFired"] = "NOTFIRED"
	assert(assertCtx, rs, t2)
	actionFired, _ = assertCtxValues["actionFired"].(string)
	if actionFired != "NOTFIRED" {
		t.Log("rule action SHOULD NOT be fired")
		t.FailNow()
	}

	// Test3: add r1 again and assert t2; rule action SHOULD be fired
	addRule(t, rs, r1)
	rs.ReplayTuplesForRule(r1.GetName())
	assertCtxValues["actionFired"] = "NOTFIRED"
	assert(assertCtx, rs, t2)
	actionFired, _ = assertCtxValues["actionFired"].(string)
	if actionFired != "FIRED" {
		t.Log("rule action SHOULD be fired")
		t.FailNow()
	}

	rs.Unregister()
}

func addRule(t *testing.T, rs model.RuleSession, rule model.Rule) {
	err := rs.AddRule(rule)
	if err != nil {
		t.Logf("[%s] error while adding rule[%s]", time.Now().Format("15:04:05.999999"), rule.GetName())
		return
	}
	t.Logf("[%s] Rule[%s] added. \n", time.Now().Format("15:04:05.999999"), rule.GetName())
}

func deleteRule(t *testing.T, rs model.RuleSession, rule model.Rule) {
	rs.DeleteRule(rule.GetName())
	t.Logf("[%s] Rule[%s] deleted. \n", time.Now().Format("15:04:05.999999"), rule.GetName())
}

func assert(ctx context.Context, rs model.RuleSession, tuple model.Tuple) {
	assertCtxValues := ctx.Value("values").(map[string]interface{})
	t, _ := assertCtxValues["test"].(*testing.T)
	err := rs.Assert(ctx, tuple)
	if err != nil {
		t.Logf("[%s] assert error %s", time.Now().Format("15:04:05.999999"), err)
	}
}

func r1Action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	assertCtxValues := ctx.Value("values").(map[string]interface{})
	test, _ := assertCtxValues["test"].(*testing.T)

	t := tuples["t1"]
	tSrt, _ := t.GetString("id")
	test.Logf("[%s] r1Action called with the tuple[%s] \n", time.Now().Format("15:04:05.999999"), tSrt)

	assertCtxValues["actionFired"] = "FIRED"
}
