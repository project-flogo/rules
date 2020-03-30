package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
	"github.com/project-flogo/rules/ruleapi"
)

func main() {
	err := example(false)
	if err != nil {
		panic(err)
	}
}

var (
	currentEventType string
)

func example(redis bool) error {
	//Load the tuple descriptor file (relative to GOPATH)
	tupleDescAbsFileNm := common.GetPathForResource("examples/statemachine/rulesapp.json", "./rulesapp.json")
	tupleDescriptor := common.FileToString(tupleDescAbsFileNm)
	currentEventType = "none"
	//First register the tuple descriptors
	err := model.RegisterTupleDescriptors(tupleDescriptor)
	if err != nil {
		return err
	}

	//Create a RuleSession
	store := ""
	if redis {
		store = "rsconfig.json"
	}
	rs, err := ruleapi.GetOrCreateRuleSession("asession", store)
	if err != nil {
		return err
	}

	events := make(map[string]int, 8)

	//// check if the packaage is in sitting state
	rule := ruleapi.NewRule("cPackageInSittingRule")
	err = rule.AddCondition("c1", []string{"package.state"}, cPackageInSitting, events)
	if err != nil {
		return err
	}
	serviceCfg := &config.ServiceDescriptor{
		Name:     "aPackageInSitting",
		Function: aPackageInSitting,
		Type:     "function",
	}
	aService, err := ruleapi.NewActionService(serviceCfg)
	if err != nil {
		return err
	}
	rule.SetActionService(aService)
	rule.SetContext(events)
	rule.SetPriority(1)
	err = rs.AddRule(rule)
	if err != nil {
		return err
	}

	//// check if the packaage is in Delayed state
	rule2 := ruleapi.NewRule("packageInDelayedRule")
	err = rule2.AddCondition("c1", []string{"package.state"}, cPackageInDelayed, events)
	if err != nil {
		return err
	}
	serviceCfg2 := &config.ServiceDescriptor{
		Name:     "aPackageInDelayed",
		Function: aPackageInDelayed,
		Type:     "function",
	}
	aService2, err := ruleapi.NewActionService(serviceCfg2)
	if err != nil {
		return err
	}
	rule2.SetActionService(aService2)
	rule2.SetContext(events)
	rule2.SetPriority(1)
	err = rs.AddRule(rule2)
	if err != nil {
		return err
	}

	//// check if the packaage is in moving state
	rule3 := ruleapi.NewRule("packageInMovingRule")
	err = rule3.AddCondition("c1", []string{"package.state"}, cPackageInMoving, events)
	if err != nil {
		return err
	}
	serviceCfg3 := &config.ServiceDescriptor{
		Name:     "aPackageInMoving",
		Function: aPackageInMoving,
		Type:     "function",
	}
	aService3, err := ruleapi.NewActionService(serviceCfg3)
	if err != nil {
		return err
	}
	rule3.SetActionService(aService3)
	rule3.SetContext(events)
	rule3.SetPriority(1)
	err = rs.AddRule(rule3)
	if err != nil {
		return err
	}

	//// check if the packaage is in Dropped state
	rule4 := ruleapi.NewRule("packageInDroppedRule")
	err = rule4.AddCondition("c1", []string{"package.state"}, cPackageInDropped, events)
	if err != nil {
		return err
	}
	serviceCfg4 := &config.ServiceDescriptor{
		Name:     "aPackageInDropped",
		Function: aPackageInDropped,
		Type:     "function",
	}
	aService4, err := ruleapi.NewActionService(serviceCfg4)
	if err != nil {
		return err
	}
	rule4.SetActionService(aService4)
	rule4.SetContext(events)
	rule4.SetPriority(1)
	err = rs.AddRule(rule4)
	if err != nil {
		return err
	}

	//// check if the packaage is in normal state and print
	rule5 := ruleapi.NewRule("printPackageRule")
	err = rule5.AddCondition("c1", []string{"package"}, cPackageEvent, events)
	if err != nil {
		return err
	}
	serviceCfg5 := &config.ServiceDescriptor{
		Name:     "aPrintPackage",
		Function: aPrintPackage,
		Type:     "function",
	}
	aService5, err := ruleapi.NewActionService(serviceCfg5)
	if err != nil {
		return err
	}
	rule5.SetActionService(aService5)
	rule5.SetContext(events)
	rule5.SetPriority(2)
	err = rs.AddRule(rule5)
	if err != nil {
		return err
	}

	// check for moveevent
	rule6 := ruleapi.NewRule("printMoveEventRule")
	err = rule6.AddCondition("c1", []string{"moveevent"}, cMoveEvent, events)
	if err != nil {
		return err
	}
	serviceCfg6 := &config.ServiceDescriptor{
		Name:     "aPrintMoveEvent",
		Function: aPrintMoveEvent,
		Type:     "function",
	}
	aService6, err := ruleapi.NewActionService(serviceCfg6)
	if err != nil {
		return err
	}
	rule6.SetActionService(aService6)
	rule6.SetContext(events)
	rule6.SetPriority(3)
	err = rs.AddRule(rule6)
	if err != nil {
		return err
	}

	// check if the package exists for received moveevent and join it with the package
	rule7 := ruleapi.NewRule("joinMoveEventAndPackageEventRule")
	err = rule7.AddCondition("c1", []string{"moveevent", "package"}, cMoveEventPkg, events)
	if err != nil {
		return err
	}
	serviceCfg7 := &config.ServiceDescriptor{
		Name:     "aJoinMoveEventAndPackage",
		Function: aJoinMoveEventAndPackage,
		Type:     "function",
	}
	aService7, err := ruleapi.NewActionService(serviceCfg7)
	if err != nil {
		return err
	}
	rule7.SetActionService(aService7)
	rule7.SetContext(events)
	rule7.SetPriority(4)
	err = rs.AddRule(rule7)
	if err != nil {
		return err
	}

	// check for movetimeoutevent for package
	rule8 := ruleapi.NewRule("aMoveTimeoutEventRule")
	err = rule8.AddCondition("c1", []string{"movetimeoutevent"}, cMoveTimeoutEvent, events)
	if err != nil {
		return err
	}
	serviceCfg8 := &config.ServiceDescriptor{
		Name:     "aMoveTimeoutEvent",
		Function: aMoveTimeoutEvent,
		Type:     "function",
	}
	aService8, err := ruleapi.NewActionService(serviceCfg8)
	if err != nil {
		return err
	}
	rule8.SetActionService(aService8)
	rule8.SetContext(events)
	rule8.SetPriority(5)
	err = rs.AddRule(rule8)
	if err != nil {
		return err
	}

	//Join movetimeoutevent and package
	rule9 := ruleapi.NewRule("joinMoveTimeoutEventAndPackage")
	err = rule9.AddCondition("c1", []string{"movetimeoutevent", "package"}, cMoveTimeoutEventPkg, events)
	if err != nil {
		return err
	}
	serviceCfg9 := &config.ServiceDescriptor{
		Name:     "aJoinMoveTimeoutEventAndPackage",
		Function: aJoinMoveTimeoutEventAndPackage,
		Type:     "function",
	}
	aService9, err := ruleapi.NewActionService(serviceCfg9)
	if err != nil {
		return err
	}
	rule9.SetActionService(aService9)
	rule9.SetContext(events)
	rule9.SetPriority(6)
	err = rs.AddRule(rule9)
	if err != nil {
		return err
	}

	rs.SetStartupFunction(AssertThisPackage)

	//set a transaction handler
	rs.RegisterRtcTransactionHandler(txHandler, nil)
	//Start the rule session
	err = rs.Start(nil)
	if err != nil {
		return err
	}

	t6, err := model.NewTupleWithKeyValues("moveevent", "PACKAGE1")
	if err != nil {
		return err
	}
	t6.SetString(nil, "packageid", "PACKAGE1")
	t6.SetString(nil, "targetstate", "sitting")
	err = rs.Assert(nil, t6)
	if err != nil {
		return err
	}

	finished := make(chan bool)
	go sleepfunc(finished)

	t1, err := model.NewTupleWithKeyValues("moveevent", "PACKAGE2")
	if err != nil {
		return err
	}
	t1.SetString(nil, "packageid", "PACKAGE2")
	t1.SetString(nil, "targetstate", "normal")
	err = rs.Assert(nil, t1)
	if err != nil {
		return err
	}
	fmt.Println("package 2", t1)

	t2, err := model.NewTupleWithKeyValues("moveevent", "PACKAGE2")
	if err != nil {
		return err
	}
	t2.SetString(nil, "packageid", "PACKAGE2")
	t2.SetString(nil, "targetstate", "sitting")
	err = rs.Assert(nil, t2)
	if err != nil {
		return err
	}

	t3, err := model.NewTupleWithKeyValues("moveevent", "PACKAGE2")
	if err != nil {
		return err
	}
	t3.SetString(nil, "packageid", "PACKAGE2")
	t3.SetString(nil, "targetstate", "moving")
	err = rs.Assert(nil, t3)
	if err != nil {
		return err
	}

	t4, err := model.NewTupleWithKeyValues("moveevent", "PACKAGE2")
	if err != nil {
		return err
	}
	t4.SetString(nil, "packageid", "PACKAGE2")
	t4.SetString(nil, "targetstate", "dropped")
	err = rs.Assert(nil, t4)
	if err != nil {
		return err
	}

	<-finished

	//delete the rule
	rs.DeleteRule(rule2.GetName())

	//unregister the session
	rs.Unregister()

	return nil
}

