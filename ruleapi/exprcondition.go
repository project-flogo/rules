package ruleapi

import (
	"reflect"
	"strconv"

	"github.com/project-flogo/core/data/property"

	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/expression"
	"github.com/project-flogo/core/data/expression/script"
	"github.com/project-flogo/core/data/resolve"
	"github.com/project-flogo/rules/common/model"
)

var td tuplePropertyResolver
var resolver resolve.CompositeResolver
var factory expression.Factory

func init() {
	td = tuplePropertyResolver{}
	//resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{".": &td})
	resolver = resolve.NewCompositeResolver(map[string]resolve.Resolver{
		".":        &td,
		"env":      &resolve.EnvResolver{},
		"property": &property.Resolver{},
		"loop":     &resolve.LoopResolver{},
	})
	factory = script.NewExprFactory(resolver)
}

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
	if name == "" {
		cndIdx := len(rule.GetConditions()) + 1
		name = "c_" + strconv.Itoa(cndIdx)
	}
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
		//e, err := expression.ParseExpression(cnd.cExpr)
		exprn, err := factory.NewExpr(cnd.cExpr)
		if err != nil {
			return result, err
		}

		scope := tupleScope{tuples}
		res, err := exprn.Eval(&scope)
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

func (ts *tupleScope) GetValue(name string) (value interface{}, exists bool) {
	return false, true
}

func (ts *tupleScope) SetValue(name string, value interface{}) error {
	return nil
}

// SetAttrValue sets the value of the specified attribute
func (ts *tupleScope) SetAttrValue(name string, value interface{}) error {
	return nil
}

///////////////////////////////////////////////////////////
type tuplePropertyResolver struct {
}

//func (t *tuplePropertyResolver) Resolve(toResolve string, scope data.Scope) (value interface{}, err error) {
func (t *tuplePropertyResolver) Resolve(scope data.Scope, item string, field string) (interface{}, error) {
	//toResolve = toResolve[1:]
	//aliasAndProp := strings.Split(toResolve, ".")
	//
	//var v interface{}
	//if ts != nil {
	//	tuple := ts.tuples[model.TupleType(aliasAndProp[0])].(model.Tuple)
	//	if tuple != nil {
	//
	//		p := tuple.GetTupleDescriptor().GetProperty(aliasAndProp[1])
	//		switch p.PropType {
	//		case data.TypeString:
	//			v, err = tuple.GetString(aliasAndProp[1])
	//		case data.TypeInteger:
	//			v, err = tuple.GetInt(aliasAndProp[1])
	//		case data.TypeLong:
	//			v, err = tuple.GetLong(aliasAndProp[1])
	//		case data.TypeDouble:
	//			v, err = tuple.GetDouble(aliasAndProp[1])
	//		case data.TypeBoolean:
	//			v, err = tuple.GetBool(aliasAndProp[1])
	//		}
	//	}`
	//}
	//return v, err
	ts := scope.(*tupleScope)
	tuple := ts.tuples[model.TupleType(field)]
	m := tuple.GetMap()
	return m, nil

}

func (*tuplePropertyResolver) GetResolverInfo() *resolve.ResolverInfo {
	return resolve.NewResolverInfo(false, false)
}
