package tests

import (
	"context"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"testing"
)

//New asserted in action (forward chain)
func Test_T3(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R3")
	rule.AddCondition("R3_c1", []string{"t1.none"}, trueCondition, nil)
	rule.SetAction(R3_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.RegisterRtcTransactionHandler(t3Handler, t)
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.Assert(context.TODO(), t1)
	rs.Unregister()
}

func R3_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1, _ := model.NewTupleWithKeyValues("t1", "t11")
	rs.Assert(ctx, t1)
}

func t3Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	t := handlerCtx.(*testing.T)

	lA := len(rtxn.GetRtcAdded())
	if lA != 1 {
		t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
		printTuples(t, "Added", rtxn.GetRtcAdded())

	} else {
		//ok
		tuples, _ := rtxn.GetRtcAdded()["t1"]
		if tuples != nil {
			if len(tuples) != 2 {
				t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 2, len(tuples))
				printTuples(t, "Added", rtxn.GetRtcAdded())
			}
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
		printTuples(t, "Deleted", rtxn.GetRtcDeleted())
	}
}
