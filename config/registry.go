package config

import (
	"errors"
	"reflect"

	"github.com/project-flogo/rules/common/model"
)

var (
	actionFunctions     = make(map[string]model.ActionFunction)
	conditionEvaluators = make(map[string]model.ConditionEvaluator)
	startupFunctions    = make(map[string]model.StartupRSFunction)
)

// RegisterActionFunction registers the specified ActionFunction
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

// GetActionFunction gets specified ActionFunction
func GetActionFunction(id string) model.ActionFunction {
	return actionFunctions[id]
}

// GetActionFunctionID get ActionFunction id based on the function reference
func GetActionFunctionID(actionFn model.ActionFunction) string {
	actionFnToCheck := reflect.ValueOf(actionFn)
	for key, value := range actionFunctions {
		valueFn := reflect.ValueOf(value)
		if valueFn.Pointer() == actionFnToCheck.Pointer() {
			return key
		}
	}
	return ""
}

// RegisterConditionEvaluator registers the specified ConditionEvaluator
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

// GetConditionEvaluator gets specified ConditionEvaluator
func GetConditionEvaluator(id string) model.ConditionEvaluator {
	return conditionEvaluators[id]
}

// GetConditionEvaluatorID gets ConditionEvaluator Id based on the function reference
func GetConditionEvaluatorID(conditionEvaluator model.ConditionEvaluator) string {
	conditionEvaluatorToCheck := reflect.ValueOf(conditionEvaluator)
	for key, value := range conditionEvaluators {
		valueFn := reflect.ValueOf(value)
		if valueFn.Pointer() == conditionEvaluatorToCheck.Pointer() {
			return key
		}
	}
	return ""
}

// RegisterStartupRSFunction registers the specified StartupRSFunction
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

// GetStartupRSFunction gets registered StartupRSFunction
func GetStartupRSFunction(rsName string) (startupFn model.StartupRSFunction) {
	return startupFunctions[rsName]
}
