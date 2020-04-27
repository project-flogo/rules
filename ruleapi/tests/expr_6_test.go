package tests

import (
	"context"
	"os"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//1 enviroment variable operation
func Test_6_Expr(t *testing.T) {

	defaultVal := os.Getenv("name_rules_test_6")
	os.Setenv("name_rules_test_6", "test")
	defer func() {
		os.Setenv("name_rules_test_6", defaultVal)
	}()

	actionCount := map[string]int{"count": 0}
	rs, _ := createRuleSession()
	r1 := ruleapi.NewRule("r1")
	r1.AddExprCondition("c1", "($.t1.p3 == $env['name_rules_test_6'])", nil)
	r1.SetAction(a6)
	r1.SetContext(actionCount)

	rs.AddRule(r1)

	rs.Start(nil)

	var ctx context.Context

	t1, _ := model.NewTupleWithKeyValues("t1", "t1")
	t1.SetInt(context.TODO(), "p1", 1)
	t1.SetString(context.TODO(), "p3", "test")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t1)

	rs.Unregister()
	count := actionCount["count"]
	if count != 1 {
		t.Errorf("expected [%d], got [%d]\n", 1, count)
	}

}

func a6(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value(TestKey{}).(*testing.T)
	t.Logf("Test_6_Expr executed!")
	actionCount := ruleCtx.(map[string]int)
	count := actionCount["count"]
	actionCount["count"] = count + 1
}
