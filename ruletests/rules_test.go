package rulesapp

import (
	"context"
	"fmt"
	"testing"

	"github.com/tibmatt/bego/common/model"
	"github.com/tibmatt/bego/ruleapi"
)

func TestOne(t *testing.T) {

	fmt.Println("** Welcome to BEGo **")

	//Create Rule, define conditiond and set action callback
	rule := ruleapi.NewRule("* Ensure n1.name is Bob*")
	fmt.Printf("Rule added: [%s]\n", rule.GetName())
	rule.AddCondition("c1", []model.TupleType{"n1"}, checkForBob) // check for name "Bob" in n1
	rule.SetAction(bobRuleFired)

	//Create Rule, define conditiond and set action callback
	rule2 := ruleapi.NewRule("* name == Tom *")
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())
	rule2.AddCondition("c1", []model.TupleType{"n1"}, checkForTom) // check for name "Bob" in n1
	rule2.SetAction(tomRuleFired)

	//Create a RuleSession and add the above Rule
	ruleSession := ruleapi.GetOrCreateRuleSession("testsession")
	ruleSession.AddRule(rule)
	ruleSession.AddRule(rule2)

	//Now assert a few facts and see if the Rule Action callback fires.
	fmt.Println("Asserting n1 tuple with name=Bob")
	tuple1 := model.NewTuple("n1")
	tuple1.SetString(nil, "name", "Bob")
	ruleSession.Assert(nil, tuple1)

	//Retract them
	ruleSession.Retract(nil, tuple1)

	//You may delete the rule
	ruleSession.DeleteRule(rule.GetName())

}

func checkForBob(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	//This conditions filters on name="Bob"
	tuple := tuples["n1"]
	if tuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := tuple.GetString("name")
	return name == "Bob"
}

func checkForTom(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	//This conditions filters on name="Bob"
	tuple := tuples["n1"]
	if tuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := tuple.GetString("name")
	return name == "Tom"
}

func checkSameNames(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	// fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	tuple1 := tuples["n1"]
	tuple2 := tuples["n2"]
	if tuple1 == nil || tuple2 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return false
	}
	name1 := tuple1.GetString("name")
	name2 := tuple2.GetString("name")
	return name1 == name2
}

func myActionFn(ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule [%s] fired\n", ruleName)
	tuple1 := tuples["n1"]
	tuple2 := tuples["n2"]
	if tuple1 == nil || tuple2 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
	}
	name1 := tuple1.GetString("name")
	name2 := tuple2.GetString("name")
	fmt.Printf("n1.name = [%s], n2.name = [%s]\n", name1, name2)
}

func bobRuleFired(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule [%s] fired\n", ruleName)
	tuple1 := tuples["n1"].(model.MutableTuple)
	if tuple1 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
	}
	name1 := tuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)
	// assertTom(ctx, rs)
	//tuple1.SetInt(ctx, "age", 36)
}

func tomRuleFired(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Tom Rule fired: [%s]\n", ruleName)
	tuple1 := tuples["n1"]

	name1 := tuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)
}

func checkForTomAction2(ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	tuple1 := tuples["n1"]

	name1 := tuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)

	return

}

func assertTom(ctx context.Context, rs model.RuleSession) {
	fmt.Println("Asserting n1 tuple with name=Tom")
	tuple1 := model.NewTuple("n1")
	tuple1.SetString(ctx, "name", "Tom")
	rs.Assert(ctx, tuple1)
}
