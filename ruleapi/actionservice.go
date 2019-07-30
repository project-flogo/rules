package ruleapi

import (
	"fmt"

	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data/mapper"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/core/support"
	logger "github.com/project-flogo/core/support/log"
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

// NewActionService creates new action service
func NewActionService(config *config.Service) (*model.ActionService, error) {

	if config.Ref == "" {
		return nil, fmt.Errorf("activity not specified for action service")
	}

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
		initCtx := newInitContext(config.Name, config.Settings, logger.ChildLogger(logger.RootLogger(), "ACTIONSERVICE"))
		pa, err := f(initCtx)
		if err != nil {
			return nil, fmt.Errorf("unable to create stage '%s' : %s", config.Ref, err.Error())
		}
		act = pa
	}

	aService := &model.ActionService{}
	aService.Name = config.Name
	aService.Act = act

	return aService, nil
}
