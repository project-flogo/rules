package ruleapi

import (
	"strconv"

	"github.com/project-flogo/rules/common/model"
)

type conditionImpl struct {
	name        string
	rule        model.Rule
	identifiers []model.TupleType
	cfn         model.ConditionEvaluator
	ctx         model.RuleContext
}

//NewCondition ... a new Condition
func newCondition(name string, rule model.Rule, identifiers []model.TupleType, cfn model.ConditionEvaluator, ctx model.RuleContext) model.Condition {
	c := conditionImpl{}
	c.initConditionImpl(name, rule, identifiers, cfn, ctx)
	return &c
}

func (cnd *conditionImpl) initConditionImpl(name string, rule model.Rule, identifiers []model.TupleType, cfn model.ConditionEvaluator, ctx model.RuleContext) {
	if name == "" {
		cndIdx := len(rule.GetConditions()) + 1
		name = "c_" + strconv.Itoa(cndIdx)
	}
	cnd.name = name
	cnd.rule = rule
	cnd.identifiers = append(cnd.identifiers, identifiers...)
	cnd.cfn = cfn
	cnd.ctx = ctx
}

func (cnd *conditionImpl) GetIdentifiers() []model.TupleType {
	return cnd.identifiers
}
func (cnd *conditionImpl) GetContext() model.RuleContext {
	return cnd.ctx
}

//func (cnd *conditionImpl) GetEvaluator() model.ConditionEvaluator {
//	return cnd.cfn
//}

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

func (cnd *conditionImpl) Evaluate(condName string, ruleNm string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) (bool, error) {
	result := false
	if cnd.cfn != nil {
		result = cnd.cfn(condName, ruleNm, tuples, ctx)
	}

	return result, nil
}
