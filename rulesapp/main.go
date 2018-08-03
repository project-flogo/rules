package main

import (
	"context"
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
	"os"
	"strings"
	"io/ioutil"
	"log"
)

func main() {

	fmt.Println("** Welcome to BEGo **")

	////A json string describing the types of tuples that will be asserted.
	////The descriptor is used for property type/value validations as well as other information such as expiry of the tuple, etc.
	////Expiry = -1 means explicit retraction. 0 means, retract at the end of RTC, non-zero means retract after that much timeout
	////in milliseconds
	////In a real application, this type descriptor will usually be externalized to a file
	//tupleDescriptor := "[{\"Name\": \"n1\",\"Expiry\": -1,\"Props\": {\"name\": {\"Name\": \"name\",\"PropType\": \"string\"}}}," +
	//	                "{\"Name\": \"n2\",\"Expiry\": -1,\"Props\": {\"name\": {\"Name\": \"name\",\"PropType\": \"string\"}}}]"

	tupleDescriptorFileNm := getAbsPathForResource("src/github.com/TIBCOSoftware/bego/rulesapp/rulesapp.json")

	tupleDescriptor := fileToString(tupleDescriptorFileNm)

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//Create a RuleSession and register the type descriptors.
	rs := ruleapi.GetOrCreateRuleSession("asession")
	rs.RegisterTupleDescriptors(tupleDescriptor)

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
	rs.AddRule(rule)
	rs.AddRule(rule2)
	rs.AddRule(rule3)

	// ctx := context.Background()

	//Now assert a few facts and see if the Rule Action callback fires.
	fmt.Println("Asserting n1 tuple with name=Bob")
	streamTuple1 := model.NewStreamTuple("n1")
	streamTuple1.SetString(nil, rs,"name", "Bob")
	rs.Assert(nil, streamTuple1)

	fmt.Println("Asserting n1 tuple with name=Fred")
	streamTuple2 := model.NewStreamTuple("n1")
	streamTuple2.SetString(nil, rs,"name", "Fred")
	rs.Assert(nil, streamTuple2)

	fmt.Println("Asserting n2 tuple with name=Fred")
	streamTuple3 := model.NewStreamTuple("n2")
	streamTuple3.SetString(nil, rs,"name", "Fred")
	rs.Assert(nil, streamTuple3)

	fmt.Println("Asserting n2 tuple with name=Bob")
	streamTuple4 := model.NewStreamTuple("n2")
	streamTuple4.SetString(nil, rs,"name", "Bob")
	rs.Assert(nil, streamTuple4)

	fmt.Println("Asserting n1 tuple with name=Tom")
	streamTuple5 := model.NewStreamTuple("n1")
	streamTuple5.SetString(nil, rs,"name", "Tom")
	rs.Assert(nil, streamTuple5)

	//Retract them
	rs.Retract(nil, streamTuple1)
	rs.Retract(nil, streamTuple2)
	rs.Retract(nil, streamTuple3)
	rs.Retract(nil, streamTuple4)
	rs.Retract(nil, streamTuple5)

	//You may delete the rule
	rs.DeleteRule(rule.GetName())

	rs.Unregister()

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

func getAbsPathForResource (resourcepath string) string {
	GOPATH := os.Getenv("GOPATH")
	paths := strings.Split(GOPATH, ":")
	for _, path:= range paths {
		absPath := path + "/" + resourcepath
		_, err := os.Stat(absPath)
		if err == nil {
			return absPath
		}
	}
	return ""
}

func fileToString(fileName string)string {

	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(dat)
}