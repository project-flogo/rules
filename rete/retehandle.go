package rete

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
)

const (
	one 		   = 0x01
	rtcIsAsserted  = one
	rtcIsModified  = one << 1
	rtcIsRetracted = one << 2

	isNew          = one << 3

)

//Holds a tuple reference and related state
type reteHandle interface {
	setTuple(tuple model.Tuple)
	getTuple() model.Tuple
	addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable)
	removeJoinTableRowRefs(changedProps map[string]bool)
	removeJoinTable(joinTableVar joinTable)

	setAsserted()
	setModified()
	setRetracted()
	setNew()

	isAsserted() bool
	isModified() bool
	isRetracted() bool
	isNew() bool


}

type handleImpl struct {
	tuple         model.Tuple
	tablesAndRows map[joinTable]*list.List

	rtcStatus	  uint8
}

func (hdl *handleImpl) setTuple(tuple model.Tuple) {
	hdl.tuple = tuple
}

func (hdl *handleImpl) initHandleImpl() {
	hdl.tablesAndRows = make(map[joinTable]*list.List)
	hdl.rtcStatus = 0x00
}

func (hdl *handleImpl) getTuple() model.Tuple {
	return hdl.tuple
}

func getOrCreateHandle(ctx context.Context, tuple model.Tuple) reteHandle {
	reteCtxVar := getReteCtx(ctx)
	return reteCtxVar.getNetwork().getOrCreateHandle(ctx, tuple)
}

func (hdl *handleImpl) addJoinTableRowRef(joinTableRowVar joinTableRow, joinTableVar joinTable) {

	rowsForJoinTable := hdl.tablesAndRows[joinTableVar]
	if rowsForJoinTable == nil {
		rowsForJoinTable = list.New()
		hdl.tablesAndRows[joinTableVar] = rowsForJoinTable
	}
	rowsForJoinTable.PushBack(joinTableRowVar)

}

func (hdl *handleImpl) removeJoinTableRowRefs(changedProps map[string]bool) {

	tuple := hdl.tuple
	alias := tuple.GetTupleType()

	emptyJoinTables := list.New()

	for joinTable, listOfRows := range hdl.tablesAndRows {

		toDelete := false
		if changedProps != nil {
			rule := joinTable.getRule()
			depProps, found := rule.GetDeps()[alias]
			if found { // rule depends on this type
				for changedProp := range changedProps {
					_, foundProp := depProps[changedProp]
					if foundProp {
						toDelete = true
						break
					}
				}
			}
		} else {
			toDelete = true
		}

		if !toDelete {
			continue
		}

		for e := listOfRows.Front(); e != nil; e = e.Next() {
			row := e.Value.(joinTableRow)
			joinTable.removeRow(row)
		}
		if joinTable.len() == 0 {
			emptyJoinTables.PushBack(joinTable)
		}
	}

	for e := emptyJoinTables.Front(); e != nil; e = e.Next() {
		emptyJoinTable := e.Value.(joinTable)
		delete(hdl.tablesAndRows, emptyJoinTable)
	}
}

//Used when a rule is deleted. See Network.RemoveRule
func (hdl *handleImpl) removeJoinTable(joinTableVar joinTable) {
	_, ok := hdl.tablesAndRows[joinTableVar]
	if ok {
		delete(hdl.tablesAndRows, joinTableVar)
	}
}
func (hdl *handleImpl) setAsserted() {
	hdl.rtcStatus = hdl.rtcStatus | rtcIsAsserted
}
func (hdl *handleImpl) setModified() {
	hdl.rtcStatus = hdl.rtcStatus | rtcIsModified
}
func (hdl *handleImpl) setRetracted() {
	hdl.rtcStatus = hdl.rtcStatus | rtcIsRetracted
}
func (hdl *handleImpl) setNew() {
	hdl.rtcStatus = hdl.rtcStatus | isNew
}

func (hdl *handleImpl) isAsserted() bool {
	return (hdl.rtcStatus & rtcIsAsserted) == rtcIsAsserted
}
func (hdl *handleImpl) isModified() bool {
	return (hdl.rtcStatus & rtcIsModified) == rtcIsModified
}
func (hdl *handleImpl) isRetracted() bool {
	return (hdl.rtcStatus & rtcIsRetracted) == rtcIsRetracted
}
func (hdl *handleImpl) isNew() bool {
	return (hdl.rtcStatus & isNew) == isNew
}