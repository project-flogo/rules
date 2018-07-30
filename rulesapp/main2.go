package main

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
)

func main2() {

	fmt.Println("My first BEGo test!")

	//Create Rule, define conditiond and set action callback
	rule := ruleapi.NewRule("first rule")
	rule.AddCondition("c1", []model.TupleTypeAlias{"n1"}, myC11)
	rule.AddCondition("c2", []model.TupleTypeAlias{"n1", "n2", "n3"}, myC31)
	rule.AddCondition("c3", []model.TupleTypeAlias{"n1", "n2"}, myC21)
	rule.SetAction(myActionFn)

	//Create a RuleSession and add the above Rule
	ruleSession := ruleapi.GetOrCreateRuleSession("asession")
	ruleSession.AddRule(rule)

	ctx := context.Background()

	//Now assert a few facts and see if the Rule Action callback fires.
	streamTuple1 := model.NewStreamTuple("n1")
	ruleSession.Assert(ctx, streamTuple1)

	streamTuple2 := model.NewStreamTuple("n1")
	ruleSession.Assert(ctx, streamTuple2)

	streamTuple3 := model.NewStreamTuple("n1")
	ruleSession.Assert(ctx, streamTuple3)

	streamTuple4 := model.NewStreamTuple("n2")
	ruleSession.Assert(ctx, streamTuple4)

	streamTuple5 := model.NewStreamTuple("n2")
	ruleSession.Assert(ctx, streamTuple5)

	streamTuple6 := model.NewStreamTuple("n3")
	ruleSession.Assert(ctx, streamTuple6)

	streamTuple7 := model.NewStreamTuple("n3")
	ruleSession.Assert(ctx, streamTuple7)

	ruleSession.Retract(ctx, streamTuple7)

	fmt.Println("After retracting n3")

	streamTuple8 := model.NewStreamTuple("n1")
	ruleSession.Assert(ctx, streamTuple8)

	ruleSession.DeleteRule("first rule")
	fmt.Println("After deleting 'first rule'")

	ruleSession.Assert(ctx, streamTuple1)
	ruleSession.Assert(ctx, streamTuple4)
	ruleSession.Assert(ctx, streamTuple6)

	// streamTuple8 := model.NewStreamTuple("n2")
	// ruleSession.Assert(ctx, streamTuple8)
	// fmt.Println("Rules fired after retracting n1s...")
	// ruleSession.Retract(streamTuple1)
	// ruleSession.Retract(streamTuple2)
	// ruleSession.Retract(streamTuple3)
	// ruleSession.Retract(streamTuple4)
	// ruleSession.Retract(streamTuple5)

	// ruleSession.Assert(ctx, streamTuple8)

}

func myC11(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	return true
}

func myC21(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	return true
}

func myC31(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	fmt.Printf("Condition [%s] of Rule [%s] has [%d] tuples\n", condName, ruleName, len(tuples))
	return true
}

func myActionFn1(ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	fmt.Printf("My Rule [%s] fired\n", ruleName)
	//for key := range tuples {
	//	fmt.Println("\tmatched tuple entry:" + string(key))
	//}
	//fmt.Printf("Rule [%s] fired and done!\n", ruleName)
}
