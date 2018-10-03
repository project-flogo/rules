package tests

import (
	"testing"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"fmt"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/expression/expr"
	"github.com/TIBCOSoftware/flogo-lib/core/mapper/exprmapper/funcexprtype"
	"strings"
	"github.com/project-flogo/rules/ruleapi"
)

func Test_1_Con(t *testing.T) {

	createRuleSession()

	//a1, _ := data.NewAttribute("name", data.TypeString, "n1")
	//a2, _ := data.NewAttribute("age2", data.TypeInteger, 10)
	//a3, _ := data.NewAttribute("age", data.TypeInteger, 10)
	//
	//attr := []*data.Attribute{a1, a2, a3}
	//
	//simpleScope := data.NewSimpleScope(attr, nil)
	//
	////fmt.Printf("bc index %d\n", strings.Index("abc", "bc"))
	//
	//r := data.GetBasicResolver()
	//
	//e, err := expression.ParseExpression(`($.age == $.age2) && ($.age2 == 11)`)
	//exprn := e.(*expr.Expression)
	//getRefs(exprn)
	//if err == nil {
	//	v, err2 := e.EvalWithScope(simpleScope, r)
	//	if err2 != nil {
	//		t.Error(err2)
	//	} else {
	//		t.Logf("Value: %t\n", v)
	//	}
	//} else {
	//	t.Error(err)
	//}

	r1 := ruleapi.NewRule("r1")
	err := r1.AddExprCondition("c1", "($t1.p1 == $t2.p2) && ($t1.p3 > $t2.p2)", nil)
	if err != nil {
		t.Errorf("error: %s\n", err)
	}


}

func getRefs(e *expr.Expression) (map[string]bool, error) {
	refs := make(map[string]bool)
	err := getRefRecursively(e, refs)
	return refs, err
}

func getRefRecursively (e *expr.Expression, refs map[string]bool) (error) {

	if e == nil {
		return nil
	}
	err := getRefsInternal(e.Left, refs)
	if err != nil {
		return err
	}
	err = getRefsInternal(e.Right, refs)
	if err != nil {
		return err
	}
	return nil
}

func getRefsInternal(e *expr.Expression, refs map[string]bool) (error) {
	if e.Type == funcexprtype.EXPRESSION {
		getRefRecursively(e, refs)
	} else if e.Type == funcexprtype.ARRAYREF {
		value := e.Value.(string)

		split := strings.Split(value, ".")
		if split != nil && len(split) != 2 {
			return fmt.Errorf("Invalid tokens [%s]", value)
		}

		refs[value] = true
	}
	return nil
}

type RuleExpressionResolver struct {

}

func (r *RuleExpressionResolver) Resolve(toResolve string, scope data.Scope) (value interface{}, err error) {
	fmt.Printf("Resolve: [%s]\n", toResolve)
	return nil, nil
}

