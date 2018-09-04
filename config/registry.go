package config

import (
	"errors"

	"github.com/project-flogo/rules/common/model"
)

var (
	actionFunctions     = make(map[string]model.ActionFunction)
	conditionEvaluators = make(map[string]model.ConditionEvaluator)
	startupFunctions    = make(map[string]model.StartupRSFunction)
)

// Register registers the specified ActionFunction
func RegisterActionFunction(id string, actionFunction model.ActionFunction) error {

	if actionFunction == nil {
		return errors.New("cannot register 'nil' ActionFunction")
	}

	if _, dup := actionFunctions[id]; dup {
		return errors.New("ActionFunction already registered: " + id)
	}

	actionFunctions[id] = actionFunction

	return nil
}

// Get gets specified ActionFunction
func GetActionFunction(id string) model.ActionFunction {
	return actionFunctions[id]
}

// Register registers the specified ConditionEvaluator
func RegisterConditionEvaluator(id string, conditionEvaluator model.ConditionEvaluator) error {

	if conditionEvaluator == nil {
		return errors.New("cannot register 'nil' ConditionEvaluator")
	}

	if _, dup := conditionEvaluators[id]; dup {
		return errors.New("ConditionEvaluator already registered: " + id)
	}

	conditionEvaluators[id] = conditionEvaluator

	return nil
}

// Get gets specified ConditionEvaluator
func GetConditionEvaluator(id string) model.ConditionEvaluator {
	return conditionEvaluators[id]
}

// Register registers the specified StartupRSFunction
func RegisterStartupRSFunction(rsName string, startupFn model.StartupRSFunction) error {

	if startupFn == nil {
		return errors.New("cannot register 'nil' StartupRSFunction")
	}

	if _, dup := startupFunctions[rsName]; dup {
		return errors.New("StartupRSFunction already registered: " + rsName)
	}

	startupFunctions[rsName] = startupFn

	return nil
}

func GetStartupRSFunction(rsName string) (startupFn model.StartupRSFunction) {
	return startupFunctions[rsName]
}
