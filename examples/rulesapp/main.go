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
	err := example(true)
	if err != nil {
		panic(err)
	}
}

func example(redis bool) error {
	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := common.GetPathForResource("examples/rulesapp/rulesapp.json", "./rulesapp.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		return err
	}

	//Create a RuleSession
	store := ""
	if redis {
		store = "rsconfig.json"
	}
	rs, err := ruleapi.GetOrCreateRuleSession("asession", store)
	if err != nil {
		return err
	}

	events := make(map[string]int, 8)
	//// check for name "Bob" in n1
	rule := ruleapi.NewRule("n1.name == Bob")
	err = rule.AddCondition("c1", []string{"n1"}, checkForBob, events)
	if err != nil {
		return err
	}
	serviceCfg := &config.ServiceDescriptor{
		Name:     "checkForBobAction",
		Function: checkForBobAction,
		Type:     "function",
	}
	aService, err := ruleapi.NewActionService(serviceCfg)
	if err != nil {
		return err
	}
	rule.SetActionService(aService)
	rule.SetContext(events)
	err = rs.AddRule(rule)
	if err != nil {
		return err
	}

	// check for name "Bob" in n1, match the "name" field in n2,
	// in effect, fire the rule when name field in both tuples is "Bob"
	rule2 := ruleapi.NewRule("n1.name == Bob && n1.name == n2.name")
	err = rule2.AddCondition("c1", []string{"n1"}, checkForBob, events)
	if err != nil {
		return err
	}
	err = rule2.AddCondition("c2", []string{"n1", "n2"}, checkSameNamesCondition, events)
	if err != nil {
		return err
	}
	serviceCfg2 := &config.ServiceDescriptor{
		Name:     "checkSameNamesAction",
		Function: checkSameNamesAction,
		Type:     "function",
	}
	aService2, err := ruleapi.NewActionService(serviceCfg2)
	if err != nil {
		return err
	}
	rule2.SetActionService(aService2)
	rule2.SetContext(events)
	err = rs.AddRule(rule2)
	if err != nil {
		return err
	}

	//set a transaction handler
	rs.RegisterRtcTransactionHandler(txHandler, nil)
	//Start the rule session
	err = rs.Start(nil)
	if err != nil {
		return err
	}

	//Now assert a "n1" tuple
	t1, err := model.NewTupleWithKeyValues("n1", "Tom")
	if err != nil {
		return err
	}
	t1.SetString(nil, "name", "Tom")
	err = rs.Assert(nil, t1)
	if err != nil {
		return err
	}
	t11 := rs.GetStore().GetTupleByKey(t1.GetKey())
	if t11 == nil {
		return fmt.Errorf("Warn: Tuple should be in store[%s]", t11.GetKey())
	}

	//Now assert a "n1" tuple
	t2, err := model.NewTupleWithKeyValues("n1", "Bob")
	if err != nil {
		return err
	}
	t2.SetString(nil, "name", "Bob")
	err = rs.Assert(nil, t2)
	if err != nil {
		return err
	}

	//Now assert a "n2" tuple
	t3, err := model.NewTupleWithKeyValues("n2", "Bob")
	if err != nil {
		return err
	}
	t3.SetString(nil, "name", "Bob")
	err = rs.Assert(nil, t3)
	if err != nil {
		return err
	}

	//Retract tuples
	err = rs.Retract(nil, t1)
	if err != nil {
		return err
	}
	err = rs.Retract(nil, t2)
	if err != nil {
		return err
	}
	err = rs.Retract(nil, t3)
	if err != nil {
		return err
	}

	//delete the rule
	rs.DeleteRule(rule2.GetName())

	//unregister the session, i.e; cleanup
	rs.Unregister()

	if events["checkForBob"] != 4 {
		return fmt.Errorf("checkForBob should have been called 4 times")
	}
	if events["checkForBobAction"] != 1 {
		return fmt.Errorf("checkForBobAction should have been called once")
	}
	if events["checkSameNamesCondition"] != 1 {
		return fmt.Errorf("checkSameNamesCondition should have been called once")
	}
	if events["checkSameNamesAction"] != 1 {
		return fmt.Errorf("checkSameNamesAction should have been called once")
	}

	return nil
}

func checkForBob(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//This conditions filters on name="Bob"
	t1 := tuples["n1"]
	if t1 == nil {
		return false
	}
	name, err := t1.GetString("name")
	if err != nil {
		return false
	}
	if name == "" {
		return false
	}
	events := ctx.(map[string]int)
	count := events["checkForBob"]
	events["checkForBob"] = count + 1
	return name == "Bob"
}

func checkForBobAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples["n1"]
	if t1 == nil {
		return
	}
	name, err := t1.GetString("name")
	if err != nil {
		return
	}
	if name == "" {
		return
	}
	events := ruleCtx.(map[string]int)
	count := events["checkForBobAction"]
	events["checkForBobAction"] = count + 1
}

func checkSameNamesCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		return false
	}
	name1, err := t1.GetString("name")
	if err != nil {
		return false
	}
	if name1 == "" {
		return false
	}
	name2, err := t2.GetString("name")
	if err != nil {
		return false
	}
	if name2 == "" {
		return false
	}
	events := ctx.(map[string]int)
	count := events["checkSameNamesCondition"]
	events["checkSameNamesCondition"] = count + 1
	return name1 == name2
}

func checkSameNamesAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		return
	}
	name1, err := t1.GetString("name")
	if err != nil {
		return
	}
	if name1 == "" {
		return
	}
	name2, err := t2.GetString("name")
	if err != nil {
		return
	}
	if name2 == "" {
		return
	}
	events := ruleCtx.(map[string]int)
	count := events["checkSameNamesAction"]
	events["checkSameNamesAction"] = count + 1
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
