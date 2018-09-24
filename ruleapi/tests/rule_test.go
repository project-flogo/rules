package trackntrace

import (
	"fmt"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"github.com/project-flogo/rules/common"

	"testing"
	"io/ioutil"
	"log"
	"context"
)

func Test_T1_AssertTTLZero(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R1")
	rule.AddCondition("R1_c1", []string{"t1.none"}, truecondition, nil)
	rule.SetAction(R1_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.RegisterRtcTransactionHandler(T1Handler, t)
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t1", "t1")
	rs.Assert(context.TODO(), t1)
}

func Test_T2_AssertTTLNegative(t *testing.T) {

	rs, _ := createRuleSession()

	rule := ruleapi.NewRule("R1")
	rule.AddCondition("R1_c1", []string{"t2.none"}, truecondition, nil)
	rule.SetAction(R1_action)
	rule.SetPriority(1)
	rs.AddRule(rule)
	t.Logf("Rule added: [%s]\n", rule.GetName())

	rs.RegisterRtcTransactionHandler(T2Handler, t)
	rs.Start(nil)

	t1, _ := model.NewTupleWithKeyValues("t2", "t2")
	rs.Assert(context.TODO(), t1)
}


//conditions and actions
func truecondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}

func R1_action(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
}

func T1Handler (ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	t := handlerCtx.(*testing.T)

	lA := len(rtxn.GetRtcAdded())
	if  lA != 1 {
		t.Errorf ("RtcAdded: Expected [%d], got [%d]\n", 1, lA)
	}
	lM := len(rtxn.GetRtcModified())
	if  lM != 0 {
		t.Errorf ("RtcModified: Expected [%d], got [%d]\n", 0, lM)
	}
	lD := len(rtxn.GetRtcDeleted())
	if  lD != 0 {
		t.Errorf ("RtcDeleted: Expected [%d], got [%d]\n", 0, lD)
	}
}

func T2Handler (ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	t := handlerCtx.(*testing.T)

	lA := len(rtxn.GetRtcAdded())
	if  lA != 0 {
		t.Errorf ("RtcAdded: Expected [%d], got [%d]\n", 0, lA)
	}
	lM := len(rtxn.GetRtcModified())
	if  lM != 0 {
		t.Errorf ("RtcModified: Expected [%d], got [%d]\n", 0, lM)
	}
	lD := len(rtxn.GetRtcDeleted())
	if  lD != 0 {
		t.Errorf ("RtcDeleted: Expected [%d], got [%d]\n", 0, lD)
	}
}

func printTuples (oprn string, tupleMap map[string]map[string]model.Tuple) {

	for k, v := range tupleMap {
		fmt.Printf ("%s tuples for type [%s]\n", oprn, k)
		for k1, _ := range v {
			fmt.Printf("    tuples key [%s]\n", k1)
		}
	}
}
func printModified (modified map[string]map[string]model.RtcModified) {

	for k, v := range modified {
		fmt.Printf ("%s tuples for type [%s]\n", "Modified", k)
		for k1, _ := range v {
			fmt.Printf("    tuples key [%s]\n", k1)
		}
	}
}

func createRuleSession() (model.RuleSession, error) {
	rs, _ := ruleapi.GetOrCreateRuleSession("test-session")

	tupleDescFileAbsPath := common.GetAbsPathForResource("src/github.com/project-flogo/rules/ruleapi/tests/rule_test.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = model.RegisterTupleDescriptors(string(dat))
	if err != nil {
		return nil, err
	}
	return rs, nil
}



