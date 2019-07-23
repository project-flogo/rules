package tests

import (
	"context"
	"testing"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

//Check if all combination of tuples t1 and t3 are triggering actions
func Test_T14(t *testing.T) {

	rs, _ := createRuleSession()

	actionMap := make(map[string]string)

	rule := ruleapi.NewRule("I1")
	rule.AddCondition("I1_c1", []string{"t1.none", "t3.none"}, trueCondition, nil)
	rule.SetAction(i1_action)
	rule.SetPriority(1)
	rule.SetContext(actionMap)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t10")
	rs.Assert(context.TODO(), t1)

	t2, _ := model.NewTupleWithKeyValues("t1", "t11")
	rs.Assert(context.TODO(), t2)

	t3, _ := model.NewTupleWithKeyValues("t3", "t12")
	rs.Assert(context.TODO(), t3)

	if len(actionMap) != 2 {
		t.Errorf("Expecting [2] actions, got [%d]", len(actionMap))
		t.FailNow()
	}

	rs.Unregister()

}

func i1_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples[model.TupleType("t1")].(model.MutableTuple)
	id1, _ := t1.GetString("id")
	//fmt.Println(" tuples for type\n", t1)

	t2 := tuples[model.TupleType("t3")].(model.MutableTuple)
	id3, _ := t2.GetString("id")
	//fmt.Println(" tuples for type\n", t2)

	if id1 == "t11" && id3 == "t12" {
		firedMap := ruleCtx.(map[string]string)
		firedMap["A1"] = "Fired"
	}
	if id1 == "t10" && id3 == "t12" {
		firedMap := ruleCtx.(map[string]string)
		firedMap["A2"] = "Fired"
	}
}
