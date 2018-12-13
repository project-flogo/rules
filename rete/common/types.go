package common

import (
	"context"
	"github.com/project-flogo/rules/common/model"
)

type RtcOprn int

const (
	ADD RtcOprn = 1 + iota
	RETRACT
	MODIFY
	DELETE
)

//Network ... these are used by RuleSession
type Network interface {
	AddRule(model.Rule) error
	String() string
	RemoveRule(string) model.Rule
	GetRules() []model.Rule
	Assert(ctx context.Context, rs model.RuleSession, tuple model.Tuple, changedProps map[string]bool, mode RtcOprn)
	Retract(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode RtcOprn)
	GetAssertedTuple(key model.TupleKey) model.Tuple
	RegisterRtcTransactionHandler(txnHandler model.RtcTransactionHandler, txnContext interface{})
	SetConfig(config map[string]string)
	GetConfigValue(key string) string
	GetConfig() map[string]string
}
