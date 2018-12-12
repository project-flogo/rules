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

//Network ... the rete network
type Network interface {
	AddRule(model.Rule) error
	String() string
	RemoveRule(string) model.Rule
	GetRules() []model.Rule
	Assert(ctx context.Context, rs model.RuleSession, tuple model.Tuple, changedProps map[string]bool, mode RtcOprn)
	Retract(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode RtcOprn)
	GetAssertedTuple(key model.TupleKey) model.Tuple
	GetAssertedTupleByStringKey(key string) model.Tuple
	RegisterRtcTransactionHandler(txnHandler model.RtcTransactionHandler, txnContext interface{})
	SetConfig(config map[string]string)
	GetConfigValue(key string) string
	GetConfig() map[string]string

	//private
	//retractInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode RtcOprn)
	//assertInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode RtcOprn)
	//getOrCreateHandle(ctx context.Context, tuple model.Tuple) ReteHandle
	//getHandle(tuple model.Tuple) ReteHandle
	//IncrementAndGetId() int
	//GetJoinTable(joinTableID int) JoinTable
	//getFactory() TypeFactory

	//AddToAllJoinTables (jT JoinTable)
}
