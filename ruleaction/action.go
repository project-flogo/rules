package ruleaction

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
	"github.com/TIBCOSoftware/flogo-lib/core/action"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"github.com/TIBCOSoftware/flogo-lib/core/trigger"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"strconv"
	//"io/ioutil"
	//"log"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	ACTION_REF = "github.com/TIBCOSoftware/bego/ruleaction"
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

	ruleAction.rs.RegisterTupleDescriptors(string(tdss))

	loadPkgRulesWithDeps(ruleAction.rs)

	//loadRules(ruleAction.rs)

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
	//fmt.Printf("Received event from stream source [%s]", h.Name)

	streamSrc := model.TupleTypeAlias(h.Name)

	streamTuple := model.NewStreamTuple(streamSrc) //n1 -> will be replaced by contextual information coming in the data

	queryParams := inputs["queryParams"].Value().(map[string]string)

	//map input data into stream tuples, only string. ignore the rest for now
	for key, value := range queryParams {
		//fmt.Printf("[%s]\n", "a")
		if key == "balance" {
			f, _ := strconv.ParseFloat(value, 64)
			streamTuple.SetFloat(ctx, a.rs, key, f)
		} else {
			streamTuple.SetString(ctx, a.rs, key, value)
		}
	}

	a.rs.Assert(ctx, streamTuple)
	//fmt.Printf("[%s]\n", "b")
	return nil, nil
}

//func loadRules(rs model.RuleSession) {
//
//	rule := ruleapi.NewRule("customer-event")
//	rule.AddCondition("customer", []model.TupleTypeAlias{"customerevent"}, truecondition) // check for name "Bob" in n1
//	rule.SetAction(customerAction)
//	rule.SetPriority(1)
//	rs.AddRule(rule)
//	fmt.Printf("Rule added: [%s]\n", rule.GetName())
//
//	rule2 := ruleapi.NewRule("debit-event")
//	rule2.AddCondition("debitevent", []model.TupleTypeAlias{"debitevent"}, truecondition)
//	rule2.SetAction(debitEvent)
//	rule2.SetPriority(2)
//	rs.AddRule(rule2)
//	fmt.Printf("Rule added: [%s]\n", rule2.GetName())
//
//	rule3 := ruleapi.NewRule("customer-debit")
//	rule3.AddCondition("customerdebit", []model.TupleTypeAlias{"debitevent", "customerevent"}, customerdebitjoincondition)
//	rule3.SetAction(debitAction)
//	rule3.SetPriority(3)
//	rs.AddRule(rule3)
//	fmt.Printf("Rule added: [%s]\n", rule3.GetName())
//
//	rule4 := ruleapi.NewRule("check-balance")
//	rule4.AddCondition("customerdebit", []model.TupleTypeAlias{"customerevent"}, checkBalance) // check for name "Bob" in n1
//	rule4.SetAction(balanceAlert)
//	rule4.SetPriority(-1)
//	rs.AddRule(rule4)
//	fmt.Printf("Rule added: [%s]\n", rule4.GetName())
//
//
//	//rule5 := ruleapi.NewRule("timeout-rule")
//	//rule5.AddCondition("packagetimeout", []model.TupleTypeAlias{"packagetimeout"}, truecondition) // check for name "Bob" in n1
//	//rule5.SetAction(packageTimeout)
//	//rule5.SetPriority(-1)
//	//rs.AddRule(rule5)
//	//fmt.Printf("Rule added: [%s]\n", rule5.GetName())
//	//
//
//
//}
func loadRulesWithDeps(rs model.RuleSession) {

	rule := ruleapi.NewRule("customer-event")
	rule.AddConditionWithDependency("customer", []string{"customerevent.none"}, truecondition) // check for name "Bob" in n1
	rule.SetAction(customerAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	rule2 := ruleapi.NewRule("debit-event")
	rule2.AddConditionWithDependency("debitevent", []string{"debitevent.none"}, truecondition)
	rule2.SetAction(debitEvent)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	rule3 := ruleapi.NewRule("customer-debit")
	rule3.AddConditionWithDependency("customerdebit", []string{"debitevent.name", "customerevent.name"}, customerdebitjoincondition)
	rule3.SetAction(debitAction)
	rule3.SetPriority(3)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	//rule4 := ruleapi.NewRule("check-balance")
	//rule4.AddConditionWithDependency("customerdebit", []string{"customerevent.balance", "customerevent.status"}, checkBalance) // check for name "Bob" in n1
	//rule4.SetAction(balanceAlert)
	//rule4.SetPriority(-1)
	//rs.AddRule(rule4)
	//fmt.Printf("Rule added: [%s]\n", rule4.GetName())

	//rule5 := ruleapi.NewRule("timeout-rule")
	//rule5.AddCondition("packagetimeout", []model.TupleTypeAlias{"packagetimeout"}, truecondition) // check for name "Bob" in n1
	//rule5.SetAction(packageTimeout)
	//rule5.SetPriority(-1)
	//rs.AddRule(rule5)
	//fmt.Printf("Rule added: [%s]\n", rule5.GetName())
	//

}

func truecondition(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	return true
}
func customerAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	streamTuple := tuples["customerevent"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := streamTuple.GetString("name")
		fmt.Printf("Received a customer event with customer name [%s]\n", name)
	}
}
func debitEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	streamTuple := tuples["debitevent"]
	if streamTuple == nil {
		fmt.Println("Should not get a nil tuple in FilterCondition! This is an error")
	} else {
		name := streamTuple.GetString("name")
		amount := streamTuple.GetString("debit")

		fmt.Printf("Received a debit event for customer [%s], amount [%s]\n", name, amount)
	}
}

