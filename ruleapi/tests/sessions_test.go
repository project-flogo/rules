package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

func TestClearSessions(t *testing.T) {
	ruleapi.GetOrCreateRuleSession("test", "")
	ruleapi.ClearSessions()
}

func TestAssert(t *testing.T) {
	rs, _ := createRuleSession()
	rule := ruleapi.NewRule("R2")
	rule.AddCondition("R2_c1", []string{"t4.none"}, trueCondition, nil)
	rule.SetActionService(createActionServiceFromFunction(t, emptyAction))
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
	deleteRuleSession(t, rs, t1)
}

func TestRace(t *testing.T) {
	rs, _ := createRuleSession()
	defer rs.Unregister()
	rule := ruleapi.NewRule("R2")
	rule.AddCondition("R2_c1", []string{"t4.none"}, trueCondition, nil)
	rule.SetActionService(createActionServiceFromFunction(t, emptyAction))
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())
	rs.Start(nil)

	done := make(chan bool, 8)
	withTTL := func() {
		for i := 0; i < 10; i++ {
			t1, _ := model.NewTupleWithKeyValues("t4", fmt.Sprintf("ttl%d", i))
			err := rs.Assert(context.TODO(), t1)
			if err != nil {
				t.Fatalf("err should be nil: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}
	withDelete := func() {
		for i := 0; i < 10; i++ {
			t1, _ := model.NewTupleWithKeyValues("t3", fmt.Sprintf("delete%d", i))
			err := rs.Assert(context.TODO(), t1)
			if err != nil {
				t.Fatalf("err should be nil: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
			rs.Delete(context.TODO(), t1)
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}
	addOnly := func() {
		for i := 0; i < 10; i++ {
			t1, _ := model.NewTupleWithKeyValues("t3", fmt.Sprintf("add%d", i))
			err := rs.Assert(context.TODO(), t1)
			if err != nil {
				t.Fatalf("err should be nil: %v", err)
			}
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}
	go withTTL()
	go withDelete()
	go addOnly()
	for i := 0; i < 3; i++ {
		<-done
	}
	time.Sleep(3 * time.Second)
}
