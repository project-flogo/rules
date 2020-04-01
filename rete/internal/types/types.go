package types

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/common"
)

type Network interface {
	common.Network
	GetHandleWithTuple(ctx context.Context, tuple model.Tuple) ReteHandle
	AssertInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error
	RetractInternal(ctx context.Context, tuple model.Tuple, changedProps map[string]bool, mode common.RtcOprn) error
	GetPrefix() string
	GetIdGenService() IdGen
	GetLockService() LockService
	GetJtService() JtService
	GetHandleService() HandleService
	GetJtRefService() JtRefsService
	GetTupleStore() model.TupleStore
}

type ConflictRes interface {
	AddAgendaItem(rule model.Rule, tupleMap map[model.TupleType]model.Tuple)
	ResolveConflict(ctx context.Context)
	DeleteAgendaFor(ctx context.Context, tuple model.Tuple, changeProps map[string]bool)
}

type ReteCtx interface {
	GetConflictResolver() ConflictRes
	GetOpsList() *list.List
	GetNetwork() Network
	GetRuleSession() model.RuleSession
	OnValueChange(tuple model.Tuple, prop string)

	GetRtcAdded() map[string]model.Tuple
	GetRtcModified() map[string]model.RtcModified
	GetRtcDeleted() map[string]model.Tuple

	AddToRtcAdded(tuple model.Tuple)
	AddToRtcModified(tuple model.Tuple)
	AddToRtcDeleted(tuple model.Tuple)
	AddRuleModifiedToOpsList()

	Normalize()
	CopyRuleModifiedToRtcModified()
	ResetModified()

	PrintRtcChangeList()
}

type JoinTable interface {
	NwElemId
	GetName() string
	GetRule() model.Rule

	AddRow(handles []ReteHandle) JoinTableRow
	RemoveRow(rowID int) JoinTableRow
	GetRow(ctx context.Context, rowID int) JoinTableRow
	GetRowIterator(ctx context.Context) JointableRowIterator

	GetRowCount() int
	RemoveAllRows(ctx context.Context) //used when join table needs to be deleted
}

type JoinTableRow interface {
	NwElemId
	Write()
	GetHandles() []ReteHandle
}

type ReteHandleStatus uint

const (
	ReteHandleStatusUnknown ReteHandleStatus = iota
	ReteHandleStatusCreating
	ReteHandleStatusCreated
	ReteHandleStatusDeleting
	ReteHandleStatusRetracting
	ReteHandleStatusRetracted
)

type ReteHandle interface {
	NwElemId
	SetTuple(tuple model.Tuple)
	GetTuple() model.Tuple
	GetTupleKey() model.TupleKey
	SetStatus(status ReteHandleStatus)
	GetStatus() ReteHandleStatus
	Unlock()
}

type JtRefsService interface {
	NwService
	AddEntry(handle ReteHandle, jtName string, rowID int)
	RemoveEntry(handle ReteHandle, jtName string, rowID int)
	GetRowIterator(ctx context.Context, handle ReteHandle) JointableIterator
}

type JtService interface {
	NwService
	GetOrCreateJoinTable(nw Network, rule model.Rule, identifiers []model.TupleType, name string) JoinTable
	GetJoinTable(name string) JoinTable
}

type HandleService interface {
	NwService
	RemoveHandle(tuple model.Tuple) ReteHandle
	GetHandle(ctx context.Context, tuple model.Tuple) ReteHandle
	GetHandleByKey(ctx context.Context, key model.TupleKey) ReteHandle
	GetHandleWithTuple(nw Network, tuple model.Tuple) ReteHandle
	GetOrCreateLockedHandle(nw Network, tuple model.Tuple) (ReteHandle, bool)
	GetLockedHandle(nw Network, tuple model.Tuple) (handle ReteHandle, locked, dne bool)
	GetAllHandles(nw Network) map[string]ReteHandle
}

type IdGen interface {
	NwService
	GetMaxID() int
	GetNextID() int
}

type LockService interface {
	NwService
	Lock()
	Unlock()
}

type JointableIterator interface {
	HasNext() bool
	Next() (JoinTableRow, JoinTable)
	Remove()
}

type JointableRowIterator interface {
	HasNext() bool
	Next() JoinTableRow
	Remove()
}
