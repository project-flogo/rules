package rete

import "github.com/TIBCOSoftware/bego/common/model"

//condition ... is a rete condtion
type condition interface {
	getName() string
	getIdentifiers() []identifier
	// Eval([]model.StreamTuple) bool

	//Stringer.String interface
	String() string
	getEvaluator() model.ConditionEvaluator
	getRule() Rule
}

type conditionImpl struct {
	name        string
	rule        Rule
	identifiers []identifier
	cfn         model.ConditionEvaluator
}

func newCondition(name string, rule Rule, identifiers []string, cfn model.ConditionEvaluator) condition {
	c := conditionImpl{}
	c.initConditionImpl(name, rule, identifiers, cfn)
	return &c
}

func (conditionImplVar *conditionImpl) initConditionImpl(name string, rule Rule, identifiers []string, cfn model.ConditionEvaluator) {
	conditionImplVar.name = name
	conditionImplVar.rule = rule
	for i := 0; i < len(identifiers); i++ {
		idName := identifiers[i]
		idr := newIdentifier(idName)
		conditionImplVar.identifiers = append(conditionImplVar.identifiers, idr)
	}
	conditionImplVar.cfn = cfn
}

func (conditionImplVar *conditionImpl) getIdentifiers() []identifier {
	return conditionImplVar.identifiers
}

func (conditionImplVar *conditionImpl) getEvaluator() model.ConditionEvaluator {
	return conditionImplVar.cfn
}

func (conditionImplVar *conditionImpl) String() string {
	return "[Condition: name:" + conditionImplVar.name + ", idrs:" + IdentifiersToString(conditionImplVar.identifiers) + "]"
}

func (conditionImplVar *conditionImpl) getName() string {
	return conditionImplVar.name
}

func (conditionImplVar *conditionImpl) getRule() Rule {
	return conditionImplVar.rule
}
