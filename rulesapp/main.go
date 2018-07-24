package main

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
)

func main() {

	fmt.Println("** Welcome to BEGo **")

	//Create Rule, define conditiond and set action callback
	rule := ruleapi.NewRule("* Ensure n1.name is Bob and n2.name matches n1.name ie Bob in this case *")
	fmt.Printf("Rule added: [%s]\n", rule.GetName())
	rule.AddCondition("c1", []model.TupleTypeAlias{"n1"}, checkForBob)          // check for name "Bob" in n1
	rule.AddCondition("c2", []model.TupleTypeAlias{"n1", "n2"}, checkSameNames) // match the "name" field in both tuples
	//in effect, fire the rule when name field in both tuples in "Bob"
	rule.SetAction(myActionFn)

	//Create Rule, define conditiond and set action callback
	rule2 := ruleapi.NewRule("* name == Tom *")
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())
	rule2.AddCondition("c1", []model.TupleTypeAlias{"n1"}, checkForTom) // check for name "Bob" in n1
	//in effect, fire the rule when name field in both tuples in "Bob"
	rule2.SetAction(checkForTomAction)
	rule2.SetPriority(100)

	//Create Rule, define conditiond and set action callback
	rule3 := ruleapi.NewRule("2* name == Tom *")
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())
	rule3.AddCondition("c1", []model.TupleTypeAlias{"n1"}, checkForTom) // check for name "Bob" in n1
	//in effect, fire the rule when name field in both tuples in "Bob"
	rule3.SetAction(checkForTomAction2)
	rule3.SetPriority(1000)

	//Create a RuleSession and add the above Rule
	ruleSession := ruleapi.NewRuleSession()
	ruleSession.AddRule(rule)
	ruleSession.AddRule(rule2)
	ruleSession.AddRule(rule3)

	// ctx := context.Background()

	//Now assert a few facts and see if the Rule Action callback fires.
	fmt.Println("Asserting n1 tuple with name=Bob")
	streamTuple1 := model.NewStreamTuple("n1")
	streamTuple1.SetString(nil, "name", "Bob")
	ruleSession.Assert(nil, streamTuple1)

	fmt.Println("Asserting n1 tuple with name=Fred")
	streamTuple2 := model.NewStreamTuple("n1")
	streamTuple2.SetString(nil, "name", "Fred")
	ruleSession.Assert(nil, streamTuple2)

	fmt.Println("Asserting n2 tuple with name=Fred")
	streamTuple3 := model.NewStreamTuple("n2")
	streamTuple3.SetString(nil, "name", "Fred")
	ruleSession.Assert(nil, streamTuple3)

	fmt.Println("Asserting n2 tuple with name=Bob")
	streamTuple4 := model.NewStreamTuple("n2")
	streamTuple4.SetString(nil, "name", "Bob")
	ruleSession.Assert(nil, streamTuple4)

	fmt.Println("Asserting n1 tuple with name=Tom")
	streamTuple5 := model.NewStreamTuple("n1")
	streamTuple5.SetString(nil, "name", "Tom")
	ruleSession.Assert(nil, streamTuple5)

	//Retract them
	ruleSession.Retract(nil, streamTuple1)
	ruleSession.Retract(nil, streamTuple2)
	ruleSession.Retract(nil, streamTuple3)
	ruleSession.Retract(nil, streamTuple4)
	ruleSession.Retract(nil, streamTuple5)

	//You may delete the rule
	ruleSession.DeleteRule(rule.GetName())

}

func checkForBob(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	//This conditions filters on name="Bob"
	streamTuple := tuples["n1"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := streamTuple.GetString("name")
	return name == "Bob"
}

func checkForTom(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	//This conditions filters on name="Bob"
	streamTuple := tuples["n1"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := streamTuple.GetString("name")
	return name == "Tom"
}

func checkSameNames(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
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

func myActionFn(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	streamTuple1 := tuples["n1"]
	streamTuple2 := tuples["n2"]
	if streamTuple1 == nil || streamTuple2 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
	}
	name1 := streamTuple1.GetString("name")
	name2 := streamTuple2.GetString("name")
	fmt.Printf("n1.name = [%s], n2.name = [%s]\n", name1, name2)
}

func checkForTomAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	streamTuple1 := tuples["n1"]

	name1 := streamTuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)
}

func checkForTomAction2(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	streamTuple1 := tuples["n1"]

	name1 := streamTuple1.GetString("name")
	fmt.Printf("n1.name = [%s]\n", name1)

	return

}
