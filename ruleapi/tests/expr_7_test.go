package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/core/data/property"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//1 property variable operation
func Test_7_Expr(t *testing.T) {

	// set test properties
	testProperties := make(map[string]interface{})
	testProperties["name"] = "Bob"
	testProperties["age"] = 10
	propertyManager := property.NewManager(testProperties)
	property.SetDefaultManager(propertyManager)

	actionCount := map[string]int{"count": 0}
	rs, _ := createRuleSession()
	r1 := ruleapi.NewRule("r1")
	r1.AddExprCondition("c1", "($.t1.p3 == $property['name'])", nil)
	r1.AddExprCondition("c2", "($.t1.p1 == $property['age'])", nil)
	r1.SetAction(a7)
	r1.SetContext(actionCount)

	rs.AddRule(r1)

	rs.Start(nil)

	var ctx context.Context

	t1, _ := model.NewTupleWithKeyValues("t1", "t1")
	t1.SetInt(context.TODO(), "p1", 10)
	t1.SetString(context.TODO(), "p3", "Bob")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t1)

	rs.Unregister()
	count := actionCount["count"]
	if count != 1 {
		t.Errorf("expected [%d], got [%d]\n", 1, count)
	}
}

func a7(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value(TestKey{}).(*testing.T)
	t.Logf("Test_7_Expr executed!")
	actionCount := ruleCtx.(map[string]int)
	count := actionCount["count"]
	actionCount["count"] = count + 1
}
