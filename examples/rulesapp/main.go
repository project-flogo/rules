package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/ruleapi"
)

func main() {

	fmt.Println("** rulesapp: Example usage of the Rules module/API **")

	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := common.GetPathForResource("examples/rulesapp/rulesapp.json", "./rulesapp.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		panic(err)
	}

	//Create a RuleSession
	rs, err := ruleapi.GetOrCreateRuleSession("asession")
	if err != nil {
		panic(err)
	}

	//// check for name "Bob" in n1
	rule := ruleapi.NewRule("n1.name == Bob")
	rule.AddCondition("c1", []string{"n1"}, checkForBob, nil)
	serviceCfg := &config.ServiceDescriptor{
		Name:     "checkForBobAction",
		Function: checkForBobAction,
		Type:     "function",
	}
	aService, err := ruleapi.NewActionService(serviceCfg)
	if err != nil {
		panic(err)
	}
	rule.SetActionService(aService)
	rule.SetContext("This is a test of context")
	err = rs.AddRule(rule)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	// check for name "Bob" in n1, match the "name" field in n2,
	// in effect, fire the rule when name field in both tuples is "Bob"
	rule2 := ruleapi.NewRule("n1.name == Bob && n1.name == n2.name")
	rule2.AddCondition("c1", []string{"n1"}, checkForBob, nil)
	rule2.AddCondition("c2", []string{"n1", "n2"}, checkSameNamesCondition, nil)
	serviceCfg2 := &config.ServiceDescriptor{
		Name:     "checkSameNamesAction",
		Function: checkSameNamesAction,
		Type:     "function",
	}
	aService2, _ := ruleapi.NewActionService(serviceCfg2)
	rule2.SetActionService(aService2)
	err = rs.AddRule(rule2)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//set a transaction handler
	rs.RegisterRtcTransactionHandler(txHandler, nil)
	//Start the rule session
	err = rs.Start(nil)
	if err != nil {
		panic(err)
	}

	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Tom")
	t1, _ := model.NewTupleWithKeyValues("n1", "Tom")
	t1.SetString(nil, "name", "Tom")
	err = rs.Assert(nil, t1)
	if err != nil {
		panic(err)
	}
	t11 := rs.GetStore().GetTupleByKey(t1.GetKey())
	if t11 == nil {
		panic(fmt.Errorf("Warn: Tuple should be in store[%s]", t11.GetKey()))
	}

	//Now assert a "n1" tuple
	fmt.Println("Asserting n1 tuple with name=Bob")
	t2, _ := model.NewTupleWithKeyValues("n1", "Bob")
	t2.SetString(nil, "name", "Bob")
	err = rs.Assert(nil, t2)
	if err != nil {
		panic(err)
	}

	//Now assert a "n2" tuple
	fmt.Println("Asserting n2 tuple with name=Bob")
	t3, _ := model.NewTupleWithKeyValues("n2", "Bob")
	t3.SetString(nil, "name", "Bob")
	err = rs.Assert(nil, t3)
	if err != nil {
		panic(err)
	}

	//Retract tuples
	err = rs.Retract(nil, t1)
	if err != nil {
		panic(err)
	}
	err = rs.Retract(nil, t2)
	if err != nil {
		panic(err)
	}
	err = rs.Retract(nil, t3)
	if err != nil {
		panic(err)
	}

	//delete the rule
	rs.DeleteRule(rule2.GetName())

	//unregister the session, i.e; cleanup
	rs.Unregister()

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
		fmt.Println("Should not get nil tuples here in JoinCondition1! This is an error")
		return
	}
}

func checkSameNamesCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition2! This is an error")
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

func getFileContent(filePath string) string {
	absPath := common.GetAbsPathForResource(filePath)
	return common.FileToString(absPath)
}

func txHandler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	store := rs.GetStore()
	store.SaveTuples(rtxn.GetRtcAdded())

	store.SaveModifiedTuples(rtxn.GetRtcModified())

	store.DeleteTuples(rtxn.GetRtcDeleted())

}
