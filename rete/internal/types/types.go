package types

import (
	"container/list"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/common"
)

type Network interface {
	common.Network
	IncrementAndGetId() int
	GetJoinTable(joinTableID int) JoinTable
	AddToAllJoinTables(jT JoinTable)
}

type NwElemId interface {
	SetID(nw Network)
	GetID() int
}
type NwElemIdImpl struct {
	ID int
	Nw Network
}

func (ide *NwElemIdImpl) SetID(nw Network) {
	ide.Nw = nw
	ide.ID = nw.IncrementAndGetId()
}
func (ide *NwElemIdImpl) GetID() int {
	return ide.ID
}

type JoinTable interface {
	NwElemId
	GetRule() model.Rule

	AddRow(handles []ReteHandle) JoinTableRow
	RemoveRow(rowID int) JoinTableRow
	GetRow(rowID int) JoinTableRow
	GetRowIterator() RowIterator

	GetRowCount() int
}
type JoinTableRow interface {
	NwElemId
	GetHandles() []ReteHandle
}

type ReteHandle interface {
	NwElemId
	SetTuple(tuple model.Tuple)
	GetTuple() model.Tuple
	AddJoinTableRowRef(joinTableRowVar JoinTableRow, joinTableVar JoinTable)
	RemoveJoinTableRowRefs(changedProps map[string]bool)
	RemoveJoinTable(joinTableID int)
	GetTupleKey() model.TupleKey
	GetRefTableIterator() HdlTblIterator
}

type RowIterator interface {
	HasNext() bool
	Next() JoinTableRow
}

type JoinTableRefsInHdl interface {
	AddEntry(jointTableID int, rowID int)
	RemoveEntry(jointTableID int)
	GetIterator() HdlTblIterator
}

type HdlTblIterator interface {
	HasNext() bool
	Next() (int, *list.List)
}

type JoinTableCollection interface {
	GetJoinTable(joinTableID int) JoinTable
	AddJoinTable(joinTable JoinTable)
}

type HandleCollection interface {
	AddHandle(hdl ReteHandle)
	RemoveHandle(tuple model.Tuple) ReteHandle
	GetHandle(tuple model.Tuple) ReteHandle
	GetHandleByKey(key model.TupleKey) ReteHandle
}
