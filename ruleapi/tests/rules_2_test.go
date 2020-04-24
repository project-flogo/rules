package tests

import (
	"context"
	"fmt"
	"testing"

	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

func Test_Two(t *testing.T) {

	// fmt.Println("** rulesapp: Example usage of the Rules module/API **")

	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := common.GetAbsPathForResource("src/github.com/project-flogo/rules/examples/rulesapp/rulesapp.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	// fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	//Create a RuleSession
	rs, _ := ruleapi.GetOrCreateRuleSession("asession")
	actionFireCount := make(map[string]int)

	//// check for name "Bob" in n1
	rule := ruleapi.NewRule("rule1")
	rule.AddCondition("c1", []string{"n1"}, checkForBob, nil)
	rule.AddCondition("c2", []string{"n1"}, checkForName, nil)

	rule.SetAction(checkForBobAction)
	rule.SetContext(actionFireCount)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	//Start the rule session
	rs.Start(nil)
	t1, _ := model.NewTupleWithKeyValues("n1", "Tom")
	rs.Assert(context.TODO(), t1)

	t2, _ := model.NewTupleWithKeyValues("n1", "Bob")
	rs.Assert(context.TODO(), t2)

	//Retract tuples
	rs.Retract(context.TODO(), t1)
	rs.Retract(context.TODO(), t2)

	if cnt, found := actionFireCount["count"]; found {
		if cnt > 1 {
			t.Logf("checkForBobAction fired more than once [%d]", cnt)
			t.FailNow()
		}
	}

	//delete the rule
	rs.DeleteRule(rule.GetName())

	//unregister the session, i.e; cleanup
	rs.Unregister()
}

func checkForName(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//This conditions filters on name="Bob"
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name, _ := t1.GetString("name")
	return len(name) != 0
}

func checkForBob(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//This conditions filters on name="Bob"
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name, _ := t1.GetString("name")
	return name == "Bob"
}

func checkForBobAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return
	}
	name, _ := t1.GetString("name")
	fmt.Println("name=", name)

	actionFiredCount := ruleCtx.(map[string]int)
	if cnt, found := actionFiredCount["count"]; found {
		cnt++
		actionFiredCount["count"] = cnt
	} else {
		actionFiredCount["count"] = 1
	}
}
