package ruleaction

import (
	"fmt"
	"github.com/tibmatt/bego/common/model"
	"strconv"
	"testing"
	"time"
)

func TestAction(t *testing.T) {
	rs := createRuleSessionAndRules()

	for i := 1; i < 2; i++ {
		debit := model.NewTuple(model.TupleType("debitevent"))
		debit.SetString(nil, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, "debit", fs)
		rs.Assert(nil, debit)
	}

	st1 := model.NewTuple(model.TupleType("customerevent"))
	st1.SetString(nil, "name", "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetFloat(nil, "balance", 1000)
	rs.Assert(nil, st1)
}

func TestActionTwo(t *testing.T) {
	rs := createRuleSessionAndRules()

	st1 := model.NewTuple(model.TupleType("customerevent"))
	st1.SetString(nil, "name", "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetFloat(nil, "balance", 1000)
	rs.Assert(nil, st1)

	for i := 1; i < 2; i++ {
		debit := model.NewTuple(model.TupleType("debitevent"))
		debit.SetString(nil, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, "debit", fs)
		rs.Assert(nil, debit)
	}

}

func TestActionTwoWithDep(t *testing.T) {
	rs := createRuleSessionAndRulesWD()

	st1 := model.NewTuple(model.TupleType("customerevent"))
	st1.SetString(nil, "name", "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetFloat(nil, "balance", 1000)
	rs.Assert(nil, st1)

	for i := 1; i < 3; i++ {
		debit := model.NewTuple(model.TupleType("debitevent"))
		debit.SetString(nil, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, "debit", fs)
		rs.Assert(nil, debit)
	}

}

func TestTupleTTL(t *testing.T) {
	rs := createRuleSessionAndRules()

	st1 := model.NewTuple(model.TupleType("customerevent"))
	st1.SetString(nil, "name", "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetFloat(nil, "balance", 1000)
	rs.Assert(nil, st1)

	time.Sleep(time.Second * 6)

	for i := 1; i < 3; i++ {
		debit := model.NewTuple(model.TupleType("debitevent"))
		debit.SetString(nil, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, "debit", fs)
		rs.Assert(nil, debit)
	}

}

func TestActionTimeout(t *testing.T) {
	rs := createRuleSessionAndRules()

	pt := model.NewTuple(model.TupleType("packagetimeout"))
	pt.SetString(nil, "packageid", "pkg1")
	rs.ScheduleAssert(nil, 5000, "myid", pt)

	time.Sleep(time.Minute)

}

func TestActionTimeoutCancel(t *testing.T) {
	rs := createRuleSessionAndRules()

	pt := model.NewTuple(model.TupleType("packagetimeout"))
	pt.SetString(nil, "packageid", "pkg1")
	rs.ScheduleAssert(nil, 1000, "myid", pt)

	rs.CancelScheduledAssert(nil, "myid")

	time.Sleep(time.Minute)

}

func TestActionBasicTimer(t *testing.T) {

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
	fmt.Printf("The task is to print myself [%d]\n", e.x)
}

func scheduleTask(period int, t Task) *time.Timer {
	tmr := time.AfterFunc(time.Second*time.Duration(period), func() {
		t.performOps()
	})
	return tmr
}
