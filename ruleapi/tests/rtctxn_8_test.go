package tests

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"context"
	"testing"
)

//no-identifier condition
func Test_T8(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R1")
	rule.AddCondition("R1_c1", []string{"t1.none"}, trueCondition, nil)
	rule.AddCondition("R1_c2", []string{}, falseCondition, nil)
	rule.SetAction(assertTuple)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.RegisterRtcTransactionHandler(t8Handler, t)
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t1")
	rs.Assert(context.TODO(), t1)
	rs.Unregister()

}


func assertTuple(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t, _:= model.NewTupleWithKeyValues("t1", "t2")
	rs.Assert(ctx, t)
}

func t8Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	t := handlerCtx.(*testing.T)
	if m, found := rtxn.GetRtcAdded()["t1"]; found {
		lA := len(m)
		if lA != 1 {
			t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 1, lA)
			printTuples(t,"Added", rtxn.GetRtcAdded())
		}
	}
	lM := len(rtxn.GetRtcModified())
	if lM != 0 {
		t.Errorf("RtcModified: Expected [%d], got [%d]\n", 0, lM)
		printModified(t, rtxn.GetRtcModified())
	}
	lD := len(rtxn.GetRtcDeleted())
	if lD != 0 {
		t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 0, lD)
		printTuples(t,"Deleted", rtxn.GetRtcDeleted())
	}
}
