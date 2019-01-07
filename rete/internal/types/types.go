package types

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/common"
)

type Network interface {
	common.Network
	GetPrefix() string
	GetIdGenService() IdGen
	GetJtService() JtService
	GetHandleService() HandleService
	GetJtRefService() JtRefsService
	GetTupleStore() model.TupleStore
}

type JoinTable interface {
	NwElemId
	GetName() string
	GetRule() model.Rule

	AddRow(handles []ReteHandle) JoinTableRow
	RemoveRow(rowID int) JoinTableRow
	GetRow(rowID int) JoinTableRow
	GetRowIterator() JointableRowIterator

	GetRowCount() int
	RemoveAllRows() //used when join table needs to be deleted
}

type JoinTableRow interface {
	NwElemId
	GetHandles() []ReteHandle
}

type ReteHandle interface {
	NwElemId
	SetTuple(tuple model.Tuple)
	GetTuple() model.Tuple
	GetTupleKey() model.TupleKey
}

type JtRefsService interface {
	NwService
	AddEntry(handle ReteHandle, jtName string, rowID int)
	RemoveRowEntry(handle ReteHandle, jtName string, rowID int)
	RemoveTableEntry(handle ReteHandle, jtName string)
	RemoveEntry(handle ReteHandle, jtName string)
	GetTableIterator(handle ReteHandle) JointableIterator
	GetRowIterator(handle ReteHandle, jtName string) JointableRowIterator
}

type JtService interface {
	NwService
	GetOrCreateJoinTable(nw Network, rule model.Rule, identifiers []model.TupleType, name string) JoinTable
	GetJoinTable(name string) JoinTable
}

type HandleService interface {
	NwService
	RemoveHandle(tuple model.Tuple) ReteHandle
	GetHandle(tuple model.Tuple) ReteHandle
	GetHandleByKey(key model.TupleKey) ReteHandle
	GetOrCreateHandle(nw Network, tuple model.Tuple) ReteHandle
}

type IdGen interface {
	NwService
	GetMaxID() int
	GetNextID() int
}

type JointableIterator interface {
	HasNext() bool
	Next() JoinTable
	Remove()
}

type JointableRowIterator interface {
	HasNext() bool
	Next() JoinTableRow
	Remove()
}
