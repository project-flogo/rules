package model

import (
	"context"
	"time"
)

type RuleSession interface {
	AddRule(rule Rule) (int, bool)
	DeleteRule(ruleName string)

	Assert(ctx context.Context, tuple Tuple)
	Retract(ctx context.Context, tuple Tuple)
	Unregister() //remove itself from the package map
	GetName() string
	RegisterTupleDescriptors (jsonRegistry string) //a json describing types
	//NewTuple(tupleType TupleType) MutableStreamTuple
	ValidateUpdate(alias TupleType, name string, value interface{}) bool

	DelayedAssert(ctx context.Context, delayInMillis uint64, key interface{}, tuple Tuple)
	CancelDelayedAssert (ctx context.Context, key interface{})
}

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
	AddCondition(conditionName string, idrs []TupleType, cFn ConditionEvaluator)
	AddConditionWithDependency(conditionName string, idrs []string, cFn ConditionEvaluator)
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

//TupleType An internal representation of a 'DataSource'
type TupleType string

//Tuple is a runtime representation of a stream of data
type Tuple interface {
	GetTypeAlias() TupleType
	GetString(name string) string
	GetInt(name string) int
	GetFloat(name string) float64
	GetDateTime(name string) time.Time
	GetProperties() []string
}


//MutableStreamTuple mutable part of the stream tuple
type MutableStreamTuple interface {
	Tuple
	SetString(ctx context.Context, rs RuleSession, name string, value string)
	SetInt(ctx context.Context, rs RuleSession, name string, value int)
	SetFloat(ctx context.Context, rs RuleSession, name string, value float64)
	SetDatetime(ctx context.Context, rs RuleSession, name string, value time.Time)
}

//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
//i.e, part of the server side API
type ConditionEvaluator func(string, string, map[TupleType]Tuple) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[TupleType]Tuple)

type ValueChangeHandler interface {
	OnValueChange(tuple Tuple, prop string)
}
