package rete

import (
	"container/list"
	"context"
	"github.com/TIBCOSoftware/bego/common/model"
)

type conflictRes interface {
	addAgendaItem(rule Rule, tupleMap map[model.StreamSource]model.StreamTuple)
	resolveConflict(ctx context.Context)
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

func getCrFromContext(ctx context.Context) conflictRes {
	var cr conflictRes
	reteCtxValPtr := ctx.Value(reteCTXKEY).(*reteCtxValType)
	cr = reteCtxValPtr.context["conflictRes"].(conflictRes)
	return cr
}
