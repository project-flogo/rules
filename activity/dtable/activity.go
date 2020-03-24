package dtable

import (
	"context"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi/dtable"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity decision table based rule action
type Activity struct {
	dtable *dtable.DTable
}

// New creates new decision table activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	// Read settings
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	// Read decision table from file
	dtable, err := dtable.LoadFromFile(settings.DTableFile)
	if err != nil {
		return nil, err
	}
	// dtable.print()
	err = dtable.Compile()
	if err != nil {
		return nil, err
	}
	// dtable.print()

	// Read setting from init context
	act := &Activity{
		dtable: dtable,
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

	// evaluate decision table
	a.dtable.Apply(context, tuples)

	return true, nil
}
