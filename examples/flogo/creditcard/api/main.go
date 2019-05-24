package main

import (
	"context"
	"fmt"
	"strings"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"github.com/project-flogo/rules/common"
)

func main() {

	fmt.Println("** rulesapp: Example usage of the Rules module/API **")

	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := common.GetAbsPathForResource("src/github.com/project-flogo/rules/examples/flogo/creditcard/api/card.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	//Create a RuleSession
	rs, _ := ruleapi.GetOrCreateRuleSession("asession")

	//// Input applicant data
	rule := ruleapi.NewRule("Applicant data")
	rule.AddCondition("c1", []string{"n1"}, checkForAddress, nil)
	rule.AddCondition("c1", []string{"n1"}, checkForAge, nil)
	rule.SetAction(checkForApplicationdata)
	rule.SetContext("This is application status context")
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	// Input cibil data of applicant
	rule2 := ruleapi.NewRule("Application status")
	rule2.AddCondition("c2", []string{"n1", "n2"}, checkForIdMatch, nil)
	rule2.AddCondition("c2", []string{"n1", "n2"}, checkForNameMatch, nil)
	rule2.AddCondition("c2", []string{"n2"}, checkForCibil, nil)
	rule2.AddCondition("c2", []string{"n1", "n2"}, checkForEligibleCreditlimit, nil)
	rule2.SetAction(checkForApplicationStatus)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//Start the rule session
	rs.Start(nil)

	//Now assert a "n1" tuple
	fmt.Println("Input user data")
	t1, _ := model.NewTupleWithKeyValues("n1", "Tom", 20, "SFO", 1, 50000)
	t1.SetString(nil, "name", "Tom")
	t1.SetInt(nil, "age", 20)
	t1.SetString(nil, "address", "SFO")
	t1.SetInt(nil, "id", 1)
	t1.SetInt(nil, "salary", 51000)
	rs.Assert(nil, t1)

	//Now assert a "n2" tuple
	fmt.Println("\nInput cibil data")
	t2, _ := model.NewTupleWithKeyValues("n2", "Tom", 751, 1)
	t2.SetString(nil, "name", "Bob")
	t2.SetInt(nil, "creditscore", 751)
	t2.SetInt(nil, "id", 1)
	rs.Assert(nil, t2)

	//Retract tuples
	rs.Retract(nil, t1)
	rs.Retract(nil, t2)
	// rs.Retract(nil, t3)

	//delete the rule
	rs.DeleteRule(rule.GetName())

	//unregister the session, i.e; cleanup
	rs.Unregister()
}

func checkForAddress(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	// This condition checks if address is empty or not
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	address, _ := t1.GetString("address")
	fmt.Println("Address of the applicant is", address)
	return true;
}

func checkForAge(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	// This condition matches on age of applicant
	t1 := tuples["n1"]
	if t1 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	age, _ := t1.GetInt("age")
	fmt.Println("Age of the applicant is", age)
	return age >= 18;
}

func checkForIdMatch(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	// This conditions filters on id from applicant and cibildata
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil  || t2 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	applicantid, _ := t1.GetInt("id")
	id, _ := t2.GetInt("id")
	if applicantid == id {
		fmt.Println("Applicant id match found with id", applicantid)
		return true;	
	}
	return false
}
func checkForNameMatch(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//This conditions filters on name from applicant and cibildata"
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil  || t2 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	applicantname, _ := t1.GetString("name")
	name, _ := t2.GetString("name")
	if strings.Compare(applicantname, name) == 0{
		fmt.Println("Name match found")
		return true
	}
	return false
}

func checkForCibil(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	//This conditions obtains cibil data"
	t2 := tuples["n2"]
	if t2 == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	creditScore, _ := t2.GetInt("creditScore")
	if creditScore >= 750{
		fmt.Println("Creditscore of the applicant is", creditScore)
		return true;
	}
	fmt.Println("Application cannot be processed with lower cibil score", creditScore)
	return false;
}

func checkForEligibleCreditlimit(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	// This condition matches for creditcard eligibility criteria
	t1 := tuples["n1"]
	t2 := tuples["n2"]
	if t1 == nil || t2 == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return true;
	}
	cibil, _ := t2.GetInt("creditScore")
	age, _ := t1.GetInt("age");
	sal, _ := t1.GetInt("salary");
	if age >= 18 {
		if cibil >= 750 && sal >= 20000 && sal < 50000 {
			var approvedlimit = 2500
			fmt.Println("Credit card applicant eligile with credit limit",approvedlimit);
			return true;
		}else if cibil >= 750 && sal >= 50000 {
			var approvedlimit = 2*sal
			fmt.Println("Credit card applicant eligile with credit limit",approvedlimit);
			return true;
		}else if cibil >= 750 && sal < 20000 {
			fmt.Println("Credit card application cannot be processed");
			return true;
		}
	}	
	return false;
}

func checkForApplicationdata(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	fmt.Printf("Application submitted")
}

func checkForApplicationStatus(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)
	fmt.Printf("Credit card application approved\n")
}
