package main

import (
	"strings"
	"context"
	"fmt"
	"github.com/project-flogo/rules/config"

	"github.com/project-flogo/rules/common/model"
)

//add this sample file to your flogo project
func init() {
	config.RegisterActionFunction("checkForApplicationdata", checkForApplicationdata)
	config.RegisterActionFunction("checkForApplicationStatus", checkForApplicationStatus)
	config.RegisterConditionEvaluator("checkForIdMatch", checkForIdMatch)
	config.RegisterConditionEvaluator("checkForNameMatch", checkForNameMatch)
	config.RegisterConditionEvaluator("checkForAddress", checkForAddress)
	config.RegisterConditionEvaluator("checkForAge", checkForAge)
	config.RegisterConditionEvaluator("checkForCibil", checkForCibil)
	config.RegisterConditionEvaluator("checkForEligibleCreditlimit", checkForEligibleCreditlimit)
	config.RegisterStartupRSFunction("simple", StartupRSFunction)
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
	fmt.Printf("Credit card application approved")
}

func StartupRSFunction(ctx context.Context, rs model.RuleSession, startupCtx map[string]interface{}) (err error) {

	fmt.Printf("In startup rule function..\n")
	t3, _ := model.NewTupleWithKeyValues("n1", "Bob")
	t3.SetString(nil, "name", "Bob")
	rs.Assert(nil, t3)
	return nil
}
