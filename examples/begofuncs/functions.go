package begofuncs

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/config"
)

//add this sample file to your flogo project
func init() {
	config.RegisterActionFunction("checkForBobAction", checkForBobAction)
	config.RegisterActionFunction("checkSameNamesAction", checkSameNamesAction)

	config.RegisterConditionEvaluator("checkForBob", checkForBob)
	config.RegisterConditionEvaluator("checkSameNamesCondition", checkSameNamesCondition)
	config.RegisterStartupRSFunction("simple", StartupRSFunction)
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
	fmt.Printf("Context is [%s]\n", ruleCtx)
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return
	}
}

func checkSameNamesCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return false
	}
	name1, _ := t1.GetString("name")
	name2, _ := t2.GetString("name")
	return name1 == name2
}

func checkSameNamesAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
		return
	}
	name1, _ := t1.GetString("name")
	name2, _ := t2.GetString("name")
	fmt.Printf("n1.name = [%s], n2.name = [%s]\n", name1, name2)
}

func StartupRSFunction (ctx context.Context, rs model.RuleSession, startupCtx map[string]interface{}) {

	fmt.Printf("In startup rule function..\n")
	t3, _ := model.NewTupleWithKeyValues("n2", "Bob")
	t3.SetString(nil, "name", "Bob")
	rs.Assert(nil, t3)
}