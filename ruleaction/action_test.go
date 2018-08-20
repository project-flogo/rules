package ruleaction

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/TIBCOSoftware/bego/common/model"
)

func TestAction(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}
	for i := 1; i < 2; i++ {
		debit, _ := model.NewTupleWithKeyValues(model.TupleType("debitevent"), "Bob")
		debit.SetString(nil, "name", "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, "debit", fs)
		rs.Assert(nil, debit)
	}

	st1, _ := model.NewTupleWithKeyValues(model.TupleType("customerevent"), "Bob")
	st1.SetString(nil, "name", "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetDouble(nil, "balance", 1000)
	rs.Assert(nil, st1)
}

func TestActionTwo(t *testing.T) {
	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	st1, _ := model.NewTupleWithKeyValues(model.TupleType("customerevent"), "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetDouble(nil, "balance", 1000)
	rs.Assert(nil, st1)

	for i := 1; i < 2; i++ {
		debit, _ := model.NewTupleWithKeyValues(model.TupleType("debitevent"), "Bob")
		debit.SetDouble(nil, "debit", float64(i*100))
		rs.Assert(nil, debit)
	}
}

func TestActionTwoWithDep(t *testing.T) {
	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}
	st1, _ := model.NewTupleWithKeyValues(model.TupleType("customerevent"), "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetDouble(nil, "balance", 1000)
	rs.Assert(nil, st1)

	for i := 1; i < 3; i++ {
		debit, _ := model.NewTupleWithKeyValues(model.TupleType("debitevent"), "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, "debit", fs)
		rs.Assert(nil, debit)
	}

}

func TestTupleTTL(t *testing.T) {
	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}
	st1, _ := model.NewTupleWithKeyValues(model.TupleType("customerevent"), "Bob")
	st1.SetString(nil, "status", "active")
	st1.SetDouble(nil, "balance", 1000)
	rs.Assert(nil, st1)

	time.Sleep(time.Second * 6)

	for i := 1; i < 3; i++ {
		debit, _ := model.NewTupleWithKeyValues(model.TupleType("debitevent"), "Bob")
		fs := strconv.FormatFloat(float64(i*100), 'E', -1, 32)
		debit.SetString(nil, "debit", fs)
		rs.Assert(nil, debit)
	}
}

func TestActionTimeout(t *testing.T) {
	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	pt, _ := model.NewTupleWithKeyValues(model.TupleType("packagetimeout"), "pkg1")
	rs.ScheduleAssert(nil, 5000, "myid", pt)

	time.Sleep(time.Minute)
}

func TestActionTimeoutCancel(t *testing.T) {
	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}
	pt, _ := model.NewTupleWithKeyValues(model.TupleType("packagetimeout"), "pkg1")
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
