package rete

import (
	"github.com/TIBCOSoftware/bego/common/model"
)

type agendaItem interface {
	getRule() model.Rule
	getTuples() map[model.TupleTypeAlias]model.StreamTuple
}

type agendaItemImpl struct {
	rule     model.Rule
	tupleMap map[model.TupleTypeAlias]model.StreamTuple
}

func newAgendaItem(rule model.Rule, tupleMap map[model.TupleTypeAlias]model.StreamTuple) agendaItem {
	ai := agendaItemImpl{rule, tupleMap}
	return &ai
}

func (ai *agendaItemImpl) getRule() model.Rule {
	return ai.rule
}

func (ai *agendaItemImpl) getTuples() map[model.TupleTypeAlias]model.StreamTuple {
	return ai.tupleMap
}
