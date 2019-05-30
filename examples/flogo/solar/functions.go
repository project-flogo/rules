package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/rules/config"

	"github.com/project-flogo/rules/common/model"
)

//add this sample file to your flogo project
func init() {
	config.RegisterActionFunction("solarAction", solarAction)

	config.RegisterConditionEvaluator("solarEval", solarEval)

	config.RegisterStartupRSFunction("res://rulesession:solar", StartupRSFunction)
}

func solarAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	tHouse := tuples["house"]
	tSolar := tuples["solar"]
	fmt.Printf("Eligible for a solar promotion! [%s], [%s]\n", tHouse.GetKey().String(), tSolar.GetKey().String())

	// add isActionFired to rule context
	fmt.Printf("Rule fired: [%s]\n", ruleName)
}

func solarEval(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	tHouse := tuples["house"]
	tSolar := tuples["solar"]
	if tHouse == nil || tSolar == nil {
		fmt.Println("Should not get nil tuples here in JoinCondition! This is an error")
		return false
	}
	parcelHouse, _ := tHouse.GetString("parcel")
	parcelSolar, _ := tSolar.GetString("parcel")

	isSolarHouse, _ := tHouse.GetBool("is_solar")
	billSolar, _ := tSolar.GetDouble("bill")

	return (parcelHouse == parcelSolar) &&
		(isSolarHouse == false) &&
		(billSolar > 200)
}

func StartupRSFunction(ctx context.Context, rs model.RuleSession, startupCtx map[string]interface{}) (err error) {

	fmt.Printf("In startup rule function..\n")
	return nil
}
