package main

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	fmt.Println("** Welcome to BEGo **")

	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := getAbsPathForResource("src/github.com/TIBCOSoftware/bego/rulesapp/rulesapp.json")
	tupleDescriptor := fileToString(tupleDescAbsFileNm)

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//Create a RuleSession and register the type descriptors.
	rs := ruleapi.GetOrCreateRuleSession("asession")
	model.RegisterTupleDescriptors(tupleDescriptor)

	//// check for name "Bob" in n1
	rule := ruleapi.NewRule("n1.name == Bob")
	rule.AddCondition("c1", []model.TupleType{"n1"}, checkForBob)
	rule.SetAction(checkForBobAction)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	// check for name "Bob" in n1, match the "name" field in n2,
	// in effect, fire the rule when name field in both tuples in "Bob"
	rule2 := ruleapi.NewRule("n1.name == Bob && n1.name == n2.name")
	rule2.AddCondition("c1", []model.TupleType{"n1"}, checkForBob)
	rule2.AddCondition("c2", []model.TupleType{"n1", "n2"}, checkSameNamesCondition)
	rule2.SetAction(checkSameNamesAction)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//Now assert a few facts and see if the Rule Action callback fires.
	fmt.Println("Asserting n1 tuple with name=Tom")
	t1 := model.NewTuple("n1")
	t1.SetString(nil, "name", "Tom")
	rs.Assert(nil, t1)

	//Now assert a few facts and see if the Rule Action callback fires.
	fmt.Println("Asserting n1 tuple with name=Bob")
	t2 := model.NewTuple("n1")
	t2.SetString(nil, "name", "Bob")
	rs.Assert(nil, t2)

	fmt.Println("Asserting n2 tuple with name=Bob")
	t3 := model.NewTuple("n2")
	t3.SetString(nil, "name", "Bob")
	rs.Assert(nil, t3)

	//Retract them
	rs.Retract(nil, t1)
	rs.Retract(nil, t2)
	rs.Retract(nil, t3)

	//delete the rule
	rs.DeleteRule(rule.GetName())

	//unregister the session, meaning, cleanup/remove all data associated with this session
	rs.Unregister()

}

func checkForBob(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	//This conditions filters on name="Bob"
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	name := t1.GetString("name")
	return name == "Bob"
}

func checkForBobAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return
	}
}

func checkSameNamesCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple) bool {
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return false
	}
	name1 := t1.GetString("name")
	name2 := t2.GetString("name")
	return name1 == name2
}

func checkSameNamesAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
	}
	name1 := t1.GetString("name")
	name2 := t2.GetString("name")
	fmt.Printf("n1.name = [%s], n2.name = [%s]\n", name1, name2)
}

func getAbsPathForResource(resourcepath string) string {
	GOPATH := os.Getenv("GOPATH")
	regex, err := regexp.Compile(":|;")
	if err != nil {
		return ""
	}
	paths := regex.Split(GOPATH, -1)
	if os.PathListSeparator == ';' {
		//windows
		resourcepath = strings.Replace(resourcepath, "/", string(os.PathSeparator), -1)
	}
	for _, path := range paths {
		absPath := path + string(os.PathSeparator) + resourcepath
		_, err := os.Stat(absPath)
		if err == nil {
			return absPath
		}
	}
	return ""
}

func fileToString(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(dat)
}
