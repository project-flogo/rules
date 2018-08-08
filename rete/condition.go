package rete

import "github.com/TIBCOSoftware/bego/common/model"

type conditionImpl struct {
	name        string
	rule        model.Rule
	identifiers []model.TupleType
	cfn         model.ConditionEvaluator
}

func newCondition(name string, rule model.Rule, identifiers []model.TupleType, cfn model.ConditionEvaluator) model.Condition {
	c := conditionImpl{}
	c.initConditionImpl(name, rule, identifiers, cfn)
	return &c
}

//
func (cnd *conditionImpl) initConditionImpl(name string, rule model.Rule, identifiers []model.TupleType, cfn model.ConditionEvaluator) {
	cnd.name = name
	cnd.rule = rule
	for i := 0; i < len(identifiers); i++ {
		idName := identifiers[i]
		idr := idName
		cnd.identifiers = append(cnd.identifiers, idr)
	}
	cnd.cfn = cfn
}

//
func (cnd *conditionImpl) GetIdentifiers() []model.TupleType {
	return cnd.identifiers
}

//
func (cnd *conditionImpl) GetEvaluator() model.ConditionEvaluator {
	return cnd.cfn
}

func (cnd *conditionImpl) String() string {
	return "[Condition: name:" + cnd.name + ", idrs:" + model.IdentifiersToString(cnd.identifiers) + "]"
}

func (cnd *conditionImpl) GetName() string {
	return cnd.name
}

//
func (cnd *conditionImpl) GetRule() model.Rule {
	return cnd.rule
}
