package trackntrace

import (
	"fmt"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
	"github.com/project-flogo/rules/common"

	"testing"
	"time"
	"io/ioutil"
	"log"
	"strconv"
	"context"
)

func TestPkgFlowNormal(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, _ := model.NewTupleWithKeyValues("packageevent", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	err = rs.Assert(nil, pkgEvt)
	if err != nil {
		fmt.Printf("Error...[%s]\n", err)
		return
	}
	//time.Sleep(time.Second*20)
	scanEv, _ := model.NewTupleWithKeyValues("scanevent", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	err = scanEv.SetValue(nil, "eta", 10)
	if err != nil {
		fmt.Printf("[%s]\n", err)
	}

	rs.Assert(nil, scanEv)

	scanEv1, _ := model.NewTupleWithKeyValues("scanevent", "1")
	scanEv1.SetString(nil, "curr", "ny")
	scanEv1.SetString(nil, "next", "done")
	scanEv.SetString(nil, "next", "ny")

	rs.Assert(nil, scanEv1)

}

func TestPkgFlowTimeout(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, _ := model.NewTupleWithKeyValues("packageevent", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	rs.Assert(nil, pkgEvt)
	//time.Sleep(time.Second*20)
	scanEv, _ := model.NewTupleWithKeyValues("scanevent", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	err = scanEv.SetValue(nil, "eta", 3)
	if err != nil {
		fmt.Printf("[%s]\n", err)
	}

	rs.Assert(nil, scanEv)

	time.Sleep(time.Second * time.Duration(20))

}

func TestPkgFlowNormalWithMapValues(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, err := model.NewTupleWithKeyValues("packageevent", "1")
	if err != nil {
		fmt.Printf("error: [%s]\n", err)
		return
	}
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")
	fmt.Printf("Asserting package with key [%s]\n", pkgEvt.GetKey().String())

	rs.Assert(nil, pkgEvt)

	values := make(map[string]interface{})
	values["packageid"] = "1"
	values["curr"] = "sfo"
	values["next"] = "ny"
	values["eta"] = 5

	scanEv, _ := model.NewTuple("scanevent", values)

	rs.Assert(nil, scanEv)

	values = make(map[string]interface{})
	values["packageid"] = "1"
	values["curr"] = "ny"
	values["next"] = "done"
	values["eta"] = 5

	scanEv2, _ := model.NewTuple("scanevent", values)
	rs.Assert(nil, scanEv2)

	time.Sleep(time.Second * time.Duration(20))
}

func TestDuplicateAssert(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, _ := model.NewTupleWithKeyValues("package", "1")
	pkgEvt.SetString(nil, "curr", "none")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")

	err = rs.Assert(nil, pkgEvt)
	if err != nil {
		fmt.Printf("Error...[%s]\n", err)
		return
	}

	pkgEvt2, _ := model.NewTupleWithKeyValues("package", "1")
	pkgEvt2.SetString(nil, "curr", "sfo")
	pkgEvt2.SetString(nil, "next", "ny")
	pkgEvt2.SetString(nil, "status", "normal")

	err = rs.Assert(nil, pkgEvt2)
	if err != nil {
		fmt.Printf("Error...[%s]\n", err)
		return
	}
}

func TestSameTupleInstanceAssert(t *testing.T) {

	rs, err := createRuleSessionAndRules()
	if err != nil {
		fmt.Printf("Error [%s]\n", err)
		return
	}

	loadPkgRulesWithDeps(rs)
	rs.Start(nil)

	pkgEvt, _ := model.NewTupleWithKeyValues("package", "1")
	pkgEvt.SetString(nil, "curr", "none")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")

	err = rs.Assert(nil, pkgEvt)
	if err != nil {
		fmt.Printf("Error...[%s]\n", err)
		return
	}

	err = rs.Assert(nil, pkgEvt)
	if err != nil {
		fmt.Printf("Error...[%s]\n", err)
		return
	}
}



func createRuleSessionAndRules() (model.RuleSession, error) {
	rs, _ := ruleapi.GetOrCreateRuleSession("asession")

	tupleDescFileAbsPath := common.GetAbsPathForResource("src/github.com/project-flogo/rules/examples/trackntrace/trackntrace.json")

	dat, err := ioutil.ReadFile(tupleDescFileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	err = model.RegisterTupleDescriptors(string(dat))
	if err != nil {
		return nil, err
	}
	return rs, nil
}

func loadPkgRulesWithDeps(rs model.RuleSession) {

	//handle a package event, create a package in the packageAction
	rule := ruleapi.NewRule("packageevent")
	rule.AddCondition("truecondition", []string{"packageevent.none"}, truecondition, nil)
	rule.SetAction(packageeventAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	//handle a package, print package details in the packageAction
	rule1 := ruleapi.NewRule("package")
	rule1.AddCondition("packageCondition", []string{"package.none"}, packageCondition, nil)
	rule1.SetAction(packageAction)
	rule1.SetPriority(2)
	rs.AddRule(rule1)
	fmt.Printf("Rule added: [%s]\n", rule1.GetName())

	//handle a scan event, see if there is matching package if so, do necessary things such as set off a timer
	//for the next destination, etc in the scaneventAction
	rule2 := ruleapi.NewRule("scanevent")
	rule2.AddCondition("scaneventCondition", []string{"package.packageid", "scanevent.packageid", "package.curr", "package.next"}, scaneventCondition, nil)
	rule2.SetAction(scaneventAction)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())

	//handle a timeout event, triggered by scaneventAction, mark the package as delayed in scantimeoutAction
	rule3 := ruleapi.NewRule("scantimeout")
	rule3.AddCondition("scantimeoutCondition", []string{"package.packageid", "scantimeout.packageid"}, scantimeoutCondition, nil)
	rule3.SetAction(scantimeoutAction)
	rule3.SetPriority(1)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	//notify when a package is marked as delayed, print as such in the packagedelayedAction
	rule4 := ruleapi.NewRule("packagedelayed")
	rule4.AddCondition("packageDelayedCheck", []string{"package.status"}, packageDelayedCheck, nil)
	rule4.SetAction(packagedelayedAction)
	rule4.SetPriority(1)
	rs.AddRule(rule4)
	fmt.Printf("Rule added: [%s]\n", rule4.GetName())
}

//conditions and actions
func truecondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}

func packageeventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {

	pkgEvent := tuples["packageevent"]
	pkgid, _ := pkgEvent.GetString("packageid")
	fmt.Printf("Received a new package asserting package id[%s]\n", pkgid)

	//assert a package
	pkg, _ := model.NewTupleWithKeyValues(model.TupleType("package"), pkgid)
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

	scantmout, _ := model.NewTupleWithKeyValues(model.TupleType("scantimeout"), pkgid)
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
