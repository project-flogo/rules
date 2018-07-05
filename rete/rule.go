package rete

import (
	"strconv"

	"github.com/TIBCOSoftware/bego/common/model"
)

//Rule ... a Rule interface
type Rule interface {
	GetName() string
	GetID() int
	GetConditions() []condition
	GetIdentifiers() []identifier
	GetActionFn() model.ActionFunction
	String() string
}

//MutableRule interface has methods to add conditions and actions
type MutableRule interface {
	Rule
	AddCondition(conditionName string, idrs []model.StreamSource, cFn model.ConditionEvaluator) condition
	SetAction(actionFn model.ActionFunction)
}

type ruleImpl struct {
	id          int
	name        string
	identifiers []identifier
	conditions  []condition
	actionFn    model.ActionFunction
}

//NewRule ... Create a new rule
func NewRule(name string) MutableRule {
	rule := ruleImpl{}
	rule.initRuleImpl(name)
	return &rule
}
func (rule *ruleImpl) initRuleImpl(name string) {
	currentNodeID++
	rule.id = currentNodeID
	rule.name = name
}

func (rule *ruleImpl) GetName() string {
	return rule.name
}

func (rule *ruleImpl) GetID() int {
	return rule.id
}

//GetConditions ... get the rule's condition set
func (rule *ruleImpl) GetConditions() []condition {
	return rule.conditions
}

func (rule *ruleImpl) GetIdentifiers() []identifier {
	return rule.identifiers
}

func (rule *ruleImpl) AddCondition(conditionName string, idrs []model.StreamSource, cfn model.ConditionEvaluator) condition {
	strIds := make([]string, len(idrs))
	for i := 0; i < len(idrs); i++ {
		strIds[i] = string(idrs[i])
	}
	c := newCondition(conditionName, rule, strIds, cfn)
	rule.conditions = append(rule.conditions, c)
	for _, idr := range strIds {
		// rule.addIdentifier(idr[0], idr[1])
		rule.addIdentifier(idr)
	}
	return c
}

func (rule *ruleImpl) addIdentifier(identifierName string) identifier {

	idrNew := newIdentifier(identifierName)
	//TODO: Optimize this, perhaps using a map (need it to be ordered)
	//search for the idr. if exists, skip, else add it

	for _, idr := range rule.identifiers {
		if idr.equals(idrNew) {
			return idr
		}
	}

	rule.identifiers = append(rule.identifiers, idrNew)
	return idrNew
}

func (rule *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + strconv.Itoa(rule.id) + ") " + rule.name + "\n"
	str += "\t[Conditions:\n"
	for _, cond := range rule.conditions {
		str += "\t\t" + cond.String() + "\n"
	}
	// idrs := ""
	// for i := 0; i < len(rule.identifiers); i++ {
	// 	idrs += rule.identifiers[i].String() + ", "
	// }
	str += "\t[Idrs:" + IdentifiersToString(rule.identifiers) + "]\n"
	return str
	// return str + idrs + "]\n"
}

func (rule *ruleImpl) SetAction(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

func (rule *ruleImpl) GetActionFn() model.ActionFunction {
	return rule.actionFn
}
