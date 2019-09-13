package tests

import (
	"context"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

func TestClearSessions(t *testing.T) {
	ruleapi.GetOrCreateRuleSession("test")
	ruleapi.ClearSessions()
}

func TestAssert(t *testing.T) {
	rs, _ := createRuleSession()
	rule := ruleapi.NewRule("R2")
	rule.AddCondition("R2_c1", []string{"t4.none"}, trueCondition, nil)
	rule.SetAction(emptyAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t4", "t4")
	err := rs.Assert(context.TODO(), t1)
	if err != nil {
		t.Fatalf("err should be nil: %v", err)
	}
	err = rs.Assert(context.TODO(), t1)
	if err == nil {
		t.Fatalf("err should be not be nil: %v", err)
	}

	time.Sleep(2 * time.Second)
	err = rs.Assert(context.TODO(), t1)
	if err != nil {
		t.Fatalf("err should be nil: %v", err)
	}
	rs.Unregister()
}
