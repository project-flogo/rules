package ruleapi

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support"
	logger "github.com/project-flogo/core/support/log"
	"github.com/project-flogo/core/support/test"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
)

// activity context
type initContext struct {
	settings      map[string]interface{}
	mapperFactory mapper.Factory
	logger        logger.Logger
}

func newInitContext(name string, settings map[string]interface{}, log logger.Logger) *initContext {
	return &initContext{
		settings:      settings,
		mapperFactory: mapper.NewFactory(resolve.GetBasicResolver()),
		logger:        logger.ChildLogger(log, name),
	}
}

func (i *initContext) Settings() map[string]interface{} {
	return i.settings
}

func (i *initContext) MapperFactory() mapper.Factory {
	return i.mapperFactory
}

func (i *initContext) Logger() logger.Logger {
	return i.logger
}

// rule action service
type ruleActionService struct {
	Name     string
	Function model.ActionFunction
	Act      activity.Activity
	Input    map[string]interface{}
}

// NewActionService creates new rule action service
func NewActionService(config *config.ServiceDescriptor) (model.ActionService, error) {

	if config.Function == nil && config.Ref == "" {
		return nil, fmt.Errorf("both service function & ref can not be empty")
	}

	raService := &ruleActionService{
		Name:  config.Name,
		Input: make(map[string]interface{}),
	}

	if config.Function != nil {
		raService.Function = config.Function
		return raService, nil
	}

	// inflate activity from ref
	if config.Ref[0] == '#' {
		var ok bool
		activityRef := config.Ref
		config.Ref, ok = support.GetAliasRef("activity", activityRef)
		if !ok {
			return nil, fmt.Errorf("activity '%s' not imported", activityRef)
		}
	}

	act := activity.Get(config.Ref)
	if act == nil {
		return nil, fmt.Errorf("unsupported Activity:" + config.Ref)
	}

	f := activity.GetFactory(config.Ref)

	if f != nil {
		initCtx := newInitContext(config.Name, config.Settings, logger.ChildLogger(logger.RootLogger(), "ruleaction"))
		pa, err := f(initCtx)
		if err != nil {
			return nil, fmt.Errorf("unable to create rule action service '%s' : %s", config.Name, err.Error())
		}
		act = pa
	}

	raService.Act = act

	return raService, nil
}

// SetInput sets input
func (raService *ruleActionService) SetInput(input map[string]interface{}) {
	for k, v := range input {
		raService.Input[k] = v
	}
}

// Execute execute rule action service
func (raService *ruleActionService) Execute(ctx context.Context, rs model.RuleSession, rName string, tuples map[model.TupleType]model.Tuple, rCtx model.RuleContext) (done bool, err error) {
	// invoke function and return, if available
	if raService.Function != nil {
		raService.Function(ctx, rs, rName, tuples, rCtx)
		return true, nil
	}

	// resolve inputs from tuple scope
	mFactory := mapper.NewFactory(resolve.GetBasicResolver())
	mapper, err := mFactory.NewMapper(raService.Input)
	if err != nil {
		return false, err
	}

	toupleScope := make(map[string]interface{})
	for tk, t := range tuples {
		toupleScope[string(tk)] = t.GetMap()
	}

	scope := data.NewSimpleScope(toupleScope, nil)
	resolvedInputs, err := mapper.Apply(scope)
	if err != nil {
		return false, err
	}

	// create activity context and set resolved inputs
	// TODO: implement context specific to rules instead of test package
	tc := test.NewActivityContext(raService.Act.Metadata())
	for k, v := range resolvedInputs {
		tc.SetInput(k, v)
	}

	return raService.Act.Eval(tc)
}
