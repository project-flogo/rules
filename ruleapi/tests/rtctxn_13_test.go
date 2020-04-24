package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//1 rtc->one assert triggers two rule actions each rule action deletes tuples.Verify Deleted Tuple types and Tuples count.
func Test_T13(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R13")
	rule.AddCondition("R13_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	rule.SetAction(r13_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rule1 := ruleapi.NewRule("R132")
	rule1.AddCondition("R132_c1", []string{"t3.none"}, trueCondition, nil)
	rule1.SetAction(r132_action)
	rule1.SetPriority(2)
	rs.AddRule(rule1)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	txnCtx := txnCtx{t, 0}
	rs.RegisterRtcTransactionHandler(t13Handler, &txnCtx)
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.Assert(context.TODO(), t1)

	t2, _ := model.NewTupleWithKeyValues("t3", "t12")
	rs.Assert(context.TODO(), t2)

	rs.Unregister()

}

func r13_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id, _ := t1.GetString("id")

	if id == "t12" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t11 := rs.GetAssertedTuple(tk).(model.MutableTuple)
		if t11 != nil {
			rs.Delete(ctx, t11)
		}
	}
}

func r132_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id, _ := t1.GetString("id")

	if id == "t12" {
		tk, _ := model.NewTupleKeyWithKeyValues("t3", "t12")
		t12 := rs.GetAssertedTuple(tk).(model.MutableTuple)
		if t12 != nil {
			rs.Delete(ctx, t12)
		}
	}
}

func t13Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	txnCtx := handlerCtx.(*txnCtx)
	txnCtx.TxnCnt = txnCtx.TxnCnt + 1
	t := txnCtx.Testing
	if txnCtx.TxnCnt == 2 {
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
		if lM != 0 {
			t.Errorf("RtcModified: Expected [%d], got [%d]\n", 0, lM)
			printModified(t, rtxn.GetRtcModified())
		}
		lD := len(rtxn.GetRtcDeleted())
		if lD != 2 {
			t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 2, lD)
			printTuples(t, "Deleted", rtxn.GetRtcDeleted())
		} else {
			tuples := rtxn.GetRtcDeleted()["t1"]
			if tuples != nil {
				if len(tuples) != 1 {
					t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 1, len(tuples))
					printTuples(t, "Deleted", rtxn.GetRtcDeleted())
				}
			}
			tuples3 := rtxn.GetRtcDeleted()["t3"]
			if tuples3 != nil {
				if len(tuples3) != 1 {
					t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 1, len(tuples3))
					printTuples(t, "Deleted", rtxn.GetRtcDeleted())
				}
			}
		}
	}
}
