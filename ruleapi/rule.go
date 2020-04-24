package ruleapi

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/project-flogo/rules/common/model"
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

func (rule *ruleImpl) AddIdrsToRule(idrs []model.TupleType) {
	for _, cidr := range idrs {
		//TODO: configure the rulesession
		if model.GetTupleDescriptor(cidr) == nil {
			return
		}
		if len(rule.identifiers) == 0 {
			rule.identifiers = append(rule.identifiers, cidr)
		} else {
			found := false
			for _, ridr := range rule.identifiers {
				if cidr == ridr {
					found = true
					break
				}
			}
			if !found {
				rule.identifiers = append(rule.identifiers, cidr)
			}
		}
	}
}

func (rule *ruleImpl) addExprCond(conditionName string, idrs []model.TupleType, cExpr string, ctx model.RuleContext) {
	condition := newExprCondition(conditionName, rule, idrs, cExpr, ctx)
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

func (rule *ruleImpl) AddCondition2(conditionName string, idrs []string, cFn model.ConditionEvaluator, ctx model.RuleContext) (err error) {
	typeDeps := []model.TupleType{}
	for _, idr := range idrs {
		aliasProp := strings.Split(string(idr), ".")
		alias := model.TupleType(aliasProp[0])

		if model.GetTupleDescriptor(model.TupleType(alias)) == nil {
			return fmt.Errorf("Tuple type not found [%s]", string(alias))
		}

		exists, _ := model.Contains(typeDeps, alias)
		if !exists {
			typeDeps = append(typeDeps, alias)
		}
		if len(aliasProp) == 2 { //specifically 2, else do not consider
			prop := aliasProp[1]

			td := model.GetTupleDescriptor(model.TupleType(alias))
			if prop != "none" && td.GetProperty(prop) == nil { //"none" is a special case
				return fmt.Errorf("TupleType property not found [%s]", prop)
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

func (rule *ruleImpl) AddCondition(conditionName string, idrs []string, cFn model.ConditionEvaluator, ctx model.RuleContext) (err error) {
	typeDeps, err := rule.addDeps(idrs)
	if err != nil {
		return err
	}
	rule.addCond(conditionName, typeDeps, cFn, ctx, true)
	return err
}

func (rule *ruleImpl) addDeps(idrs []string) ([]model.TupleType, error) {
	typeDeps := []model.TupleType{}
	for _, idr := range idrs {
		aliasProp := strings.Split(string(idr), ".")
		alias := model.TupleType(aliasProp[0])

		if model.GetTupleDescriptor(model.TupleType(alias)) == nil {
			return typeDeps, fmt.Errorf("Tuple type not found [%s]", string(alias))
		}

		exists, _ := model.Contains(typeDeps, alias)
		if !exists {
			typeDeps = append(typeDeps, alias)
		}
		if len(aliasProp) == 2 { //specifically 2, else do not consider
			prop := aliasProp[1]

			td := model.GetTupleDescriptor(model.TupleType(alias))
			if prop != "none" && td.GetProperty(prop) == nil { //"none" is a special case
				return typeDeps, fmt.Errorf("TupleType property not found [%s]", prop)
			}

			propMap, found := rule.deps[alias]
			if !found {
				propMap = map[string]bool{}
				rule.deps[alias] = propMap
			}
			propMap[prop] = true
		}
	}

	return typeDeps, nil
}

func (rule *ruleImpl) GetDeps() map[model.TupleType]map[string]bool {
	return rule.deps
}

func (rule *ruleImpl) AddExprCondition(conditionName string, cstr string, ctx model.RuleContext) error {

	//e, err := expression.ParseExpression(cstr)
	//if err != nil {
	//	return err
	//}
	//exprn := e.(*expr.Expression)
	//refs, err := getRefs(exprn)
	refs := getRefs(cstr)

	err := validateRefs(refs)
	if err != nil {
		return err
	}

	typeDeps, err := rule.addDeps(refs)
	if err != nil {
		return err
	}
	rule.addExprCond(conditionName, typeDeps, cstr, ctx)
	return nil

}

func validateRefs(refs []string) error {
	for _, ref := range refs {
		ref := strings.TrimPrefix(ref, "$.")
		vals := strings.Split(ref, ".")
		td := model.GetTupleDescriptor(model.TupleType(vals[0]))
		if td == nil {
			return fmt.Errorf("Invalid TupleType [%s]", vals[0])
		}
		prop := td.GetProperty(vals[1])
		if prop == nil {
			return fmt.Errorf("Property [%s] not found in TupleType [%s]", vals[1], vals[0])
		}
	}
	return nil
}

//
//func getRefs(e *expr.Expression) ([]string, error) {
//	refs := make(map[string]bool)
//	keys := []string{}
//
//	err := getRefRecursively(e, refs)
//	if err != nil {
//		return keys, err
//	}
//
//	for key := range refs {
//		keys = append (keys, key)
//	}
//	return keys, err
//}
//
//func getRefRecursively(e *expr.Expression, refs map[string]bool) error {
//
//	if e == nil {
//		return nil
//	}
//	err := getRefsInternal(e.Left, refs)
//	if err != nil {
//		return err
//	}
//	err = getRefsInternal(e.Right, refs)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func getRefsInternal(e *expr.Expression, refs map[string]bool) error {
//	if e.Type == funcexprtype.EXPRESSION {
//		getRefRecursively(e, refs)
//	} else if e.Type == funcexprtype.REF || e.Type == funcexprtype.ARRAYREF {
//		value := e.Value.(string)
//		if strings.Index(value, "$") == 0 {
//			value = value[1:len(value)]
//			split := strings.Split(value, ".")
//			if split != nil && len(split) != 2 {
//				return fmt.Errorf("Invalid tokens [%s]", value)
//			}
//
//			refs[value] = true
//		}
//	}
//	return nil
//}

func getRefs(cstr string) []string {
	keys2 := []string{}
	re := regexp.MustCompile(`\$\.(\w+\.\w+)`)
	keys := re.FindAllStringSubmatch(cstr, -1)
	for _, k := range keys {
		keys2 = append(keys2, k[1])
	}
	return keys2
}
