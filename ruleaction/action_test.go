package ruleaction

import (
	"github.com/TIBCOSoftware/bego/ruleapi"
	"testing"
	"github.com/TIBCOSoftware/bego/common/model"
	"strconv"
)


func TestAction (t *testing.T) {
	rs := ruleapi.NewRuleSession()

	loadRules(rs)

	st1 := model.NewStreamTuple(model.TupleTypeAlias("customerevent"))
	st1.SetString (nil,"name", "Bob")
	st1.SetString (nil,"status", "active")

	st1.SetFloat (nil,"balance", 1000)
	rs.Assert(nil, st1)

	for i := 1; i < 2  ; i++ {
		debit := model.NewStreamTuple(model.TupleTypeAlias("debitevent"))
		debit.SetString(nil,"name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil,"debit", fs)
		rs.Assert(nil, debit)
	}

}