package rulesapp

import (
	"context"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
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
	streamTuple1 := model.NewTuple("n1")
	streamTuple1.SetString(nil, ruleSession, "name", "Bob")
	ruleSession.Assert(nil, streamTuple1)

	//Retract them
	ruleSession.Retract(nil, streamTuple1)

	//You may delete the rule
	ruleSession.DeleteRule(rule.GetName())

}

func checkForBob(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	//This conditions filters on name="Bob"
	streamTuple := tuples["n1"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := streamTuple.GetString("name")
	return name == "Bob"
}

func checkForTom(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	//This conditions filters on name="Bob"
	streamTuple := tuples["n1"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := streamTuple.GetString("name")
	return name == "Tom"
}

func checkSameNames(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	// fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	streamTuple1 := tuples["n1"]
	streamTuple2 := tuples["n2"]
	if streamTuple1 == nil || streamTuple2 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return false
	}
	name1 := streamTuple1.GetString("name")
	name2 := streamTuple2.GetString("name")
	return name1 == name2
}

func myActionFn(ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule [%s] fired\n", ruleName)
	streamTuple1 := tuples["n1"]
	streamTuple2 := tuples["n2"]
	if streamTuple1 == nil || streamTuple2 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
	}
	name1 := streamTuple1.GetString("name")
	name2 := streamTuple2.GetString("name")
	fmt.Printf("n1.name = [%s], n2.name = [%s]\n", name1, name2)
}

func bobRuleFired(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule [%s] fired\n", ruleName)
	streamTuple1 := tuples["n1"].(model.MutableStreamTuple)
	if streamTuple1 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
	}
	name1 := streamTuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)
	// assertTom(ctx, rs)
	//streamTuple1.SetInt(ctx, "age", 36)
}

func tomRuleFired(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Tom Rule fired: [%s]\n", ruleName)
	streamTuple1 := tuples["n1"]

	name1 := streamTuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)
}

func checkForTomAction2(ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	streamTuple1 := tuples["n1"]

	name1 := streamTuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)

	return

}

func assertTom(ctx context.Context, rs model.RuleSession) {
	fmt.Println("Asserting n1 tuple with name=Tom")
	streamTuple5 := model.NewTuple("n1")
	streamTuple5.SetString(ctx, rs, "name", "Tom")
	rs.Assert(ctx, streamTuple5)
}
