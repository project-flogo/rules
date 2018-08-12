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
	AddCondition(conditionName string, idrs []TupleType, cFn ConditionEvaluator, ctx RuleContext)
	AddConditionWithDependency(conditionName string, idrs []string, cFn ConditionEvaluator, ctx RuleContext)
	SetAction(actionFn ActionFunction)
	SetPriority(priority int)
	SetContext(ctx RuleContext)
}

type Condition interface {
	GetName() string
	GetEvaluator() ConditionEvaluator
	GetRule() Rule
	GetIdentifiers() []TupleType
	GetContext() RuleContext
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
type ConditionEvaluator func(string, string, map[TupleType]Tuple, RuleContext) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[TupleType]Tuple, RuleContext)

type ValueChangeListener interface {
	OnValueChange(tuple Tuple, prop string)
}
