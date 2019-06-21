package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/rules/config"

	"github.com/project-flogo/rules/common/model"
)

//add this sample file to your flogo project
func init() {

	config.RegisterStartupRSFunction("simple", StartupRSFunction)

	// rule UserData
	config.RegisterConditionEvaluator("cUserData", cUserData)
	config.RegisterActionFunction("aUserData", aUserData)

	// rule NewUser
	config.RegisterConditionEvaluator("cNewUser", cNewUser)
	config.RegisterActionFunction("aNewUser", aNewUser)

	// rule NewUserApprove
	config.RegisterConditionEvaluator("cNewUserId", cNewUserId)
	config.RegisterConditionEvaluator("cNewUserAge", cNewUserAge)
	config.RegisterConditionEvaluator("cNewUserCibil", cNewUserCibil)
	config.RegisterActionFunction("aNewUserApprove", aNewUserApprove)

	// rule NewUserReject
	config.RegisterConditionEvaluator("cNewUserId", cNewUserId)
	config.RegisterConditionEvaluator("cNewUserAge", cNewUserAge)
	config.RegisterConditionEvaluator("cNewUserLowCibil", cNewUserLowCibil)
	config.RegisterActionFunction("aNewUserReject", aNewUserReject)

	// rule NewUserApprove1
	config.RegisterConditionEvaluator("cNewUserId", cNewUserId)
	config.RegisterConditionEvaluator("cNewUserAge", cNewUserAge)
	config.RegisterConditionEvaluator("cNewUserHighCibil", cNewUserHighCibil)
	config.RegisterActionFunction("aNewUserApprove1", aNewUserApprove1)
}

func StartupRSFunction(ctx context.Context, rs model.RuleSession, startupCtx map[string]interface{}) (err error) {

	fmt.Printf("In startup rule function..\n")
	t3, _ := model.NewTupleWithKeyValues("UserAccount", "Bob")
	t3.SetString(nil, "Name", "Bob")
	rs.Assert(nil, t3)
	return nil
}

func cUserData(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	userInfo := tuples["UserAccount"]
	if userInfo == nil {
		return false
	}
	fmt.Println("cUserData")
	return true
}

func aUserData(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName, "User information recevied")
}

func cNewUser(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	newaccount := tuples["NewAccount"]
	if newaccount != nil {
		address, _ := newaccount.GetString("Address")
		age, _ := newaccount.GetInt("Age")
		income, _ := newaccount.GetInt("Income")
		if address == "" || age < 18 || income < 10000 || age >= 45 {
			return true
		}
	}
	fmt.Println("cNewUser")
	return false
}

func aNewUser(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	newaccount := tuples["NewAccount"]
	rs.Retract(ctx, newaccount)
	fmt.Println("Applicant is not eligible to apply for creditcard")
}

func cNewUserId(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	userInfo := tuples["UserAccount"]
	newaccount := tuples["NewAccount"]
	if newaccount != nil || userInfo != nil {
		userId, _ := userInfo.GetInt("Id")
		newUserId, _ := newaccount.GetInt("Id")
		if userId == newUserId {
			fmt.Println("cNewUserId")
			return true
		}
	}
	return false
}

func cNewUserAge(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	userInfo := tuples["UserAccount"]
	newaccount := tuples["NewAccount"]
	if newaccount != nil || userInfo != nil {
		newUserAge, _ := newaccount.GetInt("Age")
		if newUserAge >= 18 && newUserAge <= 44 {
			fmt.Println("cNewUserAge")
			return true
		}
	}
	return false
}

func cNewUserCibil(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCibil"]
	if updateScore != nil {
		cibil, _ := updateScore.GetInt("creditScore")
		if cibil >= 750 && cibil < 820 {
			fmt.Println("cNewUserCibil")
			return true
		}
	}
	return false
}

func cNewUserLowCibil(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCibil"]
	if updateScore != nil {
		cibil, _ := updateScore.GetInt("creditScore")
		if cibil < 750 {
			fmt.Println("cNewUserLowCibil")
			return true
		}
	}
	return false
}

func cNewUserHighCibil(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCibil"]
	if updateScore != nil {
		cibil, _ := updateScore.GetInt("creditScore")
		if cibil >= 820 && cibil <= 900 {
			fmt.Println("cNewUserHighCibil")
			return true
		}
	}
	return false
}

func aNewUserApprove(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	newaccount := tuples["NewAccount"]
	updateScore := tuples["UpdateCibil"]
	cibil, _ := updateScore.GetInt("creditScore")
	income, _ := newaccount.GetInt("Income")
	var limit = 2 * income
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", cibil)
	userInfoMutable.SetString(ctx, "appStatus", "Approved")
	userInfoMutable.SetInt(ctx, "approvedLimit", limit)
	fmt.Println(userInfo)
}

func aNewUserApprove1(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	newaccount := tuples["NewAccount"]
	updateScore := tuples["UpdateCibil"]
	cibil, _ := updateScore.GetInt("creditScore")
	income, _ := newaccount.GetInt("Income")
	var limit = 3 * income
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", cibil)
	userInfoMutable.SetString(ctx, "appStatus", "Approved")
	userInfoMutable.SetInt(ctx, "approvedLimit", limit)
	fmt.Println(userInfo)
}

func aNewUserReject(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	newaccount := tuples["NewAccount"]
	updateScore := tuples["UpdateCibil"]
	cibil, _ := updateScore.GetInt("creditScore")
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", cibil)
	userInfoMutable.SetString(ctx, "appStatus", "Rejected")
	userInfoMutable.SetInt(ctx, "approvedLimit", 0)
	fmt.Println(userInfo)
	rs.Retract(ctx, newaccount)
}

func cNewUser2(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	userInfo := tuples["UserAccount"]
	newaccount := tuples["NewAccount"]
	if newaccount != nil || userInfo != nil {
		userAge, _ := userInfo.GetInt("Age")
		newUserAge, _ := newaccount.GetInt("Age")
		if userAge >= newUserAge {
			fmt.Println("cNewUser2")
			return true
		}
	}
	return false
}
