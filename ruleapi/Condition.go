package ruleapi

import "github.com/TIBCOSoftware/bego/common/model"

//condition ... is a rete condtion
type condition interface {
	GetName() string
	GetEvaluator() model.ConditionEvaluator
	GetRule() Rule
	GetStreamSource() []model.StreamSource
	//Stringer.String interface
	String() string
}

type conditionImpl struct {
	name        string
	rule        Rule
	identifiers []model.StreamSource
	cfn         model.ConditionEvaluator
}

//NewCondition ... a new Condition
func NewCondition(name string, rule Rule, identifiers []model.StreamSource, cfn model.ConditionEvaluator) condition {
	c := conditionImpl{}
	c.initConditionImpl(name, rule, identifiers, cfn)
	return &c
}

func (conditionImplVar *conditionImpl) initConditionImpl(name string, rule Rule, identifiers []model.StreamSource, cfn model.ConditionEvaluator) {
	conditionImplVar.name = name
	conditionImplVar.rule = rule
	conditionImplVar.identifiers = append(conditionImplVar.identifiers, identifiers...)
	conditionImplVar.cfn = cfn
}

func (conditionImplVar *conditionImpl) GetIdentifiers() []model.StreamSource {
	return conditionImplVar.identifiers
}

func (conditionImplVar *conditionImpl) GetEvaluator() model.ConditionEvaluator {
	return conditionImplVar.cfn
}

func (conditionImplVar *conditionImpl) String() string {
	return "[Condition: name:" + conditionImplVar.name + ", idrs: TODO]"
}

func (conditionImplVar *conditionImpl) GetName() string {
	return conditionImplVar.name
}

func (conditionImplVar *conditionImpl) GetRule() Rule {
	return conditionImplVar.rule
}
func (conditionImplVar *conditionImpl) GetStreamSource() []model.StreamSource {
	return conditionImplVar.identifiers
}
