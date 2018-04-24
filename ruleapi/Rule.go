package ruleapi

import "github.com/TIBCOSoftware/bego/common/model"

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetActionFn() model.ActionFunction
	String() string
	GetConditions() []condition
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []model.StreamSource, cFn model.ConditionEvaluator)
	SetActionFn(actionFn model.ActionFunction)
}

type ruleImpl struct {
	name       string
	conditions []condition
	actionFn   model.ActionFunction
}

//NewRule ... Create a new rule
func NewRule(name string) MutableRule {
	ruleImplVar := ruleImpl{}
	ruleImplVar.initRuleImpl(name)
	return &ruleImplVar
}
func (ruleImplVar *ruleImpl) initRuleImpl(name string) {
	ruleImplVar.name = name
}

func (ruleImplVar *ruleImpl) GetName() string {
	return ruleImplVar.name
}

func (ruleImplVar *ruleImpl) GetActionFn() model.ActionFunction {
	return ruleImplVar.actionFn
}

func (ruleImplVar *ruleImpl) GetConditions() []condition {
	return ruleImplVar.conditions
}

func (ruleImplVar *ruleImpl) SetActionFn(actionFn model.ActionFunction) {
	ruleImplVar.actionFn = actionFn
}

func (ruleImplVar *ruleImpl) AddCondition(conditionName string, idrs []model.StreamSource, cfn model.ConditionEvaluator) {
	condition := NewCondition(conditionName, ruleImplVar, idrs, cfn)
	ruleImplVar.conditions = append(ruleImplVar.conditions, condition)
}

func (ruleImplVar *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + ruleImplVar.name + "\n"
	return str
}
