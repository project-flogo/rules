package rete

import "github.com/TIBCOSoftware/bego/common/model"

//condition ... is a rete condtion
type condition interface {
	getName() string
	getIdentifiers() []Identifier
	// Eval([]model.StreamTuple) bool

	//Stringer.String interface
	String() string
	getEvaluator() model.ConditionEvaluator
	getRule() Rule
}

type conditionImpl struct {
	name        string
	rule        Rule
	identifiers []Identifier
	cfn         model.ConditionEvaluator
}

// //NewCondition ... a new Condition
// func NewCondition(name string) MutableCondition {
// 	c := conditionImplVar{}
// 	c.initConditionImpl(name)
// 	return &c
// }

//NewCondition ... a new Condition
func newCondition(name string, rule Rule, identifiers []string, cfn model.ConditionEvaluator) condition {
	c := conditionImpl{}
	c.initConditionImpl(name, rule, identifiers, cfn)
	return &c
}

// func (conditionImplVar *conditionImplVar) initConditionImpl(name string) {
// 	conditionImplVar.name = name
// }
func (conditionImplVar *conditionImpl) initConditionImpl(name string, rule Rule, identifiers []string, cfn model.ConditionEvaluator) {
	conditionImplVar.name = name
	conditionImplVar.rule = rule
	for i := 0; i < len(identifiers); i++ {
		// idrAti := identifiers[i]
		// idName := idrAti[0]
		// idAlias := idrAti[1]
		// idr := newIdentifier(idName, idAlias)
		idName := identifiers[i]
		idr := newIdentifier(idName)
		conditionImplVar.identifiers = append(conditionImplVar.identifiers, idr)
	}
	conditionImplVar.cfn = cfn
}

func (conditionImplVar *conditionImpl) getIdentifiers() []Identifier {
	return conditionImplVar.identifiers
}

// func (conditionImplVar *conditionImpl) Eval([]model.StreamTuple) bool {
// 	return true
// }

func (conditionImplVar *conditionImpl) getEvaluator() model.ConditionEvaluator {
	return conditionImplVar.cfn
}

func (conditionImplVar *conditionImpl) String() string {
	return "[Condition: name:" + conditionImplVar.name + ", idrs:" + IdentifiersToString(conditionImplVar.identifiers) + "]"
}

// func (conditionImplVar *conditionImplVar) AddIdentifier(name string, alias string) Identifier {
// 	idr := NewIdentifier(name, alias)
// 	conditionImplVar.identifiers = append(conditionImplVar.identifiers, idr)
// 	return idr
// }

// func conditionsToStr(conditions []condition) string {
// 	str := ""
// 	for _, conditionVar := range conditions {
// 		str += conditionVar.String() + ","
// 	}
// 	return str
// }

func (conditionImplVar *conditionImpl) getName() string {
	return conditionImplVar.name
}

func (conditionImplVar *conditionImpl) getRule() Rule {
	return conditionImplVar.rule
}
