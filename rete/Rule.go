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
	ruleImplVar := ruleImpl{}
	ruleImplVar.initRuleImpl(name)
	return &ruleImplVar
}
func (ruleImplVar *ruleImpl) initRuleImpl(name string) {
	currentNodeID++
	ruleImplVar.id = currentNodeID
	ruleImplVar.name = name
}

func (ruleImplVar *ruleImpl) GetName() string {
	return ruleImplVar.name
}

func (ruleImplVar *ruleImpl) GetID() int {
	return ruleImplVar.id
}

//GetConditions ... get the rule's condition set
func (ruleImplVar *ruleImpl) GetConditions() []condition {
	return ruleImplVar.conditions
}

func (ruleImplVar *ruleImpl) GetIdentifiers() []identifier {
	return ruleImplVar.identifiers
}

func (ruleImplVar *ruleImpl) AddCondition(conditionName string, idrs []model.StreamSource, cfn model.ConditionEvaluator) condition {
	strIds := make([]string, len(idrs))
	for i := 0; i < len(idrs); i++ {
		strIds[i] = string(idrs[i])
	}
	c := newCondition(conditionName, ruleImplVar, strIds, cfn)
	ruleImplVar.conditions = append(ruleImplVar.conditions, c)
	for _, idr := range strIds {
		// ruleImplVar.addIdentifier(idr[0], idr[1])
		ruleImplVar.addIdentifier(idr)
	}
	return c
}

func (ruleImplVar *ruleImpl) addIdentifier(identifierName string) identifier {

	idrNew := newIdentifier(identifierName)
	//TODO: Optimize this, perhaps using a map (need it to be ordered)
	//search for the idr. if exists, skip, else add it

	for _, idr := range ruleImplVar.identifiers {
		if idr.equals(idrNew) {
			return idr
		}
	}

	ruleImplVar.identifiers = append(ruleImplVar.identifiers, idrNew)
	return idrNew
}

func (ruleImplVar *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + strconv.Itoa(ruleImplVar.id) + ") " + ruleImplVar.name + "\n"
	str += "\t[Conditions:\n"
	for _, cond := range ruleImplVar.conditions {
		str += "\t\t" + cond.String() + "\n"
	}
	// idrs := ""
	// for i := 0; i < len(ruleImplVar.identifiers); i++ {
	// 	idrs += ruleImplVar.identifiers[i].String() + ", "
	// }
	str += "\t[Idrs:" + IdentifiersToString(ruleImplVar.identifiers) + "]\n"
	return str
	// return str + idrs + "]\n"
}

func (ruleImplVar *ruleImpl) SetAction(actionFn model.ActionFunction) {
	ruleImplVar.actionFn = actionFn
}

func (ruleImplVar *ruleImpl) GetActionFn() model.ActionFunction {
	return ruleImplVar.actionFn
}
