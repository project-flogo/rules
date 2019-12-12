package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

//1 rtc ->Delete multiple tuple types and verify count.
func Test_T10(t *testing.T) {

	rs, err := createRuleSession(t)
	assert.Nil(t, err)

	rule := ruleapi.NewRule("R10")
	err = rule.AddCondition("R10_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, r10_action))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	txnCtx := txnCtx{t, 0}
	rs.RegisterRtcTransactionHandler(t10Handler, &txnCtx)
	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)

	t3, err := model.NewTupleWithKeyValues("t3", "t11")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t3)
	assert.Nil(t, err)

	deleteRuleSession(t, rs)

}

func r10_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t3 := tuples[model.TupleType("t3")].(model.MutableTuple)
	//t1.SetString(ctx, "p3", "v3")
	id, _ := t3.GetString("id")
	if id == "t11" {
		tk, _ := model.NewTupleKeyWithKeyValues("t1", "t10")
		t10 := rs.GetAssertedTuple(ctx, tk).(model.MutableTuple)
		if t10 != nil {
			rs.Delete(ctx, t10)
		}

		tk1, _ := model.NewTupleKeyWithKeyValues("t3", "t11")
		t11 := rs.GetAssertedTuple(ctx, tk1).(model.MutableTuple)
		if t11 != nil {
			rs.Delete(ctx, t11)
		}
	}
}

func t10Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {
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
