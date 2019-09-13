package tests

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

var actionCnt1 uint64

//1 rtc->Schedule assert, Cancel scheduled assert and action should not be fired
func Test_T16(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R16")
	rule.AddCondition("R16_c1", []string{"t1.none"}, trueCondition, nil)
	rule.SetAction(r16_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.ScheduleAssert(context.TODO(), 1000, "1", t1)
	rs.CancelScheduledAssert(context.TODO(), "1")

	time.Sleep(2000 * time.Millisecond)

	if count := atomic.LoadUint64(&actionCnt1); count != 0 {
		t.Errorf("Expecting [0] actions, got [%d]", count)
		t.FailNow()
	}

	rs.Unregister()

}

func r16_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	atomic.AddUint64(&actionCnt1, 1)
}