func customerdebitjoincondition(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {

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

func debitAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	//fmt.Printf("Rule fired: [%s]\n", ruleName)
	customerTuple := tuples["customerevent"].(model.MutableStreamTuple)
	debitTuple := tuples["debitevent"]
	dbt := debitTuple.GetString("debit")
	debitAmt, _ := strconv.ParseFloat(dbt, 64)
	currBal := customerTuple.GetFloat("balance")
	if customerTuple.GetString("status") == "active" {
		customerTuple.SetFloat(ctx, rs, "balance", customerTuple.GetFloat("balance")-debitAmt)
	}
	fmt.Printf("Customer [%s], Balance [%f], Debit [%f], NewBalance [%f]\n", customerTuple.GetString("name"), currBal, debitAmt, customerTuple.GetFloat("balance"))
}

//func checkBalance(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
//	customerTuple := tuples["customerevent"]
//	if customerTuple == nil {
//		fmt.Println("Should not get a nil tuple here! This is an error")
//		return false
//	}
//	balance := customerTuple.GetFloat("balance")
//	issuspended := customerTuple.GetString("status")
//
//	return balance <= 0 && issuspended == "active"
//}
//
//
//func balanceAlert(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
//	//fmt.Printf("Rule fired: [%s]\n", ruleName)
//	customerTuple := tuples["customerevent"].(model.MutableStreamTuple)
//	fmt.Printf("**** Account Suspended *** Customer balance is 0 or negative ! [%s], Balance [%f]\n", customerTuple.GetString("name"), customerTuple.GetFloat("balance"))
//	customerTuple.SetString(ctx,rs,"status", "suspended")
//}
//
//func packageTimeout(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
//	packageTimeout := tuples["packagetimeout"]
//	if packageTimeout == nil {
//		fmt.Println("Should not get a nil tuple here! This is an error")
//		return
//	}
//	fmt.Printf("Received a package timeout event for package-id [%s]\n", packageTimeout.GetString("packageid"))
//}

func createRuleSessionAndRules() model.RuleSession {
	rs := ruleapi.GetOrCreateRuleSession("asession")

	tupleDescFileAbsPath := getAbsPathForResource("src/github.com/TIBCOSoftware/bego/common/model/tupledescriptor.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Tuple descriptors: [%s]\n", string(dat))
	rs.RegisterTupleDescriptors(string(dat))
	//loadRules(rs)
	return rs
}

func createRuleSessionAndRulesWD() model.RuleSession {
	rs := ruleapi.GetOrCreateRuleSession("asession")

	tupleDescFileAbsPath := getAbsPathForResource("src/github.com/TIBCOSoftware/bego/common/model/tupledescriptor.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Printf("Tuple descriptors: [%s]\n", string(dat))
	rs.RegisterTupleDescriptors(string(dat))
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

//func loadPkgRules(rs model.RuleSession) {
//
//	//handle a package event, create a package in the packageAction
//	rule := ruleapi.NewRule("packageevent")
//	rule.AddCondition("packageevent", []model.TupleTypeAlias{"packageevent"}, truecondition)
//	rule.SetAction(packageeventAction)
//	rule.SetPriority(1)
//	rs.AddRule(rule)
//	fmt.Printf("Rule added: [%s]\n", rule.GetName())
//
//	//handle a package, print package details in the packageAction
//	rule1:= ruleapi.NewRule("package")
//	rule1.AddCondition("packageevent1", []model.TupleTypeAlias{"package"}, packageCondition)
//	rule1.SetAction(packageAction)
//	rule1.SetPriority(2)
//	rs.AddRule(rule1)
//	fmt.Printf("Rule added: [%s]\n", rule1.GetName())
//
//	//handle a scan event, see if there is matching package if so, do necessary things such as set off a timer
//	//for the next destination, etc in the scaneventAction
//	rule2 := ruleapi.NewRule("scanevent")
//	rule2.AddCondition("scanevent", []model.TupleTypeAlias{"package", "scanevent"}, scaneventCondition)
//	rule2.SetAction(scaneventAction)
//	rule2.SetPriority(2)
//	rs.AddRule(rule2)
//	fmt.Printf("Rule added: [%s]\n", rule2.GetName())
//
//	//handle a timeout event, triggered by scaneventAction, mark the package as delayed in scantimeoutAction
//	rule3 := ruleapi.NewRule("scantimeout")
//	rule3.AddCondition("packageevent", []model.TupleTypeAlias{"package", "scantimeout"}, scantimeoutCondition)
//	rule3.SetAction(scantimeoutAction)
//	rule3.SetPriority(1)
//	rs.AddRule(rule3)
//	fmt.Printf("Rule added: [%s]\n", rule3.GetName())
//
//	//notify when a package is marked as delayed, print as such in the packagedelayedAction
//	rule4 := ruleapi.NewRule("packagedelayed")
//	rule4.AddCondition("packageevent", []model.TupleTypeAlias{"package"}, packageDelayedCheck)
//	rule4.SetAction(packagedelayedAction)
//	rule4.SetPriority(1)
//	rs.AddRule(rule4)
//	fmt.Printf("Rule added: [%s]\n", rule4.GetName())
//}

func loadPkgRulesWithDeps(rs model.RuleSession) {

	//handle a package event, create a package in the packageAction
	rule := ruleapi.NewRule("packageevent")
	rule.AddConditionWithDependency("truecondition", []string{"packageevent.none"}, truecondition)
	rule.SetAction(packageeventAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	//handle a package, print package details in the packageAction
	rule1 := ruleapi.NewRule("package")
	rule1.AddConditionWithDependency("packageCondition", []string{"package.none"}, packageCondition)
	rule1.SetAction(packageAction)
	rule1.SetPriority(2)
	rs.AddRule(rule1)
	fmt.Printf("Rule added: [%s]\n", rule1.GetName())

	//handle a scan event, see if there is matching package if so, do necessary things such as set off a timer
	//for the next destination, etc in the scaneventAction
	rule2 := ruleapi.NewRule("scanevent")
	rule2.AddConditionWithDependency("scaneventCondition", []string{"package.packageid", "scanevent.packageid", "package.curr", "package.next"}, scaneventCondition)
	rule2.SetAction(scaneventAction)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//handle a timeout event, triggered by scaneventAction, mark the package as delayed in scantimeoutAction
	rule3 := ruleapi.NewRule("scantimeout")
	rule3.AddConditionWithDependency("scantimeoutCondition", []string{"package.packageid", "scantimeout.packageid"}, scantimeoutCondition)
	rule3.SetAction(scantimeoutAction)
	rule3.SetPriority(1)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	//notify when a package is marked as delayed, print as such in the packagedelayedAction
	rule4 := ruleapi.NewRule("packagedelayed")
	rule4.AddConditionWithDependency("packageDelayedCheck", []string{"package.status"}, packageDelayedCheck)
	rule4.SetAction(packagedelayedAction)
	rule4.SetPriority(1)
	rs.AddRule(rule4)
	fmt.Printf("Rule added: [%s]\n", rule4.GetName())
}

func packageeventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {

	pkgEvent := tuples["packageevent"]
	pkgid := pkgEvent.GetString("packageid")
	fmt.Printf("Received a new package asserting package id[%s]\n", pkgid)

	//assert a package
	pkg := model.NewStreamTuple(model.TupleTypeAlias("package"))
	pkg.SetString(ctx, rs, "packageid", pkgEvent.GetString("packageid"))
	pkg.SetString(ctx, rs, "curr", "start")
	pkg.SetString(ctx, rs, "next", pkgEvent.GetString("next"))
	pkg.SetString(ctx, rs, "status", "normal")
	pkg.SetString(ctx, rs, "isnew", pkgEvent.GetString("isnew"))

	rs.Assert(ctx, pkg)
}

func scaneventCondition(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	scanevent := tuples["scanevent"]
	pkg := tuples["package"]

	if scanevent == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	return scanevent.GetString("packageid") == pkg.GetString("packageid") &&
		scanevent.GetString("curr") == pkg.GetString("next")
}

func scaneventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	scanevent := tuples["scanevent"]

	pkg := tuples["package"].(model.MutableStreamTuple)
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

	scantmout := model.NewStreamTuple(model.TupleTypeAlias("scantimeout"))
	scantmout.SetString(ctx, rs, "packageid", pkgid)
	scantmout.SetString(ctx, rs, "next", snext)

	//cancel a previous timeout if set, since we got a scan event for the package's next destination
	prevtmoutid := pkgid + scurr
	rs.CancelDelayedAssert(ctx, prevtmoutid)

	//start the timer only if this scanevent says that its not "done", so there are more destinations
	if snext != "done" {
		tmoutid := pkgid + snext
		rs.DelayedAssert(ctx, uint64(eta*1000), tmoutid, scantmout) //start the timer here
	}
	pkg.SetString(ctx, rs, "curr", scurr)
	pkg.SetString(ctx, rs, "next", snext)

}

func scantimeoutCondition(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	scantimeout := tuples["scantimeout"]
	pkg := tuples["package"]

	if scantimeout == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	return scantimeout.GetString("packageid") == pkg.GetString("packageid") &&
		scantimeout.GetString("next") == pkg.GetString("next")
}

func scantimeoutAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {

	pkg := tuples["package"].(model.MutableStreamTuple)

	pkgid := pkg.GetString("packageid")
	pcurr := pkg.GetString("curr")
	pnext := pkg.GetString("next")

	fmt.Printf("Package id[%s] : Scan for dest [%s] did not arrive by ETA. Package currently at [%s]\n",
		pkgid, pnext, pcurr)
	pkg.SetString(ctx, rs, "status", "delayed")
}

func packageCondition(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	pkg := tuples["package"]
	isnew := pkg.GetString("isnew")
	return isnew == "true"
}

func packageAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	pkg := tuples["package"].(model.MutableStreamTuple)
	pkgid := pkg.GetString("packageid")

	pcurr := pkg.GetString("curr")
	pnext := pkg.GetString("next")
	fmt.Printf("Received a new package id[%s], current loc [%s], next loc [%s]\n", pkgid, pcurr, pnext)
	pkg.SetString(ctx, rs, "isnew", "false")
}

func packageDelayedCheck(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	pkg := tuples["package"]
	status := pkg.GetString("status")
	return status == "delayed"
}

func packagedelayedAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {
	pkg := tuples["package"].(model.MutableStreamTuple)
	pkgid := pkg.GetString("packageid")

	fmt.Printf("Package is now delayed id[%s]\n", pkgid)
}
