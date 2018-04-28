package main

import (
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
)

func main1() {

	fmt.Println("My first BEGo test!")

	//Create Rule, define conditiond and set action callback
	rule := ruleapi.NewRule("first rule")
	rule.AddCondition("c1", []model.StreamSource{"n1"}, myC11)
	rule.AddCondition("c2", []model.StreamSource{"n1", "n2", "n3"}, myC31)
	rule.AddCondition("c3", []model.StreamSource{"n1", "n2"}, myC21)
	rule.SetActionFn(myActionFn)

	//Create a RuleSession and add the above Rule
	ruleSession := ruleapi.NewRuleSession()
	ruleSession.AddRule(rule)

	//Now assert a few facts and see if the Rule Action callback fires.
	streamTuple1 := model.NewStreamTuple("n1")
	ruleSession.Assert(streamTuple1)

	streamTuple2 := model.NewStreamTuple("n1")
	ruleSession.Assert(streamTuple2)

	streamTuple3 := model.NewStreamTuple("n1")
	ruleSession.Assert(streamTuple3)

	streamTuple4 := model.NewStreamTuple("n2")
	ruleSession.Assert(streamTuple4)

	streamTuple5 := model.NewStreamTuple("n2")
	ruleSession.Assert(streamTuple5)

	streamTuple6 := model.NewStreamTuple("n3")
	ruleSession.Assert(streamTuple6)

	streamTuple7 := model.NewStreamTuple("n3")
	ruleSession.Assert(streamTuple7)

	ruleSession.Retract(streamTuple7)

	fmt.Println("After retracting n3")

	streamTuple8 := model.NewStreamTuple("n1")
	ruleSession.Assert(streamTuple8)

	ruleSession.DeleteRule("first rule")
	fmt.Println("After deleting 'first rule'")

	ruleSession.Assert(streamTuple1)
	ruleSession.Assert(streamTuple4)
	ruleSession.Assert(streamTuple6)

	// streamTuple8 := model.NewStreamTuple("n2")
	// ruleSession.Assert(streamTuple8)
	// fmt.Println("Rules fired after retracting n1s...")
	// ruleSession.Retract(streamTuple1)
	// ruleSession.Retract(streamTuple2)
	// ruleSession.Retract(streamTuple3)
	// ruleSession.Retract(streamTuple4)
	// ruleSession.Retract(streamTuple5)

	// ruleSession.Assert(streamTuple8)

}

func myC11(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {
	fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	return true
}

func myC21(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {
	fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	return true
}

func myC31(ruleName string, condName string, tuples map[model.StreamSource]model.StreamTuple) bool {
	fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	return true
}

func myActionFn1(ruleName string, tuples map[model.StreamSource]model.StreamTuple) {
	fmt.Printf("My Rule [%s] fired\n", ruleName)
	//for key := range tuples {
	//	fmt.Println("\tmatched tuple entry:" + string(key))
	//}
	//fmt.Printf("Rule [%s] fired and done!\n", ruleName)
}
