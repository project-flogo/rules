package rete

import (
	"fmt"
	"testing"

	"context"

	"github.com/TIBCOSoftware/bego/common/model"
)

func TestNetworkSimple(t *testing.T) {

	fmt.Println("Creating a RETE network..")

	network := NewReteNetwork()

	// r1 := createR1("r1")
	// fmt.Println(r1)

	r2 := createR2("r2")
	// fmt.Println(r2)

	// network.AddRule(r1)
	network.AddRule(r2)

	// network.RemoveRule(r1.GetName())

	str := network.String()
	fmt.Println("RETE network.." + str)

	// network.AddRule(r1)

	streamTuple1 := model.NewTuple("n1")
	network.Assert(nil, nil, streamTuple1)

	streamTuple2 := model.NewTuple("n1")
	network.Assert(nil, nil, streamTuple2)

	streamTuple3 := model.NewTuple("n1")
	network.Assert(nil, nil, streamTuple3)

	streamTuple4 := model.NewTuple("n1")
	network.Assert(nil, nil, streamTuple4)

	streamTuple5 := model.NewTuple("n1")
	network.Assert(nil, nil, streamTuple5)

	// streamTuple6 := model.NewTuple("n6")
	// network.Assert(streamTuple6)

	// streamTuple7 := model.NewTuple("n7")
	// network.Assert(streamTuple7)

	streamTuple8 := model.NewTuple("n2")
	network.Assert(nil, nil, streamTuple8)
	fmt.Println("Rules fired after retracting n1s...")
	network.Retract(streamTuple1)
	network.Retract(streamTuple2)
	network.Retract(streamTuple3)
	network.Retract(streamTuple4)
	network.Retract(streamTuple5)

	network.Assert(nil, nil, streamTuple8)

	network.RemoveRule(r2.GetName())
}

func c1(conditionName string, ruleName string, tupleMap map[model.TupleType]model.Tuple) bool {
	// fmt.Printf("evaluating condition [%s] for rule [%s]\n", conditionName, ruleName)
	// for key := range tupleMap {
	// 	fmt.Println("\tcond. eval" + string(key))
	// }
	// fmt.Printf("evaluating condition [%s] for rule [%s] done!\n", conditionName, ruleName)
	return true
}

func createR1(name string) model.Rule {
	rule := NewRule(name)
	rule.AddCondition("c1", []model.TupleType{"n1"}, c1)
	rule.AddCondition("c2", []model.TupleType{"n2"}, c1)
	rule.AddCondition("c3", []model.TupleType{"n1", "n2"}, c1)
	rule.AddCondition("c4", []model.TupleType{"n4", "n3"}, c1)
	rule.AddCondition("c6", []model.TupleType{"n5", "n6"}, c1)
	rule.AddCondition("c5", []model.TupleType{"n1", "n2", "n3"}, c1)
	rule.AddCondition("c6", []model.TupleType{"n4", "n5", "n6"}, c1)
	rule.SetAction(r1Action)
	fmt.Println(rule)
	return rule
}

func createR2(name string) model.Rule {
	rule := NewRule(name)
	rule.AddCondition("c1", []model.TupleType{"n1"}, c1)
	rule.AddCondition("c2", []model.TupleType{"n2"}, c1)
	// rule.AddCondition("c3", []string{"n1", "n2")
	// rule.AddCondition("c4", []string{"n1", "n3")
	// rule.AddCondition("c6", []string{"n5", "n6")
	// rule.AddCondition("c5", []string{"n1", "n2", "n3")
	// rule.AddCondition("c6", []string{"n4", "n5", "n6")
	rule.SetAction(r2Action)
	fmt.Println(rule)
	return rule
}

func r1Action(context context.Context, ruleSession model.RuleSession, ruleName string, tupleMap map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule [%s] fired", ruleName)
	for key := range tupleMap {
		fmt.Println("\tmatched tuple entry:" + string(key))
	}
	fmt.Printf("Rule [%s] fired and done!\n", ruleName)
}

func r2Action(context context.Context, ruleSession model.RuleSession, ruleName string, tupleMap map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule [%s] fired\n", ruleName)
	for key := range tupleMap {
		fmt.Println("\tmatched tuple entry:" + string(key))
	}
	fmt.Printf("Rule [%s] fired and done!\n", ruleName)
}
