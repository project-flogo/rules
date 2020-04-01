package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

func TestClearSessions(t *testing.T) {
	_, err := ruleapi.GetOrCreateRuleSession("test", "")
	assert.Nil(t, err)
	ruleapi.ClearSessions()
}

func TestAssert(t *testing.T) {
	rs, err := createRuleSession(t)
	assert.Nil(t, err)
	rule := ruleapi.NewRule("R2")
	err = rule.AddCondition("R2_c1", []string{"t4.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, emptyAction))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())
	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t4", "t4")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.NotNil(t, err)

	time.Sleep(2 * time.Second)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)
	deleteRuleSession(t, rs, t1)
}

func TestRace(t *testing.T) {
	rs, err := createRuleSession(t)
	assert.Nil(t, err)
	defer rs.Unregister()
	rule := ruleapi.NewRule("R2")
	err = rule.AddCondition("R2_c1", []string{"t4.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, emptyAction))
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())
	err = rs.Start(nil)
	assert.Nil(t, err)

	done := make(chan bool, 8)
	withTTL := func() {
		for i := 0; i < 10; i++ {
			t1, err := model.NewTupleWithKeyValues("t4", fmt.Sprintf("ttl%d", i))
			assert.Nil(t, err)
			err = rs.Assert(context.TODO(), t1)
			assert.Nil(t, err)
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}
	withDelete := func() {
		for i := 0; i < 10; i++ {
			t1, err := model.NewTupleWithKeyValues("t3", fmt.Sprintf("delete%d", i))
			assert.Nil(t, err)
			err = rs.Assert(context.TODO(), t1)
			assert.Nil(t, err)
			time.Sleep(10 * time.Millisecond)
			err = rs.Delete(context.TODO(), t1)
			assert.Nil(t, err)
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}
	addOnly := func() {
		for i := 0; i < 10; i++ {
			t1, err := model.NewTupleWithKeyValues("t3", fmt.Sprintf("add%d", i))
			assert.Nil(t, err)
			err = rs.Assert(context.TODO(), t1)
			assert.Nil(t, err)
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
