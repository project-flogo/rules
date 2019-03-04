package ruleapi

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
	"github.com/project-flogo/rules/common/model"
	"reflect"
	"strings"
)

type exprConditionImpl struct {
	name        string
	rule        model.Rule
	identifiers []model.TupleType
	cExpr       string
	ctx         model.RuleContext
}

func newExprCondition(name string, rule model.Rule, identifiers []model.TupleType, cExpr string, ctx model.RuleContext) model.Condition {
	c := exprConditionImpl{}
	c.initExprConditionImpl(name, rule, identifiers, cExpr, ctx)
	return &c
}

func (cnd *exprConditionImpl) initExprConditionImpl(name string, rule model.Rule, identifiers []model.TupleType, cExpr string, ctx model.RuleContext) {
	cnd.name = name
	cnd.rule = rule
	cnd.identifiers = append(cnd.identifiers, identifiers...)
	cnd.cExpr = cExpr
	cnd.ctx = ctx
}

func (cnd *exprConditionImpl) GetIdentifiers() []model.TupleType {
	return cnd.identifiers
}
func (cnd *exprConditionImpl) GetContext() model.RuleContext {
	return cnd.ctx
}

func (cnd *exprConditionImpl) GetEvaluator() model.ConditionEvaluator {
	return nil
}

func (cnd *exprConditionImpl) String() string {
	return "[Condition: name:" + cnd.name + ", idrs: TODO]"
}

func (cnd *exprConditionImpl) GetName() string {
	return cnd.name
}

func (cnd *exprConditionImpl) GetRule() model.Rule {
	return cnd.rule
}
func (cnd *exprConditionImpl) GetTupleTypeAlias() []model.TupleType {
	return cnd.identifiers
}

func (cnd *exprConditionImpl) Evaluate(condName string, ruleNm string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) (bool, error) {
	result := false
	if cnd.cExpr != "" {
		e, err := expression.ParseExpression(cnd.cExpr)
		exprn := e.(*expr.Expression)
		if err != nil {
			return result, err
		}
		td := tuplePropertyResolver{}
		scope := tupleScope{tuples}
		res, err := exprn.EvalWithData(tuples, &scope, &td)
		if err != nil {
			return false, err
		} else if reflect.TypeOf(res).Kind() == reflect.Bool {
			result = res.(bool)
		}
	}

	return result, nil
}

//////////////////////////////////////////////////////////
type tupleScope struct {
	tuples map[model.TupleType]model.Tuple
}

func (ts *tupleScope) GetAttr(name string) (attr *data.Attribute, exists bool) {
	return nil, false
}

// SetAttrValue sets the value of the specified attribute
func (ts *tupleScope) SetAttrValue(name string, value interface{}) error {
	return nil
}

///////////////////////////////////////////////////////////
type tuplePropertyResolver struct {
}

func (t *tuplePropertyResolver) Resolve(toResolve string, scope data.Scope) (value interface{}, err error) {

	toResolve = toResolve[1:]
	aliasAndProp := strings.Split(toResolve, ".")

	ts := scope.(*tupleScope)
	var v interface{}
	if ts != nil {
		tuple := ts.tuples[model.TupleType(aliasAndProp[0])].(model.Tuple)
		if tuple != nil {

			p := tuple.GetTupleDescriptor().GetProperty(aliasAndProp[1])
			switch p.PropType {
			case data.TypeString:
				v, err = tuple.GetString(aliasAndProp[1])
			case data.TypeInteger:
				v, err = tuple.GetInt(aliasAndProp[1])
			case data.TypeLong:
				v, err = tuple.GetLong(aliasAndProp[1])
			case data.TypeDouble:
				v, err = tuple.GetDouble(aliasAndProp[1])
			case data.TypeBoolean:
				v, err = tuple.GetBool(aliasAndProp[1])
			}
		}
	}
	return v, err
}
