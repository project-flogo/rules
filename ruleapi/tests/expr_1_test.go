package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//1 condition, 1 expression
func Test_1_Expr(t *testing.T) {
	actionCount := map[string]int{"count": 0}
	rs, err := createRuleSession()
	if err != nil {
		t.Fatal(err)
	}
	r1 := ruleapi.NewRule("r1")
	err = r1.AddExprCondition("c1", "$.t2.p2 > $.t1.p1", nil)
	if err != nil {
		t.Fatal(err)
	}
	r1.SetActionService(createActionServiceFromFunction(t, a1))
	r1.SetContext(actionCount)
	err = rs.AddRule(r1)
	if err != nil {
		t.Fatal(err)
	}
	err = rs.Start(nil)
	if err != nil {
		t.Fatal(err)
	}
	var ctx context.Context

	t1, err := model.NewTupleWithKeyValues("t1", "t1")
	if err != nil {
		t.Fatal(err)
	}
	err = t1.SetInt(nil, "p1", 1)
	if err != nil {
		t.Fatal(err)
	}
	err = t1.SetDouble(nil, "p2", 1.3)
	if err != nil {
		t.Fatal(err)
	}
	err = t1.SetString(nil, "p3", "t3")
	if err != nil {
		t.Fatal(err)
	}

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	err = rs.Assert(ctx, t1)
	if err != nil {
		t.Fatal(err)
	}

	t2, err := model.NewTupleWithKeyValues("t2", "t2")
	if err != nil {
		t.Fatal(err)
	}
	err = t2.SetInt(nil, "p1", 1)
	if err != nil {
		t.Fatal(err)
	}
	err = t2.SetDouble(nil, "p2", 1.0001)
	if err != nil {
		t.Fatal(err)
	}
	err = t2.SetString(nil, "p3", "t3")
	if err != nil {
		t.Fatal(err)
	}

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	err = rs.Assert(ctx, t2)
	if err != nil {
		t.Fatal(err)
	}

	deleteRuleSession(t, rs, t1)

	count := actionCount["count"]
	if count != 1 {
		t.Errorf("expected [%d], got [%d]\n", 1, count)
	}
}

func a1(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value(TestKey{}).(*testing.T)
	t.Logf("Test_1_Expr executed!")
	actionCount := ruleCtx.(map[string]int)
	count := actionCount["count"]
	actionCount["count"] = count + 1
}
