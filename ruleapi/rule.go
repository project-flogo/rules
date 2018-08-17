package ruleapi

import (
	"strings"

	"errors"

	"github.com/TIBCOSoftware/bego/common/model"
)

type ruleImpl struct {
	name        string
	identifiers []model.TupleType
	conditions  []model.Condition
	actionFn    model.ActionFunction
	priority    int
	deps        map[model.TupleType]map[string]bool
	ctx         model.RuleContext
}

func (rule *ruleImpl) GetContext() model.RuleContext {
	return rule.ctx
}

func (rule *ruleImpl) SetContext(ctx model.RuleContext) {
	rule.ctx = ctx
}

//NewRule ... Create a new rule
func NewRule(name string) model.MutableRule {
	rule := ruleImpl{}
	rule.initRuleImpl(name)
	return &rule
}

func (rule *ruleImpl) initRuleImpl(name string) {
	rule.name = name
	rule.identifiers = []model.TupleType{}
	rule.conditions = []model.Condition{}
	rule.deps = make(map[model.TupleType]map[string]bool)
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

func (rule *ruleImpl) addCond(conditionName string, idrs []model.TupleType, cfn model.ConditionEvaluator, ctx model.RuleContext, setIdr bool) {
	condition := newCondition(conditionName, rule, idrs, cfn, ctx)
	rule.conditions = append(rule.conditions, condition)

	for _, cidr := range idrs {
		if len(rule.identifiers) == 0 {
			rule.identifiers = append(rule.identifiers, cidr)
		} else {
			for _, ridr := range rule.identifiers {
				if cidr != ridr {
					rule.identifiers = append(rule.identifiers, cidr)
					break
				}
			}
		}
	}

}

func (rule *ruleImpl) GetPriority() int {
	return rule.priority
}

func (rule *ruleImpl) SetPriority(priority int) {
	rule.priority = priority
}

func (rule *ruleImpl) String() string {
	str := ""
	str += "[Rule: (" + ") " + rule.name + "\n"
	//str += "[Rule: (" + strconv.Itoa(rule.id) + ") " + rule.name + "\n"

	str += "\t[Conditions:\n"
	for _, cond := range rule.conditions {
		str += "\t\t" + cond.String() + "\n"
	}
	str += "\t[Idrs:" + model.IdentifiersToString(rule.identifiers) + "]\n"
	return str
}

func (rule *ruleImpl) GetIdentifiers() []model.TupleType {
	return rule.identifiers
}

func (rule *ruleImpl) SetAction(actionFn model.ActionFunction) {
	rule.actionFn = actionFn
}

func (rule *ruleImpl) AddCondition(conditionName string, idrs []string, cFn model.ConditionEvaluator, ctx model.RuleContext) (err error) {
	typeDeps := []model.TupleType{}
	for _, idr := range idrs {
		aliasProp := strings.Split(string(idr), ".")
		alias := model.TupleType(aliasProp[0])

		if model.GetTupleDescriptor(model.TupleType(alias)) == nil {
			return errors.New("Tuple type not found " + string(alias))
		}

		typeDeps = append(typeDeps, alias)
		if len(aliasProp) == 2 { //specifically 2, else do not consider
			prop := aliasProp[1]

			td := model.GetTupleDescriptor(model.TupleType(alias))
			if prop != "none" && td.GetProperty(prop) == nil { //"none" is a special case
				return errors.New("TupleType property not found " + prop)
			}

			propMap, found := rule.deps[alias]
			if !found {
				propMap = map[string]bool{}
				rule.deps[alias] = propMap
			}
			propMap[prop] = true
		}
	}

	rule.addCond(conditionName, typeDeps, cFn, ctx, true)
	return err
}

func (rule *ruleImpl) GetDeps() map[model.TupleType]map[string]bool {
	return rule.deps
}
