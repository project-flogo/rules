package tests

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"golang.org/x/net/context"
	"testing"
)

//1 condition, 1 expression
func Test_1_Expr(t *testing.T) {

	rs, _ := createRuleSession()
	r1 := ruleapi.NewRule("r1")
	r1.AddExprCondition("c1", "$.t2.p2 > $.t1.p1", nil)
	r1.SetAction(a1)
	rs.AddRule(r1)

	rs.Start(nil)

	var ctx context.Context

	t1, _ := model.NewTupleWithKeyValues("t1", "t1")
	t1.SetInt(nil, "p1", 1)
	t1.SetDouble(nil, "p2", 1.3)
	t1.SetString(nil, "p3", "t3")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t1)

	t2, _ := model.NewTupleWithKeyValues("t2", "t2")
	t2.SetInt(nil, "p1", 1)
	t2.SetDouble(nil, "p2", 1.0001)
	t2.SetString(nil, "p3", "t3")

	ctx = context.WithValue(context.TODO(), TestKey{}, t)
	rs.Assert(ctx, t2)
	rs.Unregister()
}

func a1(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t := ctx.Value(TestKey{}).(*testing.T)
	t.Logf("Test_1_Expr executed!")
}

//
// These standalone tests are not relevant anymore as the expression API has changed
//
//func Test_Eval (t *testing.T) {
//	expr, _ := expression.ParseExpression("1 == 1.23")
//	i, err := expr.Eval()
//	if err != nil {
//		t.Fatalf("error %s\n", err)
//	}
//	res := i.(bool)
//	if res {
//		t.Errorf("Expected false, got : %t\n ", res)
//	}
//}
//
//func Test_Eval2 (t *testing.T) {
//	expr, _ := expression.ParseExpression("1 < 1.23")
//	i, err := expr.Eval()
//	if err != nil {
//		t.Fatalf("error %s\n", err)
//	}
//	res := i.(bool)
//	if !res {
//		t.Errorf("Expected true, got : %t\n ", res)
//	}
//}
//
//func Test_Eval3 (t *testing.T) {
//	expr, _ := expression.ParseExpression("1.23 == 1")
//	i, err := expr.Eval()
//	if err != nil {
//		t.Fatalf("error %s\n", err)
//	}
//	res := i.(bool)
//	if res {
//		t.Errorf("Expected false, got : %t\n ", res)
//	}
//}
//
//func Test_Eval4 (t *testing.T) {
//	expr, _ := expression.ParseExpression("1.23 > 1")
//	i, err := expr.Eval()
//	if err != nil {
//		t.Fatalf("error %s\n", err)
//	}
//	res := i.(bool)
//	if !res {
//		t.Errorf("Expected true, got : %t\n ", res)
//	}
//}
