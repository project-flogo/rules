package model

import (
	"context"
)

//TupleType Each tuple is of a certain type, described by TypeDescriptor
type TupleType string

// TupleKey for each TupleDescriptor
type TupleKey interface {
	String() string
	GetTupleDescriptor() TupleDescriptor
}

// RuleContext associated with every rule
type RuleContext interface{}

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetIdentifiers() []TupleType
	GetConditions() []Condition
	GetActionFn() ActionFunction
	String() string
	GetPriority() int
	GetDeps() map[TupleType]map[string]bool
	GetContext() RuleContext
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []string, cFn ConditionEvaluator, ctx RuleContext) (err error)
	SetAction(actionFn ActionFunction)
	SetPriority(priority int)
	SetContext(ctx RuleContext)
}

//Condition interface to maintain/get various condition properties
type Condition interface {
	GetName() string
	GetEvaluator() ConditionEvaluator
	GetRule() Rule
	GetIdentifiers() []TupleType
	GetContext() RuleContext
	String() string
}

// RuleSession to maintain rules and assert tuples against those rules
type RuleSession interface {
	GetName() string

	AddRule(rule Rule) (int, bool)
	DeleteRule(ruleName string)

	Assert(ctx context.Context, tuple Tuple)
	Retract(ctx context.Context, tuple Tuple)

	ScheduleAssert(ctx context.Context, delayInMillis uint64, key interface{}, tuple Tuple)
	CancelScheduledAssert(ctx context.Context, key interface{})

	Unregister()
}

//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
//i.e, part of the server side API
type ConditionEvaluator func(string, string, map[TupleType]Tuple, RuleContext) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[TupleType]Tuple, RuleContext)

// ValueChangeListener to pickup and process tuple value changes
type ValueChangeListener interface {
	OnValueChange(tuple Tuple, prop string)
}
