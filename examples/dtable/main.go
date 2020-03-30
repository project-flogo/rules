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
	tupleDescAbsFileNm := common.GetPathForResource("examples/dtable/rulesapp.json", "./rulesapp.json")
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

	// student care information rule
	rule1 := ruleapi.NewRule("studentcare")
	err = rule1.AddExprCondition("c2", "$.student.careRequired", nil)
	if err != nil {
		return err
	}
	printService := &config.ServiceDescriptor{
		Name:     "printstudentinfo",
		Function: printStudentInfo,
		Type:     "function",
	}
	aService1, err := ruleapi.NewActionService(printService)
	if err != nil {
		return err
	}
	rule1.SetActionService(aService1)
	rule1.SetPriority(2)

	err = rs.AddRule(rule1)
	if err != nil {
		return err
	}

	// student analysis rule
	rule2 := ruleapi.NewRule("studentanalysis")
	err = rule2.AddExprCondition("c2", "$.studentanalysis.name == $.student.name", nil)
	if err != nil {
		return err
	}

	settings := make(map[string]interface{})
	settings["filename"] = "dtable-file.xlsx"

	dtableService := &config.ServiceDescriptor{
		Name:     "dtableservice",
		Type:     "decisiontable",
		Settings: settings,
	}
	aService2, err := ruleapi.NewActionService(dtableService)
	if err != nil {
		return err
	}
	rule2.SetActionService(aService2)
	rule2.SetPriority(1)

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

	// assert student info
	s1, err := model.NewTupleWithKeyValues("student", "s1")
	if err != nil {
		return err
	}
	s1.SetString(nil, "grade", "GRADE-C")
	s1.SetString(nil, "class", "X-A")
	s1.SetBool(nil, "careRequired", false)
	err = rs.Assert(nil, s1)
	if err != nil {
		return err
	}

	// assert another student info
	s2, err := model.NewTupleWithKeyValues("student", "s2")
	if err != nil {
		return err
	}
	s2.SetString(nil, "grade", "GRADE-B")
	s2.SetString(nil, "class", "X-A")
	s2.SetBool(nil, "careRequired", false)
	err = rs.Assert(nil, s2)
	if err != nil {
		return err
	}

	// assert studentanalysis event
	se1, err := model.NewTupleWithKeyValues("studentanalysis", "s1")
	if err != nil {
		return err
	}
	err = rs.Assert(nil, se1)
	if err != nil {
		return err
	}

	// assert studentanalysis event
	se2, err := model.NewTupleWithKeyValues("studentanalysis", "s2")
	if err != nil {
		return err
	}
	err = rs.Assert(nil, se2)
	if err != nil {
		return err
	}

	//unregister the session, i.e; cleanup
	rs.Unregister()

	return nil
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

func printStudentInfo(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	student := tuples["student"].ToMap()
	fmt.Println("Student Name: ", student["name"], " Comments: ", student["comments"])
}
