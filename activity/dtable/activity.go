package dtable

import (
	"context"
	"fmt"
	"strconv"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
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

	context := ctx.GetInput("ctx").(context.Context)
	tuples := ctx.GetInput("tuples").(map[model.TupleType]model.Tuple)

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

				tds := mutableTuple.GetTupleDescriptor()
				strVal := fmt.Sprintf("%v", act.Value)

				switch tds.GetProperty(act.Field).PropType {
				case data.TypeString:
					mutableTuple.SetString(context, act.Field, strVal)
				case data.TypeBool:
					b, err := strconv.ParseBool(strVal)
					if err == nil {
						mutableTuple.SetBool(context, act.Field, b)
					}
				case data.TypeInt:
					i, err := strconv.ParseInt(strVal, 10, 64)
					if err == nil {
						mutableTuple.SetInt(context, act.Field, int(i))
					}
				case data.TypeInt32:
					i, err := strconv.ParseInt(strVal, 10, 64)
					if err == nil {
						mutableTuple.SetInt(context, act.Field, int(i))
					}
				case data.TypeInt64:
					i, err := strconv.ParseInt(strVal, 10, 64)
					if err == nil {
						mutableTuple.SetLong(context, act.Field, i)
					}
				case data.TypeFloat32:
					f, err := strconv.ParseFloat(strVal, 32)
					if err == nil {
						mutableTuple.SetDouble(context, act.Field, f)
					}
				case data.TypeFloat64:
					f, err := strconv.ParseFloat(strVal, 64)
					if err == nil {
						mutableTuple.SetDouble(context, act.Field, f)
					}
				default:
					mutableTuple.SetValue(context, act.Field, act.Value)

				}

			}
		}
	}
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
