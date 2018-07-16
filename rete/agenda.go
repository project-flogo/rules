package rete

import (
	"github.com/TIBCOSoftware/bego/common/model"
)

type agendaItem interface {
	getRule() Rule
	getTuples() map[model.StreamSource]model.StreamTuple
}

type agendaItemImpl struct {
	rule     Rule
	tupleMap map[model.StreamSource]model.StreamTuple
}

func newAgendaItem(rule Rule, tupleMap map[model.StreamSource]model.StreamTuple) agendaItem {
	ai := agendaItemImpl{rule, tupleMap}
	return &ai
}

func (ai *agendaItemImpl) getRule() Rule {
	return ai.rule
}

func (ai *agendaItemImpl) getTuples() map[model.StreamSource]model.StreamTuple {
	return ai.tupleMap
}
