package tests

import (
	"context"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

var actionCnt int

//1 rtc->Scheduled assert, Action should be fired after the delay time.
func Test_T15(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R15")
	rule.AddCondition("R15_c1", []string{"t1.none"}, trueCondition, nil)
	rule.SetAction(r15_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.ScheduleAssert(context.TODO(), 1000, "1", t1)

	if actionCnt != 0 {
		t.Errorf("Expecting [0] actions, got [%d]", actionCnt)
		t.FailNow()
	}
	time.Sleep(2000 * time.Millisecond)

	if actionCnt != 1 {
		t.Errorf("Expecting [1] actions, got [%d]", actionCnt)
		t.FailNow()
	}

	rs.Unregister()

}

func r15_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	actionCnt++
}
