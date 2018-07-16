package ruleapi

import "github.com/TIBCOSoftware/bego/common/model"

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetActionFn() model.ActionFunction
	String() string
	GetConditions() []condition
	GetPriority() int
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []model.StreamSource, cFn model.ConditionEvaluator)
	SetActionFn(actionFn model.ActionFunction)
	SetPriority(priority int)
}

type ruleImpl struct {
	name       string
	conditions []condition
	actionFn   model.ActionFunction
	priority   int
}

//NewRule ... Create a new rule
func NewRule(name string) MutableRule {
	rule := ruleImpl{}
	rule.initRuleImpl(name)
	return &rule
}
func (rule *ruleImpl) initRuleImpl(name string) {
	rule.name = name
}

func (rule *ruleImpl) GetName() string {
	return rule.name
}

func (rule *ruleImpl) GetActionFn() model.ActionFunction {
	return rule.actionFn
}

func (rule *ruleImpl) GetConditions() []condition {
	return rule.conditions
}

func (rule *ruleImpl) SetActionFn(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

func (rule *ruleImpl) AddCondition(conditionName string, idrs []model.StreamSource, cfn model.ConditionEvaluator) {
	condition := NewCondition(conditionName, rule, idrs, cfn)
	rule.conditions = append(rule.conditions, condition)
}

func (rule *ruleImpl) GetPriority() int {
	return rule.priority
}

func (rule *ruleImpl) SetPriority(priority int) {
	rule.priority = priority
}

func (rule *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + rule.name + "\n"
	return str
}
