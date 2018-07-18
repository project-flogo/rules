package rete

import (
	"container/list"
	"context"

	"github.com/TIBCOSoftware/bego/common/model"
)

type conflictRes interface {
	addAgendaItem(rule Rule, tupleMap map[model.StreamSource]model.StreamTuple)
	resolveConflict(ctx context.Context)
	deleteAgendaFor(ctx context.Context, tuple model.StreamTuple)
}

type conflictResImpl struct {
	agendaList list.List
}

func newConflictRes() conflictRes {
	cr := conflictResImpl{}
	cr.initCR()
	return &cr
}

func (cr *conflictResImpl) initCR() {
	cr.agendaList = list.List{}
}

func (cr *conflictResImpl) addAgendaItem(rule Rule, tupleMap map[model.StreamSource]model.StreamTuple) {
	item := newAgendaItem(rule, tupleMap)
	v := rule.GetPriority()
	found := false
	for e := cr.agendaList.Front(); e != nil; e = e.Next() {
		curr := e.Value.(agendaItem)
		if v <= curr.getRule().GetPriority() {
			cr.agendaList.InsertBefore(item, e)
			found = true
			break
		}
	}
	if !found {
		cr.agendaList.PushBack(item)
	}
}

func (cr *conflictResImpl) resolveConflict(ctx context.Context) {
	var item agendaItem

	front := cr.agendaList.Front()
	for front != nil {
		val := cr.agendaList.Remove(front)
		if val != nil {
			item = val.(agendaItem)
			actionTuples := item.getTuples()
			actionFn := item.getRule().GetActionFn()
			if actionFn != nil {
				actionFn(item.getRule().GetName(), actionTuples)
			}
		}
		front = cr.agendaList.Front()
	}
}

func (cr *conflictResImpl) deleteAgendaFor(ctx context.Context, modifiedTuple model.StreamTuple) {

	hdlModified := getOrCreateHandle(ctx, modifiedTuple)

	for e := cr.agendaList.Front(); e != nil; {
		item := e.Value.(agendaItem)
		next := e.Next()
		for _, tuple := range item.getTuples() {
			hdl := getOrCreateHandle(ctx, tuple)
			if hdl == hdlModified { //this agendaitem has the modified tuple, remove the agenda item!
				cr.agendaList.Remove(e)
				break
			}
		}
		e = next
	}

}
