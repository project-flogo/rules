package ruleapi

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)


type exprConditionImpl struct {
	name        string
	rule        model.Rule
	identifiers []model.TupleType
	cExpr         string
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

func (cnd *exprConditionImpl) Evaluate (condName string, ruleNm string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) (bool, error) {
	result := false
	if cnd.cExpr != "" {
		e, err := expression.ParseExpression(cnd.cExpr)
		exprn := e.(*expr.Expression)
		if err != nil {
			return result, err
		}
		td := TuplePropertyResolver{}
		res, err := exprn.EvalWithScope(nil, &td)

		result = res.(bool)
	}

	return result, nil
}

type TuplePropertyResolver struct {

}

func (t *TuplePropertyResolver) Resolve(toResolve string, scope data.Scope) (value interface{}, err error) {


	return nil, nil
}