package ruleapi

import (
	"testing"

	_ "github.com/project-flogo/contrib/activity/log"
	"github.com/project-flogo/rules/config"
	"github.com/stretchr/testify/assert"
)

func TestNewActionService(t *testing.T) {
	cfg := &config.Service{
		Name:        "MyLogService",
		Description: "my log service",
	}
	aService, err := NewActionService(cfg)
	assert.NotNil(t, err)
	assert.Equal(t, "activity not specified for action service", err.Error())
	assert.Nil(t, aService)

	// activity reference
	cfg.Ref = "github.com/project-flogo/contrib/activity/log"
	aService, err = NewActionService(cfg)
	assert.Nil(t, err)
	assert.NotNil(t, aService)

	// set input
	input := map[string]interface{}{"message": "=$.n1.name"}
	aService.SetInput(input)

	// TODO: test aService.Execute()
}
