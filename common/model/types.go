package model

import (
	"context"
)

//TupleType Each tuple is of a certain type, described by TypeDescriptor
type TupleType string

type TupleKey interface {
	String() string
	GetTupleDescriptor() TupleDescriptor
}

type ConditionContext interface{}

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetIdentifiers() []TupleType
	GetConditions() []Condition
	GetActionFn() ActionFunction
	String() string
	GetPriority() int
	GetDeps() map[TupleType]map[string]bool
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []TupleType, cFn ConditionEvaluator, ctx ConditionContext)
	AddConditionWithDependency(conditionName string, idrs []string, cFn ConditionEvaluator, ctx ConditionContext)
	SetAction(actionFn ActionFunction)
	SetPriority(priority int)
}

type Condition interface {
	GetName() string
	GetEvaluator() ConditionEvaluator
	GetRule() Rule
	GetIdentifiers() []TupleType
	String() string
}

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
type ConditionEvaluator func(string, string, map[TupleType]Tuple) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[TupleType]Tuple)

type ValueChangeListener interface {
	OnValueChange(tuple Tuple, prop string)
}
