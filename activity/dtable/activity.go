package dtable

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity decision table based rule action
type Activity struct {
	tasks []*DecisionTable
}

// New creates new decision table activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	// Read settings
	settings := &Settings{}
	err := settings.FromMap(ctx.Settings())
	if err != nil {
		return nil, err
	}
	fmt.Println("settings: ", settings)
	fmt.Println("settings: ", settings.Make[0])

	// Read setting from init context
	act := &Activity{
		tasks: settings.Make,
	}
	return act, nil
}

// Metadata activity metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activity.ToMetadata(&Input{})
}

// Eval implements decision table action
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	fmt.Println("DECISION TABLE ENTER")

	ctx1 := ctx.GetInput("ctx").(context.Context)
	// rs := ctx.GetInput("rulesession").(model.RuleSession)
	// rName := ctx.GetInput("rulename").(string)
	tuples := ctx.GetInput("tuples").(map[model.TupleType]model.Tuple)
	// rCtx := ctx.GetInput("rulecontext").(model.RuleContext)

	// Run tasks
	for _, task := range a.tasks {

		conditionEval := false
		for _, cond := range task.DtConditions {

			eval := evaluateCondition(cond, tuples)
			if eval {
				conditionEval = true
				continue
			} else {
				conditionEval = false
				break
			}
		}

		if conditionEval {
			for _, act := range task.DtActions {
				tuple := tuples[model.TupleType(act.Tuple)]
				if tuple == nil {
					continue
				}
				mutableTuple := tuple.(model.MutableTuple)
				mutableTuple.SetString(ctx1, act.Field, act.Value)
			}
		}

	}

	fmt.Println("DECISION TABLE EXIT")
	return false, nil
}

func evaluateCondition(cond *DtCondition, tuples map[model.TupleType]model.Tuple) bool {

	condExprsn := "$." + cond.Tuple + "." + cond.Field + " " + cond.Expr

	condExprs := ruleapi.NewExprCondition(condExprsn)
	res, err := condExprs.Evaluate("", "", tuples, "")
	if err != nil {
		return false
	}

	return res
}
