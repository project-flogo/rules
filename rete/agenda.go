package rete

import (
	"github.com/project-flogo/rules/common/model"
)

type agendaItem interface {
	getRule() model.Rule
	getTuples() map[model.TupleType]model.Tuple
}

type agendaItemImpl struct {
	rule     model.Rule
	tupleMap map[model.TupleType]model.Tuple
}

func newAgendaItem(rule model.Rule, tupleMap map[model.TupleType]model.Tuple) agendaItem {
	ai := agendaItemImpl{rule, tupleMap}
	return &ai
}

func (ai *agendaItemImpl) getRule() model.Rule {
	return ai.rule
}

func (ai *agendaItemImpl) getTuples() map[model.TupleType]model.Tuple {
	return ai.tupleMap
}
