package ruleaction

import (
	"testing"
	"github.com/TIBCOSoftware/bego/common/model"
	"strconv"
	"fmt"
	"time"
)


func TestAction (t *testing.T) {
	rs := createRuleSessionAndRules()

	for i := 1; i < 2  ; i++ {
		debit := model.NewStreamTuple(model.TupleTypeAlias("debitevent"))
		debit.SetString(nil, rs,"name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, rs,"debit", fs)
		rs.Assert(nil, debit)
	}

	st1 := model.NewStreamTuple(model.TupleTypeAlias("customerevent"))
	st1.SetString (nil, rs,"name", "Bob")
	st1.SetString (nil, rs,"status", "active")
	st1.SetFloat (nil, rs,"balance", 1000)
	rs.Assert(nil, st1)
}

func TestActionTwo (t *testing.T) {
	rs := createRuleSessionAndRules()

	st1 := model.NewStreamTuple(model.TupleTypeAlias("customerevent"))
	st1.SetString(nil, rs, "name", "Bob")
	st1.SetString(nil, rs, "status", "active")
	st1.SetFloat(nil, rs, "balance", 1000)
	rs.Assert(nil, st1)

	for i := 1; i < 2; i++ {
		debit := model.NewStreamTuple(model.TupleTypeAlias("debitevent"))
		debit.SetString(nil, rs, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, rs, "debit", fs)
		rs.Assert(nil, debit)
	}

}

func TestActionTwoWithDep (t *testing.T) {
	rs := createRuleSessionAndRulesWD()

	st1 := model.NewStreamTuple(model.TupleTypeAlias("customerevent"))
	st1.SetString(nil, rs, "name", "Bob")
	st1.SetString(nil, rs, "status", "active")
	st1.SetFloat(nil, rs, "balance", 1000)
	rs.Assert(nil, st1)

	for i := 1; i < 3; i++ {
		debit := model.NewStreamTuple(model.TupleTypeAlias("debitevent"))
		debit.SetString(nil, rs, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, rs, "debit", fs)
		rs.Assert(nil, debit)
	}

}

func TestTupleTTL (t *testing.T) {
	rs := createRuleSessionAndRules()

	st1 := model.NewStreamTuple(model.TupleTypeAlias("customerevent"))
	st1.SetString(nil, rs, "name", "Bob")
	st1.SetString(nil, rs, "status", "active")
	st1.SetFloat(nil, rs, "balance", 1000)
	rs.Assert(nil, st1)

	time.Sleep(time.Second * 6)

	for i := 1; i < 3; i++ {
		debit := model.NewStreamTuple(model.TupleTypeAlias("debitevent"))
		debit.SetString(nil, rs, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, rs, "debit", fs)
		rs.Assert(nil, debit)
	}

}



func TestActionTimeout (t *testing.T) {
	rs := createRuleSessionAndRules()

	pt := model.NewStreamTuple(model.TupleTypeAlias("packagetimeout"))
	pt.SetString(nil, rs, "packageid", "pkg1")
	rs.DelayedAssert(nil, 5000,"myid", pt)

	time.Sleep(time.Minute)

}

func TestActionTimeoutCancel (t *testing.T) {
	rs := createRuleSessionAndRules()

	pt := model.NewStreamTuple(model.TupleTypeAlias("packagetimeout"))
	pt.SetString(nil, rs, "packageid", "pkg1")
	rs.DelayedAssert(nil, 1000,"myid", pt)

	rs.CancelDelayedAssert(nil, "myid")

	time.Sleep(time.Minute)

}


func TestActionBasicTimer (t *testing.T) {

	e := &EventTimer{151}

	scheduleTask(5, e)

	time.Sleep(time.Minute)
}

type Task interface {
	performOps()
}

type EventTimer struct {
	x int
}

func (e *EventTimer) performOps() {
	fmt.Printf ("The task is to print myself [%d]\n", e.x)
}

func scheduleTask (period int, t Task) *time.Timer {

	tmr := time.AfterFunc(time.Second * time.Duration(period), func() {
		t.performOps()
	})

	return tmr
}