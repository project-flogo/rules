package ruleapi

import (
	"github.com/tibmatt/bego/common/model"
)

type conditionImpl struct {
	name        string
	rule        model.Rule
	identifiers []model.TupleType
	cfn         model.ConditionEvaluator
	ctx         model.ConditionContext
}

//NewCondition ... a new Condition
func newCondition(name string, rule model.Rule, identifiers []model.TupleType, cfn model.ConditionEvaluator, ctx model.ConditionContext) model.Condition {
	c := conditionImpl{}
	c.initConditionImpl(name, rule, identifiers, cfn, ctx)
	return &c
}

func (cnd *conditionImpl) initConditionImpl(name string, rule model.Rule, identifiers []model.TupleType, cfn model.ConditionEvaluator, ctx model.ConditionContext) {
	cnd.name = name
	cnd.rule = rule
	cnd.identifiers = append(cnd.identifiers, identifiers...)
	cnd.cfn = cfn
	cnd.ctx = ctx
}

func (cnd *conditionImpl) GetIdentifiers() []model.TupleType {
	return cnd.identifiers
}
func (cnd *conditionImpl) GetContext() model.ConditionContext {
	return cnd.ctx
}

func (cnd *conditionImpl) GetEvaluator() model.ConditionEvaluator {
	return cnd.cfn
}

func (cnd *conditionImpl) String() string {
	return "[Condition: name:" + cnd.name + ", idrs: TODO]"
}

func (cnd *conditionImpl) GetName() string {
	return cnd.name
}

func (cnd *conditionImpl) GetRule() model.Rule {
	return cnd.rule
}
func (cnd *conditionImpl) GetTupleTypeAlias() []model.TupleType {
	return cnd.identifiers
}
