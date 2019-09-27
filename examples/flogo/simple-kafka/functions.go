package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/rules/config"

	"github.com/project-flogo/rules/common/model"
)

func init() {
	config.RegisterConditionEvaluator("checkForGrocery", checkForGrocery)
	config.RegisterActionFunction("groceryAction", groceryAction)

	config.RegisterConditionEvaluator("checkForFurniture", checkForFurniture)
	config.RegisterActionFunction("furnitureAction", furnitureAction)

}

func checkForGrocery(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	ordr := tuples["order"]
	if ordr == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	typeVal, _ := ordr.GetString("type")
	return typeVal == "grocery"
}

func groceryAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)

	t1 := tuples["order"]
	if t1 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
		return
	}

	price, _ := t1.GetDouble("totalPrice")
	if price >= 2000 {
		fmt.Println("Congratulations you are eligible for Rs. 500 gift coupon")
	}

}

func checkForFurniture(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	ordr := tuples["order"]
	if ordr == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
		return false
	}
	typeVal, _ := ordr.GetString("type")
	return typeVal == "furniture"
}

func furnitureAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	fmt.Printf("Rule fired: [%s]\n", ruleName)

	t1 := tuples["order"]
	if t1 == nil {
		fmt.Println("Should not get nil tuples here in Action! This is an error")
		return
	}

	price, _ := t1.GetDouble("totalPrice")
	if price >= 20000 {
		fmt.Println("Congratulations you are eligible for Rs. 1000 gift coupon")
	}
}
