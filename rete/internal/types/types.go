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

type NwElemId interface {
	SetID(nw Network)
	GetID() int
	GetNw() Network
}
type NwElemIdImpl struct {
	ID int
	Nw Network
}

func (ide *NwElemIdImpl) SetID(nw Network) {
	ide.Nw = nw
	ide.ID = nw.GetIdGenService().GetNextID()
}
func (ide *NwElemIdImpl) GetID() int {
	return ide.ID
}
func (ide *NwElemIdImpl) GetNw() Network {
	return ide.Nw
}

type JoinTable interface {
	NwElemId
	GetName() string
	GetRule() model.Rule

	AddRow(handles []ReteHandle) JoinTableRow
	RemoveRow(rowID int) JoinTableRow
	GetRow(rowID int) JoinTableRow
	GetRowIterator() RowIterator

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

type RowIterator interface {
	HasNext() bool
	Next() JoinTableRow
	Remove() // remove underneath current element
}

type RowIDIterator interface {
	HasNext() bool
	Next() JoinTableRow
}

type NwService interface {
	model.Service
	GetNw() Network
}

type JtRefsService interface {
	NwService
	AddEntry(handle ReteHandle, jtName string, rowID int)
	RemoveRowEntry(handle ReteHandle, jtName string, rowID int)
	RemoveTableEntry(handle ReteHandle, jtName string)
	RemoveEntry(handle ReteHandle, jtName string)
	GetIterator(handle ReteHandle) HdlTblIterator
	GetRowIterator(handle ReteHandle, jtName string) RowIterator
}

type HdlTblIterator interface {
	HasNext() bool
	Next() JoinTable
	Remove()
}

type JtService interface {
	NwService
	GetOrCreateJoinTable(nw Network, rule model.Rule, identifiers []model.TupleType, name string) JoinTable
	GetJoinTable(name string) JoinTable
	//AddJoinTable(joinTable JoinTable)
	//RemoveJoinTable(name string)
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

type NwServiceImpl struct {
	Nw Network
}

func (nws *NwServiceImpl) GetNw() Network {
	return nws.Nw
}
