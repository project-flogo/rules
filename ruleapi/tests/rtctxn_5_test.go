package tests

import (
	"context"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"testing"
)

//3 rtcs, 1st rtc ->asserted, 2nd rtc ->modified the 1st one, 3rd rtc ->deleted the 2nd one
func Test_T5(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R5")
	rule.AddCondition("R5_c1", []string{"t1.none"}, trueCondition, nil)
	rule.SetAction(r5_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	txnCtx := txnCtx{t, 0}
	rs.RegisterRtcTransactionHandler(t5Handler, &txnCtx)
	rs.Start(nil)

	i1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.Assert(context.TODO(), i1)

	i2, _ := model.NewTupleWithKeyValues("t1", "t11")
	rs.Assert(context.TODO(), i2)

	i3, _ := model.NewTupleWithKeyValues("t1", "t13")
	rs.Assert(context.TODO(), i3)

	rs.Unregister()

}

func r5_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	//t1.SetString(ctx, "p3", "v3")
	id, _ := t1.GetString("id")
	if id == "t11" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t10 := rs.GetAssertedTuple(tk).(model.MutableTuple)
		if t10 != nil {
			t10.SetString(ctx, "p3", "v3")
			t10.SetDouble(ctx, "p2", 11.11)
		}
	} else if id == "t13" {
		//delete t11
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t11")
		t11 := rs.GetAssertedTuple(tk).(model.MutableTuple)
		if t11 != nil {
			rs.Delete(ctx, t11)
		}
	}
}

func t5Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	txnCtx := handlerCtx.(*txnCtx)
	txnCtx.TxnCnt = txnCtx.TxnCnt + 1
	t := txnCtx.Testing
	if txnCtx.TxnCnt == 1 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
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
	} else if txnCtx.TxnCnt == 2 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		}
		lM := len(rtxn.GetRtcModified())
		if lM != 1 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 1, lM)
			printModified(t, rtxn.GetRtcModified())
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 0 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 0, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		}
	} else if txnCtx.TxnCnt == 3 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		}
		lM := len(rtxn.GetRtcModified())
		if lM != 0 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 0, lM)
			printModified(t, rtxn.GetRtcModified())
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 1 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 1, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		}
	}
}
