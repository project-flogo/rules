package ruleaction

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
)

// Action ref to register the action factory
const (
	ActionRef = "github.com/TIBCOSoftware/bego/ruleaction"
)

// RuleAction wraps RuleSession
type RuleAction struct {
	rs model.RuleSession
}

// ActionFactory wrapper to register with the action
type ActionFactory struct {
}

//todo fix this
var metadata = &action.Metadata{ID: ActionRef, Async: false}

func init() {
	action.RegisterFactory(ActionRef, &ActionFactory{})
}

// Init implements action.Factory.Init
func (ff *ActionFactory) Init() error {
	return nil
}

// ActionData maintains Tuple descriptor details
type ActionData struct {
	Tds []model.TupleDescriptor
}

// New implements action.Factory.New
func (ff *ActionFactory) New(config *action.Config) (action.Action, error) {

	ruleAction := &RuleAction{}

	ruleAction.rs, _ = ruleapi.GetOrCreateRuleSession("flogosession")

	actionData := ActionData{}
	actionData.Tds = []model.TupleDescriptor{}

	err := json.Unmarshal(config.Data, &actionData)
	if err != nil {
		return nil, fmt.Errorf("failed to read rule action data '%s' error '%s'", config.Id, err.Error())
	}
	tdss, _ := json.Marshal(&actionData.Tds)
	fmt.Printf("**ACTION DATA: [%s]\n**", string(tdss))

	model.RegisterTupleDescriptors(string(tdss))

	loadPkgRulesWithDeps(ruleAction.rs)

	return ruleAction, nil
}

// Metadata get the Action's metadata
func (a *RuleAction) Metadata() *action.Metadata {
	return metadata
}

// IOMetadata get the Action's IO metadata
func (a *RuleAction) IOMetadata() *data.IOMetadata {
	return nil
}

// Run implements action.Action.Run
func (a *RuleAction) Run(ctx context.Context, inputs map[string]*data.Attribute) (map[string]*data.Attribute, error) {

	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("Unhandled Error executing rule action \n")

			// todo: useful for debugging
			logger.Debugf("StackTrace: %s", debug.Stack())
		}

	}()

	h, _ok := trigger.HandlerFromContext(ctx)
	if !_ok {
		return nil, nil
	}
	//fmt.Printf("Received event from tuple source [%s]", h.Name)

	tupleType := model.TupleType(h.Name)
	queryParams := inputs["queryParams"].Value().(map[string]string)

	tuple, _ := model.NewTupleFromStringMap(tupleType, queryParams) //n1 -> will be replaced by contextual information coming in the data

	//map input data into tuples, only string. ignore the rest for now
	//for key, value := range queryParams {
	//	//fmt.Printf("[%s]\n", "a")
	//	if key == "balance" {
	//		f, _ := strconv.ParseFloat(value, 64)
	//		tuple.SetDouble(ctx, key, f)
	//	} else {
	//		tuple.SetString(ctx, key, value)
	//	}
	//}
	//tuple.SetValue(ctx, queryParams)

	a.rs.Assert(ctx, tuple)
	//fmt.Printf("[%s]\n", "b")
	return nil, nil
}

func loadRulesWithDeps(rs model.RuleSession) {

	rule := ruleapi.NewRule("customer-event")
	rule.AddCondition("customer", []model.TupleType{"customerevent.none"}, truecondition, nil) // check for name "Bob" in n1
	rule.SetAction(customerAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	rule2 := ruleapi.NewRule("debit-event")
	rule2.AddCondition("debitevent", []model.TupleType{"debitevent.none"}, truecondition, nil)
	rule2.SetAction(debitEvent)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	rule3 := ruleapi.NewRule("customer-debit")
	rule3.AddCondition("customerdebit", []model.TupleType{"debitevent.name", "customerevent.name"}, customerdebitjoincondition, nil)
	rule3.SetAction(debitAction)
	rule3.SetPriority(3)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())
}

func truecondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}
func customerAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	tuple := tuples["customerevent"]
	if tuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name, _ := tuple.GetString("name")
		fmt.Printf("Received a customer event with customer name [%s]\n", name)
	}
}
func debitEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	tuple := tuples["debitevent"]
	if tuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name, _ := tuple.GetString("name")
		amount, _ := tuple.GetString("debit")

		fmt.Printf("Received a debit event for customer [%s], amount [%s]\n", name, amount)
	}
}

func customerdebitjoincondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {

	customerTuple := tuples["customerevent"]
	debitTuple := tuples["debitevent"]

	if customerTuple == nil || debitTuple == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	custName, _ := customerTuple.GetString("name")
	acctName, _ := debitTuple.GetString("name")

	return custName == acctName

}

func debitAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	customerTuple := tuples["customerevent"].(model.MutableTuple)
	debitTuple := tuples["debitevent"]
	dbt, _ := debitTuple.GetString("debit")
	debitAmt, _ := strconv.ParseFloat(dbt, 64)
	currBal, _ := customerTuple.GetDouble("balance")
	st, _ := customerTuple.GetString("status")
	if st == "active" {
		customerTuple.SetDouble(ctx, "balance", currBal-debitAmt)
	}
	nm, _ := customerTuple.GetString("name")
	newBal, _ := customerTuple.GetDouble("balance")
	fmt.Printf("Customer [%s], Balance [%f], Debit [%f], NewBalance [%f]\n", nm, currBal, debitAmt, newBal)
}

func createRuleSessionAndRules() model.RuleSession {
	rs, _ := ruleapi.GetOrCreateRuleSession("asession")

	tupleDescFileAbsPath := getAbsPathForResource("src/github.com/TIBCOSoftware/bego/common/model/tupledescriptor.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	model.RegisterTupleDescriptors(string(dat))
	return rs
}

func createRuleSessionAndRulesWD() model.RuleSession {
	rs, _ := ruleapi.GetOrCreateRuleSession("asession")

	tupleDescFileAbsPath := getAbsPathForResource("src/github.com/TIBCOSoftware/bego/common/model/tupledescriptor.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	model.RegisterTupleDescriptors(string(dat))
	loadRulesWithDeps(rs)
	return rs
}

func getAbsPathForResource(resourcepath string) string {
	GOPATH := os.Getenv("GOPATH")
	fmt.Printf("path[%s]\n", GOPATH)
	paths := strings.Split(GOPATH, ":")
	for _, path := range paths {
		fmt.Printf("path[%s]\n", path)
		absPath := path + "/" + resourcepath
		_, err := os.Stat(absPath)
		if err == nil {
			return absPath
		}
	}
	return ""
}

func loadPkgRulesWithDeps(rs model.RuleSession) {

	//handle a package event, create a package in the packageAction
	rule := ruleapi.NewRule("packageevent")
	rule.AddCondition("truecondition", []model.TupleType{"packageevent.none"}, truecondition, nil)
	rule.SetAction(packageeventAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	//handle a package, print package details in the packageAction
	rule1 := ruleapi.NewRule("package")
	rule1.AddCondition("packageCondition", []model.TupleType{"package.none"}, packageCondition, nil)
	rule1.SetAction(packageAction)
	rule1.SetPriority(2)
	rs.AddRule(rule1)
	fmt.Printf("Rule added: [%s]\n", rule1.GetName())

	//handle a scan event, see if there is matching package if so, do necessary things such as set off a timer
	//for the next destination, etc in the scaneventAction
	rule2 := ruleapi.NewRule("scanevent")
	rule2.AddCondition("scaneventCondition", []model.TupleType{"package.packageid", "scanevent.packageid", "package.curr", "package.next"}, scaneventCondition, nil)
	rule2.SetAction(scaneventAction)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//handle a timeout event, triggered by scaneventAction, mark the package as delayed in scantimeoutAction
	rule3 := ruleapi.NewRule("scantimeout")
	rule3.AddCondition("scantimeoutCondition", []model.TupleType{"package.packageid", "scantimeout.packageid"}, scantimeoutCondition, nil)
	rule3.SetAction(scantimeoutAction)
	rule3.SetPriority(1)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	//notify when a package is marked as delayed, print as such in the packagedelayedAction
	rule4 := ruleapi.NewRule("packagedelayed")
	rule4.AddCondition("packageDelayedCheck", []model.TupleType{"package.status"}, packageDelayedCheck, nil)
	rule4.SetAction(packagedelayedAction)
	rule4.SetPriority(1)
	rs.AddRule(rule4)
	fmt.Printf("Rule added: [%s]\n", rule4.GetName())
}

func packageeventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {

	pkgEvent := tuples["packageevent"]
	pkgid, _ := pkgEvent.GetString("packageid")
	fmt.Printf("Received a new package asserting package id[%s]\n", pkgid)

	//assert a package
	pkg, _ := model.NewTuple(model.TupleType("package"))
	pkgID, _ := pkgEvent.GetString("packageid")
	nxt, _ := pkgEvent.GetString("next")
	pkg.SetString(ctx, "packageid", pkgID)
	pkg.SetString(ctx, "curr", "start")
	pkg.SetString(ctx, "next", nxt)
	pkg.SetString(ctx, "status", "normal")

	rs.Assert(ctx, pkg)
}

func scaneventCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	scanevent := tuples["scanevent"]
	pkg := tuples["package"]

	if scanevent == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	pkgID, _ := scanevent.GetString("packageid")
	pkgID2, _ := pkg.GetString("packageid")
	curr, _ := scanevent.GetString("curr")
	nxt, _ := pkg.GetString("next")
	return pkgID == pkgID2 && curr == nxt
}

func scaneventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	scanevent := tuples["scanevent"]

	pkg := tuples["package"].(model.MutableTuple)
	pkgid, _ := pkg.GetString("packageid")

	scurr, _ := scanevent.GetString("curr")
	snext, _ := scanevent.GetString("next")
	fmt.Printf("Received a new scan event for package id[%s], current loc [%s], next loc [%s]\n", pkgid, scurr, snext)

	if scanevent == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return
	}

	etaS, _ := scanevent.GetString("eta")
	eta, _ := strconv.Atoi(etaS)

	scantmout, _ := model.NewTuple(model.TupleType("scantimeout"))
	scantmout.SetString(ctx, "packageid", pkgid)
	scantmout.SetString(ctx, "next", snext)

	//cancel a previous timeout if set, since we got a scan event for the package's next destination
	prevtmoutid := pkgid + scurr
	rs.CancelScheduledAssert(ctx, prevtmoutid)

	//start the timer only if this scanevent says that its not "done", so there are more destinations
	if snext != "done" {
		tmoutid := pkgid + snext
		rs.ScheduleAssert(ctx, uint64(eta*1000), tmoutid, scantmout) //start the timer here
	}
	pkg.SetString(ctx, "curr", scurr)
	pkg.SetString(ctx, "next", snext)

}

func scantimeoutCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	scantimeout := tuples["scantimeout"]
	pkg := tuples["package"]

	if scantimeout == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	pkgID, _ := scantimeout.GetString("packageid")
	pkgID2, _ := pkg.GetString("packageid")
	nxt, _ := scantimeout.GetString("next")
	nxt2, _ := pkg.GetString("next")
	return pkgID == pkgID2 &&
		nxt == nxt2
}

func scantimeoutAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {

	pkg := tuples["package"].(model.MutableTuple)

	pkgid, _ := pkg.GetString("packageid")
	pcurr, _ := pkg.GetString("curr")
	pnext, _ := pkg.GetString("next")

	fmt.Printf("Package id[%s] : Scan for dest [%s] did not arrive by ETA. Package currently at [%s]\n",
		pkgid, pnext, pcurr)
	pkg.SetString(ctx, "status", "delayed")
}

func packageCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	isnew, _ := pkg.GetString("isnew")
	return isnew == "true"
}

func packageAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"].(model.MutableTuple)
	pkgid, _ := pkg.GetString("packageid")

	pcurr, _ := pkg.GetString("curr")
	pnext, _ := pkg.GetString("next")
	fmt.Printf("Received a new package id[%s], current loc [%s], next loc [%s]\n", pkgid, pcurr, pnext)
	pkg.SetString(ctx, "isnew", "false")
}

func packageDelayedCheck(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	status, _ := pkg.GetString("status")
	return status == "delayed"
}

func packagedelayedAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"].(model.MutableTuple)
	pkgid, _ := pkg.GetString("packageid")

	fmt.Printf("Package is now delayed id[%s]\n", pkgid)
}
