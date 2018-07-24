package ruleapi

import (
	"github.com/TIBCOSoftware/bego/common/model"
)

type ruleImpl struct {
	name        string
	identifiers []model.TupleTypeAlias
	conditions  []model.Condition
	actionFn    model.ActionFunction
	priority    int
}

//NewRule ... Create a new rule
func NewRule(name string) model.MutableRule {
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

func (rule *ruleImpl) GetConditions() []model.Condition {
	return rule.conditions
}

func (rule *ruleImpl) SetActionFn(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

func (rule *ruleImpl) AddCondition(conditionName string, idrs []model.TupleTypeAlias, cfn model.ConditionEvaluator) {
	condition := newCondition(conditionName, rule, idrs, cfn)
	rule.conditions = append(rule.conditions, condition)
}

func (rule *ruleImpl) GetPriority() int {
	return rule.priority
}

func (rule *ruleImpl) SetPriority(priority int) {
	rule.priority = priority
}

//func (rule *ruleImpl) String() string {
//	str := ""
//	str += "[Rule: (" + rule.name + "\n"
//	return str
//}

func (rule *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + ") " + rule.name + "\n"
	//str += "[Rule: (" + strconv.Itoa(rule.id) + ") " + rule.name + "\n"

	str += "\t[Conditions:\n"
	for _, cond := range rule.conditions {
		str += "\t\t" + cond.String() + "\n"
	}
	// idrs := ""
	// for i := 0; i < len(rule.identifiers); i++ {
	// 	idrs += rule.identifiers[i].String() + ", "
	// }
	str += "\t[Idrs:" + model.IdentifiersToString(rule.identifiers) + "]\n"
	return str
	// return str + idrs + "]\n"
}

func (rule *ruleImpl) GetIdentifiers() []model.TupleTypeAlias {
	return rule.identifiers
}

func (rule *ruleImpl) SetAction(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

////IdentifiersToString Take a slice of Identifiers and return a string representation
//func IdentifiersToString(identifiers []model.TupleTypeAlias) string {
//	str := ""
//	for _, idr := range identifiers {
//		str += string(idr) + ", "
//	}
//	return str
//}
