package tests

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

var actionCnt uint64

//1 rtc->Scheduled assert, Action should be fired after the delay time.
func Test_T15(t *testing.T) {
	actionCnt = 0
	rs, err := createRuleSession(t)
	assert.Nil(t, err)

	rule := ruleapi.NewRule("R15")
	err = rule.AddCondition("R15_c1", []string{"t1.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, r15_action))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	rs.ScheduleAssert(context.TODO(), 1000, "1", t1)

	if count := atomic.LoadUint64(&actionCnt); count != 0 {
		t.Errorf("Expecting [0] actions, got [%d]", count)
		t.FailNow()
	}
	time.Sleep(2000 * time.Millisecond)

	if count := atomic.LoadUint64(&actionCnt); count != 1 {
		t.Errorf("Expecting [1] actions, got [%d]", count)
		t.FailNow()
	}

	deleteRuleSession(t, rs, t1)

}

func r15_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	atomic.AddUint64(&actionCnt, 1)
}
