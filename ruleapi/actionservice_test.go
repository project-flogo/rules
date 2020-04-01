package ruleapi

import (
	"context"
	"testing"

	_ "github.com/project-flogo/contrib/activity/log"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/stretchr/testify/assert"
)

func TestNewActionService(t *testing.T) {
	cfg := &config.ServiceDescriptor{
		Name:        "TestActionService",
		Description: "test action service",
	}
	aService, err := NewActionService(cfg)
	assert.NotNil(t, err)
	assert.Equal(t, "service type can't be empty", err.Error())
	assert.Nil(t, aService)

	// unsupported service type
	cfg.Type = "unknowntype"
	aService, err = NewActionService(cfg)
	assert.NotNil(t, err)
	assert.Equal(t, "service type - 'unknowntype' is not supported", err.Error())

	// action service with function
	cfg.Function = emptyAction
	cfg.Type = config.TypeServiceFunction
	aService, err = NewActionService(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, aService)
	cfg.Function = nil //clear for next test scenario

	// action service with activity
	cfg.Ref = "github.com/project-flogo/contrib/activity/log"
	cfg.Type = config.TypeServiceActivity
	aService, err = NewActionService(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, aService)

	// set input
	input := map[string]interface{}{"message": "=$.n1.name"}
	aService.SetInput(input)

	// TODO: test aService.Execute()
}

func emptyAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
}
