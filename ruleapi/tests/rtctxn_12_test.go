package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//1 rtc->one assert triggers two rule actions each rule action modifies tuples.Verify Tuple type and Tuples count.
func Test_T12(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R12")
	rule.AddCondition("R12_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	rule.SetAction(r122_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rule1 := ruleapi.NewRule("R122")
	rule1.AddCondition("R122_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	rule1.SetAction(r12_action)
	rule1.SetPriority(1)
	rs.AddRule(rule1)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	txnCtx := txnCtx{t, 0}
	rs.RegisterRtcTransactionHandler(t12Handler, &txnCtx)
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.Assert(context.TODO(), t1)

	t2, _ := model.NewTupleWithKeyValues("t1", "t11")
	rs.Assert(context.TODO(), t2)

	t3, _ := model.NewTupleWithKeyValues("t3", "t12")
	rs.Assert(context.TODO(), t3)

	t4, _ := model.NewTupleWithKeyValues("t3", "t13")
	rs.Assert(context.TODO(), t4)

	rs.Unregister()

}

func r12_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id, _ := t1.GetString("id")

	if id == "t13" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t10 := rs.GetAssertedTuple(tk).(model.MutableTuple)
		t10.SetDouble(ctx, "p2", 11.11)
	}
}

func r122_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id, _ := t1.GetString("id")

	if id == "t13" {
		tk1, _ := model.NewTupleKeyWithKeyValues("t1", "t11")
		t11 := rs.GetAssertedTuple(tk1).(model.MutableTuple)
		t11.SetDouble(ctx, "p2", 11.11)

		tk2, _ := model.NewTupleKeyWithKeyValues("t3", "t12")
		t12 := rs.GetAssertedTuple(tk2).(model.MutableTuple)
		t12.SetDouble(ctx, "p2", 11.11)
	}
}

func t12Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	txnCtx := handlerCtx.(*txnCtx)
	txnCtx.TxnCnt = txnCtx.TxnCnt + 1
	t := txnCtx.Testing
	if txnCtx.TxnCnt == 4 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		} else {
			tuples := rtxn.GetRtcAdded()["t1"]
			if tuples != nil {
				if len(tuples) != 0 {
					t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 0, len(tuples))
					printTuples(t, "Added", rtxn.GetRtcAdded())
				}
			}
			tuples3 := rtxn.GetRtcAdded()["t3"]
			if tuples3 != nil {
				if len(tuples3) != 1 {
					t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 1, len(tuples3))
					printTuples(t, "Added", rtxn.GetRtcAdded())
				}
			}
		}
		lM := len(rtxn.GetRtcModified())
		if lM != 2 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 2, lM)
			printModified(t, rtxn.GetRtcModified())
		} else {
			tuples := rtxn.GetRtcModified()["t1"]
			if tuples != nil {
				if len(tuples) != 2 {
					t.Errorf("RtcModified: Expected [%d], got [%d]\n", 2, len(tuples))
					printModified(t, rtxn.GetRtcModified())
				}
			}
			tuples3 := rtxn.GetRtcModified()["t3"]
			if tuples3 != nil {
				if len(tuples3) != 1 {
					t.Errorf("RtcModified: Expected [%d], got [%d]\n", 1, len(tuples3))
					printModified(t, rtxn.GetRtcModified())
				}
			}
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 0 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 0, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		}
	}
}
