package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

//1 rtc->one assert triggers two rule actions each rule action deletes tuples.Verify Deleted Tuple types and Tuples count.
func Test_T13(t *testing.T) {

	rs, err := createRuleSession(t)
	assert.Nil(t, err)

	rule := ruleapi.NewRule("R13")
	err = rule.AddCondition("R13_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, r13_action))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rule1 := ruleapi.NewRule("R132")
	err = rule1.AddCondition("R132_c1", []string{"t3.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule1.SetActionService(createActionServiceFromFunction(t, r132_action))
	rule1.SetPriority(2)
	err = rs.AddRule(rule1)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule1.GetName())

	txnCtx := txnCtx{t, 0}
	rs.RegisterRtcTransactionHandler(t13Handler, &txnCtx)
	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)

	t2, err := model.NewTupleWithKeyValues("t3", "t12")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t2)
	assert.Nil(t, err)

	deleteRuleSession(t, rs)

}

func r13_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id, _ := t1.GetString("id")

	if id == "t12" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t11 := rs.GetAssertedTuple(ctx, tk).(model.MutableTuple)
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
		t12 := rs.GetAssertedTuple(ctx, tk).(model.MutableTuple)
		if t12 != nil {
			rs.Delete(ctx, t12)
		}
	}
}

func t13Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {
	if done {
		return
	}

	txnCtx := handlerCtx.(*txnCtx)
	txnCtx.TxnCnt = txnCtx.TxnCnt + 1
	t := txnCtx.Testing
	if txnCtx.TxnCnt == 2 {
		lA := len(rtxn.GetRtcAdded())
		if lA != 1 {
			t.Errorf("RtcAdded: Types expected [%d], got [%d]\n", 1, lA)
			printTuples(t, "Added", rtxn.GetRtcAdded())
		} else {
			tuples, _ := rtxn.GetRtcAdded()["t1"]
			if tuples != nil {
				if len(tuples) != 0 {
					t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 0, len(tuples))
					printTuples(t, "Added", rtxn.GetRtcAdded())
				}
			}
			tuples3, _ := rtxn.GetRtcAdded()["t3"]
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
			tuples, _ := rtxn.GetRtcDeleted()["t1"]
			if tuples != nil {
				if len(tuples) != 1 {
					t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 1, len(tuples))
					printTuples(t, "Deleted", rtxn.GetRtcDeleted())
				}
			}
			tuples3, _ := rtxn.GetRtcDeleted()["t3"]
			if tuples3 != nil {
				if len(tuples3) != 1 {
					t.Errorf("RtcDeleted: Expected [%d], got [%d]\n", 1, len(tuples3))
					printTuples(t, "Deleted", rtxn.GetRtcDeleted())
				}
			}
		}
	}
}
