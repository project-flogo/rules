package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

//TTL != 0 asserted
func Test_T1(t *testing.T) {
	rs, err := createRuleSession()
	assert.Nil(t, err)
	rule := ruleapi.NewRule("R1")
	err = rule.AddCondition("R1_c1", []string{"t1.none"}, trueCondition, t)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, emptyAction))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.RegisterRtcTransactionHandler(t1Handler, t)
	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t1", "t1")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)
	deleteRuleSession(t, rs, t1)

}

func t1Handler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {
	if done {
		return
	}
	t := handlerCtx.(*testing.T)

	lA := len(rtxn.GetRtcAdded())
	if lA != 1 {
		t.Errorf("RtcAdded: Expected [%d], got [%d]\n", 1, lA)
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
