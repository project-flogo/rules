package tests

import (
	"context"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"testing"
)

//TTL = 0, asserted
func Test_T2(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R2")
	rule.AddCondition("R2_c1", []string{"t2.none"}, trueCondition, nil)
	rule.SetAction(emptyAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.RegisterRtcTransactionHandler(t2Handler, t)
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t2", "t2")
	rs.Assert(context.TODO(), t1)
	rs.Unregister()

}

func t2Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	t := handlerCtx.(*testing.T)

	lA := len(rtxn.GetRtcAdded())
	if lA != 0 {
		t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 0, lA)
		printTuples(t, "Added", rtxn.GetRtcAdded())
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
