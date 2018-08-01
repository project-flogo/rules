package ruleaction

import (
	"testing"
	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/ruleapi"
	"fmt"
	"golang.org/x/net/context"
	"strconv"
	"time"
)

func TestPkgFlowTimeout (t *testing.T) {

	rs := createRuleSessionAndRules()

	loadPkgRules(rs)

	pkgEvt := model.NewStreamTuple(model.TupleTypeAlias("packageevent"))
	//ctx := context.TODO()
	pkgEvt.SetString(nil, rs, "packageid", "1")
	pkgEvt.SetString(nil, rs, "next", "sfo")
	pkgEvt.SetString(nil, rs, "status", "normal")
	pkgEvt.SetString(nil, rs, "isnew", "true")

	rs.Assert(nil, pkgEvt)
	scanEv := model.NewStreamTuple(model.TupleTypeAlias("scanevent"))
	scanEv.SetString(nil, rs, "packageid", "1")
	scanEv.SetString(nil, rs, "curr", "sfo")
	scanEv.SetString(nil, rs, "next", "ny")
	scanEv.SetString(nil, rs, "eta", "5")

	rs.Assert(nil, scanEv)

	time.Sleep(time.Second*10)

}

func TestPkgFlowNormal (t *testing.T) {

	rs := createRuleSessionAndRules()

	loadPkgRules(rs)

	pkgEvt := model.NewStreamTuple(model.TupleTypeAlias("packageevent"))
	//ctx := context.TODO()
	pkgEvt.SetString(nil, rs, "packageid", "1")
	pkgEvt.SetString(nil, rs, "next", "sfo")
	pkgEvt.SetString(nil, rs, "status", "normal")
	pkgEvt.SetString(nil, rs, "isnew", "true")

	rs.Assert(nil, pkgEvt)
	scanEv := model.NewStreamTuple(model.TupleTypeAlias("scanevent"))
	scanEv.SetString(nil, rs, "packageid", "1")
	scanEv.SetString(nil, rs, "curr", "sfo")
	scanEv.SetString(nil, rs, "next", "ny")
	scanEv.SetString(nil, rs, "eta", "5")

	rs.Assert(nil, scanEv)

	scanEv1 := model.NewStreamTuple(model.TupleTypeAlias("scanevent"))
	scanEv1.SetString(nil, rs, "packageid", "1")
	scanEv1.SetString(nil, rs, "curr", "ny")
	scanEv1.SetString(nil, rs, "next", "done")
	scanEv1.SetString(nil, rs, "eta", "5")
	rs.Assert(nil, scanEv1)

	time.Sleep(time.Second*10)

}

func loadPkgRules(rs model.RuleSession) {

	rule := ruleapi.NewRule("packageevent")
	rule.AddCondition("packageevent", []model.TupleTypeAlias{"packageevent"}, truecondition) // check for name "Bob" in n1
	rule.SetAction(packageeventAction)
	rule.SetPriority(1)
	rs.AddRule(rule)
	fmt.Printf("Rule added: [%s]\n", rule.GetName())

	rule1:= ruleapi.NewRule("package")
	rule1.AddCondition("packageevent1", []model.TupleTypeAlias{"package"}, packageCondition) // check for name "Bob" in n1
	rule1.SetAction(packageAction)
	rule1.SetPriority(2)
	rs.AddRule(rule1)
	fmt.Printf("Rule added: [%s]\n", rule1.GetName())
	//
	rule2 := ruleapi.NewRule("scanevent")
	rule2.AddCondition("scanevent", []model.TupleTypeAlias{"package", "scanevent"}, scaneventCondition) // check for name "Bob" in n1
	rule2.SetAction(scaneventAction)
	rule2.SetPriority(2)
	rs.AddRule(rule2)
	fmt.Printf("Rule added: [%s]\n", rule2.GetName())
	//
	rule3 := ruleapi.NewRule("scantimeout")
	rule3.AddCondition("packageevent", []model.TupleTypeAlias{"package", "scantimeout"}, scantimeoutCondition) // check for name "Bob" in n1
	rule3.SetAction(scantimeoutAction)
	rule3.SetPriority(1)
	rs.AddRule(rule3)
	fmt.Printf("Rule added: [%s]\n", rule3.GetName())

	rule4 := ruleapi.NewRule("packagedelayed")
	rule4.AddCondition("packageevent", []model.TupleTypeAlias{"package"}, packageDelayedCheck) // check for name "Bob" in n1
	rule4.SetAction(packagedelayedAction)
	rule4.SetPriority(1)
	rs.AddRule(rule4)
	fmt.Printf("Rule added: [%s]\n", rule4.GetName())
}

func packageeventAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {

	pkgEvent := tuples["packageevent"]
	pkgid := pkgEvent.GetString("packageid")
	fmt.Printf ("Received a new package asserting package id[%s]\n", pkgid)

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
	fmt.Printf ("Received a new scan event for package id[%s], current loc [%s], next loc [%s]\n", pkgid, scurr, snext)


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
	pkg.SetString (ctx, rs, "curr", scurr)
	pkg.SetString (ctx,  rs, "next", snext)

}

func scantimeoutCondition(ruleName string, condName string, tuples map[model.TupleTypeAlias]model.StreamTuple) bool {
	scantimeout := tuples["scantimeout"]
	pkg := tuples["package"]

	if scantimeout == nil || pkg == nil {
		fmt.Println("Should not get a nil tuple here! This is an error")
		return false
	}
	return scantimeout.GetString("packageid") == pkg.GetString("packageid") &&
		scantimeout.GetString("next") == pkg.GetString("next") && pkg.GetString("status") == "normal"
}

func scantimeoutAction(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleTypeAlias]model.StreamTuple) {

	pkg := tuples["package"].(model.MutableStreamTuple)

	pkgid := pkg.GetString("packageid")
	pcurr :=  pkg.GetString("curr")
	pnext :=  pkg.GetString("next")

	fmt.Printf ("Package id[%s] : Scan for dest [%s] did not arrive by ETA. Package currently at [%s]\n",
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
	fmt.Printf ("Received a new package id[%s], current loc [%s], next loc [%s]\n", pkgid, pcurr, pnext)
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

	fmt.Printf ("Package is now delayed id[%s]\n", pkgid)
}