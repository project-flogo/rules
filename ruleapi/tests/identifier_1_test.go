package tests

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"

	"github.com/stretchr/testify/assert"
)

var count uint64

//Check if all combination of tuples t1 and t3 are triggering actions
func Test_I1(t *testing.T) {
	count = 0
	rs, err := createRuleSession()
	assert.Nil(t, err)

	rule := ruleapi.NewRule("I1")
	err = rule.AddCondition("I1_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	assert.Nil(t, err)
	rule.SetActionService(createActionServiceFromFunction(t, i1_action))

	rule.SetPriority(1)
	err = rs.AddRule(rule)
	assert.Nil(t, err)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	err = rs.Start(nil)
	assert.Nil(t, err)

	t1, err := model.NewTupleWithKeyValues("t1", "t10")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t1)
	assert.Nil(t, err)

	t2, err := model.NewTupleWithKeyValues("t1", "t11")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t2)
	assert.Nil(t, err)

	t3, err := model.NewTupleWithKeyValues("t3", "t12")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t3)
	assert.Nil(t, err)

	//Check if the 2 combinations {t1=t10,t3=t12} and {t1=t11,t3=t12} triggers action twice
	if cnt := atomic.LoadUint64(&count); cnt != 2 {
		t.Errorf("Expecting [2] actions, got [%d]", cnt)
		t.FailNow()
	}

	t4, err := model.NewTupleWithKeyValues("t3", "t13")
	assert.Nil(t, err)
	err = rs.Assert(context.TODO(), t4)
	assert.Nil(t, err)

	//Check if the 2 combinations {t1=t10,t3=t13} and {t1=t11,t3=t13} triggers action two more times making the total action count 4
	if cnt := atomic.LoadUint64(&count); cnt != 4 {
		t.Errorf("Expecting [4] actions, got [%d]", cnt)
		t.FailNow()
	}

	deleteRuleSession(t, rs, t1, t2, t3, t4)

}

func i1_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id1, _ := t1.GetString("id")

	t2 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id3, _ := t2.GetString("id")

	if id1 == "t11" && id3 == "t12" {
		atomic.AddUint64(&count, 1)
	}
	if id1 == "t10" && id3 == "t12" {
		atomic.AddUint64(&count, 1)
	}
	if id1 == "t11" && id3 == "t13" {
		atomic.AddUint64(&count, 1)
	}
	if id1 == "t10" && id3 == "t13" {
		atomic.AddUint64(&count, 1)
	}
}
