package dtable

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/rules/common/model"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity decision table based rule action
type Activity struct {
	tasks []*Task
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
		tuple := tuples[model.TupleType(task.Tuple)]
		if tuple == nil {
			fmt.Printf("tuple[%s] not found \n", task.Tuple)
			continue
		}
		fmt.Printf("tuple = %v \n", tuple)
		mutableTuple := tuple.(model.MutableTuple)
		mutableTuple.SetString(ctx1, task.Field, task.To)
		fmt.Printf("updated tuple = %v \n", tuple)
	}

	fmt.Println("DECISION TABLE EXIT")
	return false, nil
}
