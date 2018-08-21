package config

import (
	"github.com/TIBCOSoftware/bego/common/model"
	"errors"
)

var (
	actionFunctions   = make(map[string]model.ActionFunction)
	conditionEvaluators   = make(map[string]model.ConditionEvaluator)
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
