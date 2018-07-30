package model

import (
	"context"
	"time"
)

//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
//i.e, part of the server side API
type ConditionEvaluator func(string, string, map[TupleTypeAlias]StreamTuple) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[TupleTypeAlias]StreamTuple)

type RuleSession interface {
	AddRule(rule Rule) (int, bool)
	DeleteRule(ruleName string)

	Assert(ctx context.Context, tuple StreamTuple)
	Retract(ctx context.Context, tuple StreamTuple)
	Unregister() //remove itself from the package map
	GetName() string
}

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetIdentifiers() []TupleTypeAlias
	GetConditions() []Condition
	GetActionFn() ActionFunction
	String() string
	GetPriority() int
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []TupleTypeAlias, cFn ConditionEvaluator)
	SetAction(actionFn ActionFunction)
	SetPriority(priority int)
}

type Condition interface {
	GetName() string
	GetEvaluator() ConditionEvaluator
	GetRule() Rule
	GetIdentifiers() []TupleTypeAlias
	String() string
}

//TupleTypeAlias An internal representation of a 'DataSource'
type TupleTypeAlias string

//StreamTuple is a runtime representation of a stream of data
type StreamTuple interface {
	GetTypeAlias() TupleTypeAlias
	GetString(name string) string
	GetInt(name string) int
	GetFloat(name string) float64
	GetDateTime(name string) time.Time
}

type ValueChangeHandler interface {
	OnValueChange(tuple StreamTuple)
}

//MutableStreamTuple mutable part of the stream tuple
type MutableStreamTuple interface {
	StreamTuple
	SetString(ctx context.Context, name string, value string)
	SetInt(ctx context.Context, name string, value int)
	SetFloat(ctx context.Context, name string, value float64)
	SetDatetime(ctx context.Context, name string, value time.Time)
}
