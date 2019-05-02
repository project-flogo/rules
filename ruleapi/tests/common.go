package tests

import (
	"context"
	"io/ioutil"
	"log"
	"testing"

	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

func createRuleSession() (model.RuleSession, error) {
	rs, _ := ruleapi.GetOrCreateRuleSession("test")

	tupleDescFileAbsPath := common.GetAbsPathForResource("src/github.com/project-flogo/rules/ruleapi/tests/tests.json")

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

//conditions and actions
func trueCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}
func falseCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return false
}
func emptyAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
}

func printTuples(t *testing.T, oprn string, tupleMap map[string]map[string]model.Tuple) {

	for k, v := range tupleMap {
		t.Logf("%s tuples for type [%s]\n", oprn, k)
		for k1, _ := range v {
			t.Logf("    tuples key [%s]\n", k1)
		}
	}
}
func printModified(t *testing.T, modified map[string]map[string]model.RtcModified) {

	for k, v := range modified {
		t.Logf("%s tuples for type [%s]\n", "Modified", k)
		for k1, _ := range v {
			t.Logf("    tuples key [%s]\n", k1)
		}
	}
}

type txnCtx struct {
	Testing *testing.T
	TxnCnt  int
}
