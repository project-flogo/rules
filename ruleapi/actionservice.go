package ruleapi

import (
	"context"
	"fmt"

	"github.com/project-flogo/core/action"
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
	Type     string
	Function model.ActionFunction
	Act      activity.Activity
	Action   action.Action
	Input    map[string]interface{}
}

// NewActionService creates new rule action service
func NewActionService(serviceCfg *config.ServiceDescriptor) (model.ActionService, error) {

	raService := &ruleActionService{
		Name:  serviceCfg.Name,
		Type:  serviceCfg.Type,
		Input: make(map[string]interface{}),
	}

	switch serviceCfg.Type {
	default:
		return nil, fmt.Errorf("type[%s] not supported for the service[%s]", serviceCfg.Type, serviceCfg.Name)
	case config.TypeServiceFunction:
		if serviceCfg.Function == nil {
			return nil, fmt.Errorf("service[%s] function can't empty", serviceCfg.Name)
		}
		raService.Function = serviceCfg.Function
	case config.TypeServiceActivity:
		// inflate activity from ref
		if serviceCfg.Ref[0] == '#' {
			var ok bool
			activityRef := serviceCfg.Ref
			serviceCfg.Ref, ok = support.GetAliasRef("activity", activityRef)
			if !ok {
				return nil, fmt.Errorf("activity '%s' not imported", activityRef)
			}
		}

		act := activity.Get(serviceCfg.Ref)
		if act == nil {
			return nil, fmt.Errorf("unsupported Activity:" + serviceCfg.Ref)
		}

		f := activity.GetFactory(serviceCfg.Ref)

		if f != nil {
			initCtx := newInitContext(serviceCfg.Name, serviceCfg.Settings, logger.ChildLogger(logger.RootLogger(), "ruleaction"))
			pa, err := f(initCtx)
			if err != nil {
				return nil, fmt.Errorf("unable to create rule action service '%s' : %s", serviceCfg.Name, err.Error())
			}
			act = pa
		}

		raService.Act = act

	case config.TypeServiceAction:
		if serviceCfg.Ref[0] == '#' {
			var ok bool
			actionRef := serviceCfg.Ref
			serviceCfg.Ref, ok = support.GetAliasRef("action", actionRef)
			if !ok {
				return nil, fmt.Errorf("action[%s] not imported", actionRef)
			}
		}

		actionFactory := action.GetFactory(serviceCfg.Ref)
		if actionFactory == nil {
			return nil, fmt.Errorf("factory not found for the action[%s]", serviceCfg.Ref)
		}

		actionCfg := &action.Config{Settings: serviceCfg.Settings}
		var err error
		raService.Action, err = actionFactory.New(actionCfg)
		if err != nil {
			return nil, fmt.Errorf("not able create action[%s] - %s", serviceCfg.Ref, err)
		}
	}

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

	tupleScope := make(map[string]interface{})
	for tk, t := range tuples {
		tupleScope[string(tk)] = t.GetMap()
	}

	scope := data.NewSimpleScope(tupleScope, nil)
	resolvedInputs, err := mapper.Apply(scope)
	if err != nil {
		return false, err
	}

	if raService.Type == config.TypeServiceActivity {
		// create activity context and set resolved inputs
		// TODO: implement context specific to rules instead of test package
		tc := test.NewActivityContext(raService.Act.Metadata())
		for k, v := range resolvedInputs {
			tc.SetInput(k, v)
		}

		return raService.Act.Eval(tc)
	} else if raService.Type == config.TypeServiceAction {
		syncAction, syncOk := raService.Action.(action.SyncAction)
		if syncOk && syncAction != nil {
			// sync action
			results, err := syncAction.Run(ctx, resolvedInputs)
			if err != nil {
				return false, fmt.Errorf("error while running the action service[%s] - %s", raService.Name, err)
			}
			fmt.Printf("service[%s] executed successfully. Service outputs: %s \n", raService.Name, results)
			return true, nil
		}

		asyncAction, asyncOk := raService.Action.(action.AsyncAction)
		if asyncOk && asyncAction != nil {
			err := asyncAction.Run(ctx, resolvedInputs, &actionResultHandler{name: raService.Name})
			if err != nil {
				return false, fmt.Errorf("error while running the action service[%s] - %s", raService.Name, err)
			}
			return true, nil
		}
	}
	return true, nil
}

type actionResultHandler struct {
	name string
}

// HandleResult is invoked when there are results available
func (arh *actionResultHandler) HandleResult(results map[string]interface{}, err error) {
	fmt.Printf("service[%s] outputs: %s \n", arh.name, results)
}

// Done indicates that the action has completed
func (arh *actionResultHandler) Done() {
	fmt.Printf("service[%s] executed successfully asynchronously\n", arh.name)
}
