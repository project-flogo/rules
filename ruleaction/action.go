package ruleaction

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/tibmatt/bego/common/model"
	"github.com/tibmatt/bego/ruleapi"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	ACTION_REF = "github.com/tibmatt/bego/ruleaction"
)

type RuleAction struct {
	rs model.RuleSession
}
type ActionFactory struct {
}

//todo fix this
var metadata = &action.Metadata{ID: ACTION_REF, Async: false}

func init() {
	action.RegisterFactory(ACTION_REF, &ActionFactory{})
}

func (ff *ActionFactory) Init() error {
	return nil
}

type ActionData struct {
	Tds []model.TupleDescriptor
}

func (ff *ActionFactory) New(config *action.Config) (action.Action, error) {

	ruleAction := &RuleAction{}

	ruleAction.rs = ruleapi.GetOrCreateRuleSession("flogosession")

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

	tuple := model.NewTuple(tupleType) //n1 -> will be replaced by contextual information coming in the data

	queryParams := inputs["queryParams"].Value().(map[string]string)

	//map input data into tuples, only string. ignore the rest for now
	for key, value := range queryParams {
		//fmt.Printf("[%s]\n", "a")
		if key == "balance" {
			f, _ := strconv.ParseFloat(value, 64)
			tuple.SetFloat(ctx, key, f)
		} else {
			tuple.SetString(ctx, key, value)
		}
	}

	a.rs.Assert(ctx, tuple)
	//fmt.Printf("[%s]\n", "b")
	return nil, nil
}

func loadRulesWithDeps(rs model.RuleSession) {

	rule := ruleapi.NewRule("customer-event")
	rule.AddConditionWithDependency("customer", []string{"customerevent.none"}, truecondition, nil) // check for name "Bob" in n1
	rule.SetAction(customerAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	rule2 := ruleapi.NewRule("debit-event")
	rule2.AddConditionWithDependency("debitevent", []string{"debitevent.none"}, truecondition, nil)
	rule2.SetAction(debitEvent)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	rule3 := ruleapi.NewRule("customer-debit")
	rule3.AddConditionWithDependency("customerdebit", []string{"debitevent.name", "customerevent.name"}, customerdebitjoincondition, nil)
	rule3.SetAction(debitAction)
	rule3.SetPriority(3)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())
}

func truecondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.ConditionContext) bool {
	return true
}
func customerAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	tuple := tuples["customerevent"]
	if tuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := tuple.GetString("name")
		fmt.Printf("Received a customer event with customer name [%s]\n", name)
	}
}
func debitEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	tuple := tuples["debitevent"]
	if tuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := tuple.GetString("name")
		amount := tuple.GetString("debit")

		fmt.Printf("Received a debit event for customer [%s], amount [%s]\n", name, amount)
	}
}

func customerdebitjoincondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.ConditionContext) bool {

	customerTuple := tuples["customerevent"]
	debitTuple := tuples["debitevent"]

	if customerTuple == nil || debitTuple == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	custName := customerTuple.GetString("name")
	acctName := debitTuple.GetString("name")

	return custName == acctName

}

func debitAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	customerTuple := tuples["customerevent"].(model.MutableTuple)
	debitTuple := tuples["debitevent"]
	dbt := debitTuple.GetString("debit")
	debitAmt, _ := strconv.ParseFloat(dbt, 64)
	currBal := customerTuple.GetFloat("balance")
	if customerTuple.GetString("status") == "active" {
		customerTuple.SetFloat(ctx, "balance", customerTuple.GetFloat("balance")-debitAmt)
	}
	fmt.Printf("Customer [%s], Balance [%f], Debit [%f], NewBalance [%f]\n", customerTuple.GetString("name"), currBal, debitAmt, customerTuple.GetFloat("balance"))
}

func createRuleSessionAndRules() model.RuleSession {
	rs := ruleapi.GetOrCreateRuleSession("asession")

	tupleDescFileAbsPath := getAbsPathForResource("src/github.com/tibmatt/bego/common/model/tupledescriptor.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Tuple descriptors: [%s]\n", string(dat))
	model.RegisterTupleDescriptors(string(dat))
	return rs
}

