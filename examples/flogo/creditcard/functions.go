package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/rules/config"

	"github.com/project-flogo/rules/common/model"
)

//add this sample file to your flogo project
func init() {

	// rule UserData
	config.RegisterConditionEvaluator("cNewUser", cNewUser)
	config.RegisterActionFunction("aNewUser", aNewUser)

	// rule NewUser
	config.RegisterConditionEvaluator("cBadUser", cBadUser)
	config.RegisterActionFunction("aBadUser", aBadUser)

	// rule NewUserApprove
	config.RegisterConditionEvaluator("cUserIdMatch", cUserIdMatch)
	config.RegisterConditionEvaluator("cUserCibil", cUserCibil)
	config.RegisterActionFunction("aApproveWithLowLimit", aApproveWithLowLimit)

	// // rule NewUserReject
	config.RegisterConditionEvaluator("cUserIdMatch", cUserIdMatch)
	config.RegisterConditionEvaluator("cUserLowCibil", cUserLowCibil)
	config.RegisterActionFunction("aUserReject", aUserReject)

	// // rule NewUserApprove1
	config.RegisterConditionEvaluator("cUserIdMatch", cUserIdMatch)
	config.RegisterConditionEvaluator("cUserHighCibil", cUserHighCibil)
	config.RegisterActionFunction("aApproveWithHigherLimit", aApproveWithHigherLimit)
}

func cNewUser(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	newaccount := tuples["NewAccount"]
	if newaccount != nil {
		address, _ := newaccount.GetString("Address")
		age, _ := newaccount.GetInt("Age")
		income, _ := newaccount.GetInt("Income")
		if address != "" || age >= 18 || income >= 10000 || age <= 44 {
			return true
		}
	}
	return false
}

func aNewUser(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	newaccount := tuples["NewAccount"]
	Id, _ := newaccount.GetInt("Id")
	name, _ := newaccount.GetString("Name")
	address, _ := newaccount.GetString("Address")
	age, _ := newaccount.GetInt("Age")
	income, _ := newaccount.GetInt("Income")
	userInfo, _ := model.NewTupleWithKeyValues("UserAccount", Id)
	userInfo.SetString(ctx, "Name", name)
	userInfo.SetString(ctx, "Addresss", address)
	userInfo.SetInt(ctx, "Age", age)
	userInfo.SetInt(ctx, "Income", income)
	fmt.Println(userInfo)
	rs.Assert(ctx, userInfo)
	fmt.Println("User information received")

}

func cBadUser(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	newaccount := tuples["NewAccount"]
	if newaccount != nil {
		address, _ := newaccount.GetString("Address")
		age, _ := newaccount.GetInt("Age")
		income, _ := newaccount.GetInt("Income")
		if address == "" || age < 18 || income < 10000 || age >= 45 {
			return true
		}
	}
	return false
}

func aBadUser(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	newaccount := tuples["NewAccount"]
	// rs.Retract(ctx, newaccount)
	fmt.Println(newaccount)
	fmt.Println("Applicant is not eligible to apply for creditcard")
}

func cUserIdMatch(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCibil"]
	if userInfo != nil || updateScore != nil {
		userId, _ := userInfo.GetInt("Id")
		newUserId, _ := updateScore.GetInt("Id")
		if userId == newUserId {
			fmt.Println("Userid match found")
			return true
		}
	}
	return false
}

func cUserCibil(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCibil"]
	if updateScore != nil {
		cibil, _ := updateScore.GetInt("creditScore")
		if cibil >= 750 && cibil < 820 {
			fmt.Println("cUserCibil")
			return true
		}
	}
	return false
}

func cUserLowCibil(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCibil"]
	if updateScore != nil {
		cibil, _ := updateScore.GetInt("creditScore")
		if cibil < 750 {
			fmt.Println("cUserLowCibil")
			return true
		}
	}
	return false
}

func cUserHighCibil(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	updateScore := tuples["UpdateCibil"]
	if updateScore != nil {
		cibil, _ := updateScore.GetInt("creditScore")
		if cibil >= 820 && cibil <= 900 {
			fmt.Println("cUserHighCibil")
			return true
		}
	}
	return false
}

func aApproveWithLowLimit(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCibil"]
	cibil, _ := updateScore.GetInt("creditScore")
	income, _ := userInfo.GetInt("Income")
	var limit = 2 * income
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", cibil)
	userInfoMutable.SetString(ctx, "appStatus", "Approved")
	userInfoMutable.SetInt(ctx, "approvedLimit", limit)
	fmt.Println(userInfo)
}

func aApproveWithHigherLimit(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCibil"]
	cibil, _ := updateScore.GetInt("creditScore")
	income, _ := userInfo.GetInt("Income")
	var limit = 3 * income
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", cibil)
	userInfoMutable.SetString(ctx, "appStatus", "Approved")
	userInfoMutable.SetInt(ctx, "approvedLimit", limit)
	fmt.Println(userInfo)
}

func aUserReject(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Println("Rule fired:", ruleName)
	userInfo := tuples["UserAccount"]
	updateScore := tuples["UpdateCibil"]
	cibil, _ := updateScore.GetInt("creditScore")
	userInfoMutable := userInfo.(model.MutableTuple)
	userInfoMutable.SetInt(ctx, "creditScore", cibil)
	userInfoMutable.SetString(ctx, "appStatus", "Rejected")
	userInfoMutable.SetInt(ctx, "approvedLimit", 0)
	fmt.Println(userInfo)
}