func AssertThisPackage(ctx context.Context, rs model.RuleSession, startupCtx map[string]interface{}) (err error) {
	fmt.Printf("In startup rule function..\n")
	pkg, _ := model.NewTupleWithKeyValues("package", "PACKAGE1")
	pkg.SetString(nil, "state", "normal")
	rs.Assert(nil, pkg)
	return nil
}

func cPackageEvent(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	if pkg != nil {
		state, _ := pkg.GetString("state")
		return state == "normal"
	}
	return false
}

func aPrintPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"]
	pkgid, _ := pkg.GetString("id")
	fmt.Printf("Received package [%s]\n", pkgid)
}

func aPrintMoveEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["moveevent"]
	meid, _ := me.GetString("id")
	s, _ := me.GetString("targetstate")
	pkgID, _ := me.GetString("packageid")

	fmt.Printf("Received a 'moveevent' [%s] target state [%s]\n", meid, s)

	if s == "normal" {
		pkg, _ := model.NewTupleWithKeyValues("package", pkgID)
		pkg.SetString(nil, "state", "normal")
		err := rs.Assert(ctx, pkg)
		if err != nil {
			fmt.Println("Tuple already inserted: ", pkgID)
		} else {
			fmt.Println("Tuple inserted successfully: ", pkgID)
		}
	}

}

func aJoinMoveEventAndPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["moveevent"]
	mepkgid, _ := me.GetString("packageid")
	s, _ := me.GetString("targetstate")

	pkg := tuples["package"]
	pkgid, _ := pkg.GetString("id")
	pkgState, _ := pkg.GetString("state")

	if strings.Compare("sitting", s) == 0 {
		currentEventType = "sitting"
	} else {
		currentEventType = ""
	}

	if currentEventType == "sitting" {

		if pkgState == "normal" {
			fmt.Printf("Joining a 'moveevent' with packageid [%s] to package [%s], target state [%s]\n", mepkgid, pkgid, s)

			//change the package's state to "sitting"
			pkgMutable := pkg.(model.MutableTuple)
			pkgMutable.SetString(ctx, "state", "sitting")

			//very first sitting event since the last notsitting event.
			id, _ := common.GetUniqueId()
			timeoutEvent, _ := model.NewTupleWithKeyValues("movetimeoutevent", id)
			timeoutEvent.SetString(ctx, "packageid", pkgid)
			timeoutEvent.SetInt(ctx, "timeoutinmillis", 10000)
			fmt.Printf("Starting a 10s timer.. [%s]\n", pkgid)
			rs.ScheduleAssert(ctx, 10000, pkgid, timeoutEvent)
		}
	} else {
		if strings.Compare("moving", s) == 0 && pkgState == "sitting" {

			fmt.Printf("Joining a 'moveevent' with packageid [%s] to package [%s], target state [%s]\n", mepkgid, pkgid, s)

			//a non-sitting event, cancel a previous timer
			rs.CancelScheduledAssert(ctx, pkgid)

			pkgMutable := pkg.(model.MutableTuple)
			pkgMutable.SetString(ctx, "state", "moving")

		} else if strings.Compare("dropped", s) == 0 && pkgState == "moving" {

			fmt.Printf("Joining a 'moveevent' with packageid [%s] to package [%s], target state [%s]\n", mepkgid, pkgid, s)

			pkgMutable := pkg.(model.MutableTuple)
			pkgMutable.SetString(ctx, "state", "dropped")

		}

	}

}

func aMoveTimeoutEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples["movetimeoutevent"]
	id, _ := t1.GetString("id")
	pkgid, _ := t1.GetString("packageid")
	tomillis, _ := t1.GetInt("timeoutinmillis")

	if t1 != nil {
		fmt.Printf("Received a 'movetimeoutevent' id [%s], packageid [%s], timeoutinmillis [%d]\n", id, pkgid, tomillis)
	}
}

func aJoinMoveTimeoutEventAndPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["movetimeoutevent"]
	//meid, _ := me.GetString("id")
	epkgid, _ := me.GetString("packageid")
	toms, _ := me.GetInt("timeoutinmillis")

	pkg := tuples["package"].(model.MutableTuple)
	pkgid, _ := pkg.GetString("id")

	fmt.Printf("Joining a 'movetimeoutevent' [%s] to package [%s], timeout [%d]\n", epkgid, pkgid, toms)

	//change the package's state to "delayed"
	pkg.SetString(ctx, "state", "delayed")
}

func cPackageInSitting(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	if pkg != nil {
		//pkgid, _ := pkg.GetString("id")
		state, _ := pkg.GetString("state")
		return state == "sitting"
	}
	return false
}

func aPackageInSitting(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"]
	if pkg != nil {
		pkgid, _ := pkg.GetString("id")
		fmt.Printf("PACKAGE [%s] is Sitting\n", pkgid)
	}
}

func cMoveEvent(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	mpkg := tuples["moveevent"]
	if mpkg != nil {
		mpkgid, _ := mpkg.GetString("packageid")
		return len(mpkgid) != 0
	}
	return false
}

func cMoveTimeoutEvent(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	mpkg := tuples["movetimeoutevent"]
	if mpkg != nil {
		mpkgid, _ := mpkg.GetString("packageid")
		return len(mpkgid) != 0
	}
	return false
}

func cMoveTimeoutEventPkg(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	mpkg := tuples["movetimeoutevent"]
	pkg := tuples["package"]
	if mpkg != nil {
		if pkg != nil {
			pkgid, _ := pkg.GetString("id")
			mpkgid, _ := mpkg.GetString("packageid")
			return pkgid == mpkgid
		}
	}
	return false
}

func cMoveEventPkg(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	mpkg := tuples["moveevent"]
	pkg := tuples["package"]
	if mpkg != nil {
		if pkg != nil {
			pkgid, _ := pkg.GetString("id")
			mpkgid, _ := mpkg.GetString("packageid")
			return pkgid == mpkgid
		}
	}
	return false
}

func cPackageInDelayed(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	if pkg != nil {
		//pkgid, _ := pkg.GetString("id")
		state, _ := pkg.GetString("state")
		return state == "delayed"
	}
	return false
}

func aPackageInDelayed(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"]
	if pkg != nil {
		pkgid, _ := pkg.GetString("id")
		fmt.Printf("PACKAGE [%s] is Delayed\n", pkgid)
		rs.Retract(ctx, pkg)
	}
}

func cPackageInMoving(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	if pkg != nil {
		//pkgid, _ := pkg.GetString("id")
		state, _ := pkg.GetString("state")
		return state == "moving"
	}
	return false
}

func aPackageInMoving(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"]
	if pkg != nil {
		pkgid, _ := pkg.GetString("id")
		fmt.Printf("PACKAGE [%s] is Moving\n", pkgid)
	}
}

func cPackageInDropped(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	if pkg != nil {
		//pkgid, _ := pkg.GetString("id")
		state, _ := pkg.GetString("state")
		return state == "dropped"
	}
	return false
}

func aPackageInDropped(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"]
	if pkg != nil {
		pkgid, _ := pkg.GetString("id")
		fmt.Printf("PACKAGE [%s] is Dropped\n", pkgid)
		rs.Retract(ctx, pkg)
	}
}

func getFileContent(filePath string) string {
	absPath := common.GetAbsPathForResource(filePath)
	return common.FileToString(absPath)
}

func txHandler(ctx context.Context, rs model.RuleSession, rtxn model.RtcTxn, handlerCtx interface{}) {

	store := rs.GetStore()
	store.SaveTuples(rtxn.GetRtcAdded())

	store.SaveModifiedTuples(rtxn.GetRtcModified())

	store.DeleteTuples(rtxn.GetRtcDeleted())

}
func sleepfunc(finished chan bool) {
	time.Sleep(15 * time.Second)
	finished <- true
}