func createRuleSessionAndRulesWD() model.RuleSession {
	rs := ruleapi.GetOrCreateRuleSession("asession")

	tupleDescFileAbsPath := getAbsPathForResource("src/github.com/tibmatt/bego/common/model/tupledescriptor.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Tuple descriptors: [%s]\n", string(dat))
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
	rule.AddConditionWithDependency("truecondition", []string{"packageevent.none"}, truecondition, nil)
	rule.SetAction(packageeventAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	//handle a package, print package details in the packageAction
	rule1 := ruleapi.NewRule("package")
	rule1.AddConditionWithDependency("packageCondition", []string{"package.none"}, packageCondition, nil)
	rule1.SetAction(packageAction)
	rule1.SetPriority(2)
	rs.AddRule(rule1)
	fmt.Printf("Rule added: [%s]\n", rule1.GetName())

	//handle a scan event, see if there is matching package if so, do necessary things such as set off a timer
	//for the next destination, etc in the scaneventAction
	rule2 := ruleapi.NewRule("scanevent")
	rule2.AddConditionWithDependency("scaneventCondition", []string{"package.packageid", "scanevent.packageid", "package.curr", "package.next"}, scaneventCondition, nil)
	rule2.SetAction(scaneventAction)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//handle a timeout event, triggered by scaneventAction, mark the package as delayed in scantimeoutAction
	rule3 := ruleapi.NewRule("scantimeout")
	rule3.AddConditionWithDependency("scantimeoutCondition", []string{"package.packageid", "scantimeout.packageid"}, scantimeoutCondition, nil)
	rule3.SetAction(scantimeoutAction)
	rule3.SetPriority(1)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	//notify when a package is marked as delayed, print as such in the packagedelayedAction
	rule4 := ruleapi.NewRule("packagedelayed")
	rule4.AddConditionWithDependency("packageDelayedCheck", []string{"package.status"}, packageDelayedCheck, nil)
	rule4.SetAction(packagedelayedAction)
	rule4.SetPriority(1)
	rs.AddRule(rule4)
	fmt.Printf("Rule added: [%s]\n", rule4.GetName())
}

func packageeventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {

	pkgEvent := tuples["packageevent"]
	pkgid := pkgEvent.GetString("packageid")
	fmt.Printf("Received a new package asserting package id[%s]\n", pkgid)

	//assert a package
	pkg := model.NewTuple(model.TupleType("package"))
	pkg.SetString(ctx, "packageid", pkgEvent.GetString("packageid"))
	pkg.SetString(ctx, "curr", "start")
	pkg.SetString(ctx, "next", pkgEvent.GetString("next"))
	pkg.SetString(ctx, "status", "normal")
	pkg.SetString(ctx, "isnew", pkgEvent.GetString("isnew"))

	rs.Assert(ctx, pkg)
}

func scaneventCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.ConditionContext) bool {
	scanevent := tuples["scanevent"]
	pkg := tuples["package"]

	if scanevent == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	return scanevent.GetString("packageid") == pkg.GetString("packageid") &&
		scanevent.GetString("curr") == pkg.GetString("next")
}

func scaneventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	scanevent := tuples["scanevent"]

	pkg := tuples["package"].(model.MutableTuple)
	pkgid := pkg.GetString("packageid")

	scurr := scanevent.GetString("curr")
	snext := scanevent.GetString("next")
	fmt.Printf("Received a new scan event for package id[%s], current loc [%s], next loc [%s]\n", pkgid, scurr, snext)

	if scanevent == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return
	}

	etaS := scanevent.GetString("eta")
	eta, _ := strconv.Atoi(etaS)

	scantmout := model.NewTuple(model.TupleType("scantimeout"))
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

func scantimeoutCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.ConditionContext) bool {
	scantimeout := tuples["scantimeout"]
	pkg := tuples["package"]

	if scantimeout == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	return scantimeout.GetString("packageid") == pkg.GetString("packageid") &&
		scantimeout.GetString("next") == pkg.GetString("next")
}

func scantimeoutAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {

	pkg := tuples["package"].(model.MutableTuple)

	pkgid := pkg.GetString("packageid")
	pcurr := pkg.GetString("curr")
	pnext := pkg.GetString("next")

	fmt.Printf("Package id[%s] : Scan for dest [%s] did not arrive by ETA. Package currently at [%s]\n",
		pkgid, pnext, pcurr)
	pkg.SetString(ctx, "status", "delayed")
}

func packageCondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.ConditionContext) bool {
	pkg := tuples["package"]
	isnew := pkg.GetString("isnew")
	return isnew == "true"
}

func packageAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	pkg := tuples["package"].(model.MutableTuple)
	pkgid := pkg.GetString("packageid")

	pcurr := pkg.GetString("curr")
	pnext := pkg.GetString("next")
	fmt.Printf("Received a new package id[%s], current loc [%s], next loc [%s]\n", pkgid, pcurr, pnext)
	pkg.SetString(ctx, "isnew", "false")
}

func packageDelayedCheck(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.ConditionContext) bool {
	pkg := tuples["package"]
	status := pkg.GetString("status")
	return status == "delayed"
}

func packagedelayedAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple) {
	pkg := tuples["package"].(model.MutableTuple)
	pkgid := pkg.GetString("packageid")

	fmt.Printf("Package is now delayed id[%s]\n", pkgid)
}
