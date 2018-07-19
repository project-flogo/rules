package model

import (
	"context"
)

//ConditionEvaluator is a function pointer for handling condition evaluations on the server side
//i.e, part of the server side API
type ConditionEvaluator func(string, string, map[StreamSource]StreamTuple) bool

//ActionFunction is a function pointer for handling action callbacks on the server side
//i.e part of the server side API
type ActionFunction func(context.Context, RuleSession, string, map[StreamSource]StreamTuple)

type RuleSession interface {
	AddRule(rule Rule) (int, bool)
	DeleteRule(ruleName string)

	Assert(ctx context.Context, tuple StreamTuple)
	Retract(ctx context.Context, tuple StreamTuple)
	PrintNetwork()
}

type Rule interface {
	GetName() string
	GetActionFn() ActionFunction
	String() string
	GetConditions() []Condition
	GetPriority() int
}

type Condition interface {
	GetName() string
	GetEvaluator() ConditionEvaluator
	GetRule() Rule
	GetStreamSource() []StreamSource
	//Stringer.String interface
	String() string
}

// type sessionCtx interface {
// 	setRuleSession(rs RuleSession)
// 	getRuleSession() RuleSession
// }

type sessionKeyType struct {
}

var sessionCtxKEY = sessionKeyType{}

// type sessionCtxImpl struct {
// 	rs RuleSession
// }

// func newSessionCtx() sessionCtx {
// 	sCtx := sessionCtxImpl{}
// 	return &sCtx
// }

// func (sctx *sessionCtxImpl) setRuleSession(rs RuleSession) {
// 	sctx.rs = rs
// }

// func (sctx *sessionCtxImpl) getRuleSession() RuleSession {
// 	return sctx.rs
// }
