package main

import (
	"context"
	"fmt"

	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
)

var (
	lastEventType    string
	currentEventType string
)

//add this sample file to your flogo project
func init() {

	config.RegisterStartupRSFunction("res://rulesession:simple", AssertThisPackage)

	//rule printPackage
	config.RegisterConditionEvaluator("cPackageEvent", cTruecondition)
	config.RegisterActionFunction("aPrintPackage", aPrintPackage)

	//rule printMoveEvent
	config.RegisterConditionEvaluator("cMoveEvent", cTruecondition)
	config.RegisterActionFunction("aPrintMoveEvent", aPrintMoveEvent)

	//rule joinMoveEventAndPackage
	config.RegisterConditionEvaluator("cJoinMoveEventAndPackage", cTruecondition)
	config.RegisterActionFunction("aJoinMoveEventAndPackage", aJoinMoveEventAndPackage)

	//rule timeoutSitting
	config.RegisterConditionEvaluator("cMoveTimeoutEvent", cTruecondition)
	config.RegisterActionFunction("aMoveTimeoutEvent", aMoveTimeoutEvent)

	//rule joinMoveTimeoutEventAndPackage
	config.RegisterConditionEvaluator("cJoinMoveTimeoutEventAndPackage", cTruecondition)
	config.RegisterActionFunction("aJoinMoveTimeoutEventAndPackage", aJoinMoveTimeoutEventAndPackage)

	//rule packageInSitting
	config.RegisterConditionEvaluator("cPackageInSitting", cPackageInSitting)
	config.RegisterActionFunction("aPackageInSitting", aPackageInSitting)

	lastEventType = "none"
	currentEventType = "none"
}

func AssertThisPackage(ctx context.Context, rs model.RuleSession, startupCtx map[string]interface{}) (err error) {
	fmt.Printf("In startup rule function..\n")
	pkg, _ := model.NewTupleWithKeyValues("package", "PACKAGE1")
	pkg.SetString(nil, "state", "normal")
	rs.Assert(nil, pkg)
	return nil
}

func cTruecondition(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	return true
}

func aPrintPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"]
	pkgid, _ := pkg.GetString("id")
	fmt.Printf("Received package [%s]\n", pkgid)
}

func aPrintMoveEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["moveevent"]
	meid, _ := me.GetString("id")
	s, _ := me.GetDouble("sitting")
	m, _ := me.GetDouble("moving")
	d, _ := me.GetDouble("dropped")

	fmt.Printf("Received a 'moveevent' [%s] sitting [%f], moving [%f], dropped [%f]\n",
		meid, s, m, d)
}

func aJoinMoveEventAndPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["moveevent"]
	mepkgid, _ := me.GetString("packageid")
	s, _ := me.GetDouble("sitting")
	m, _ := me.GetDouble("moving")
	d, _ := me.GetDouble("dropped")

	pkg := tuples["package"]
	pkgid, _ := pkg.GetString("id")

	fmt.Printf("Joining a 'moveevent' with packageid [%s] to package [%s], sitting [%f], moving [%f], dropped [%f]\n",
		mepkgid, pkgid, s, m, d)

	if s > 0.5 {
		currentEventType = "sitting"
	} else {
		currentEventType = "notsitting"
	}

	if currentEventType == "sitting" {
		if lastEventType != "sitting" {
			//very first sitting event since the last notsitting event.
			id, _ := common.GetUniqueId()
			timeoutEvent, _ := model.NewTupleWithKeyValues("movetimeoutevent", id)
			timeoutEvent.SetString(ctx, "packageid", pkgid)
			timeoutEvent.SetInt(ctx, "timeoutinmillis", 15000)
			fmt.Printf("Starting a 15s timer.. [%s]\n", pkgid)
			rs.ScheduleAssert(ctx, 15000, pkgid, timeoutEvent)
		}
	} else { //a non-sitting event, cancel a previous timer
		rs.CancelScheduledAssert(ctx, pkgid)
	}
	lastEventType = currentEventType
}

func aMoveTimeoutEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	t1 := tuples["movetimeoutevent"]
	id, _ := t1.GetString("id")
	pkgid, _ := t1.GetString("packageid")
	tomillis, _ := t1.GetInt("timeoutinmillis")

	if t1 != nil {
		fmt.Printf("'movetimeoutevent' id [%s], packageid [%s], timeoutinmillis [%d]\n",
			id, pkgid, tomillis)
	}
}

func aJoinMoveTimeoutEventAndPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["movetimeoutevent"]
	//meid, _ := me.GetString("id")
	epkgid, _ := me.GetString("packageid")
	toms, _ := me.GetInt("timeoutinmillis")

	pkg := tuples["package"].(model.MutableTuple)
	pkgid, _ := pkg.GetString("id")

	fmt.Printf("Joining a 'movetimeoutevent' [%s] to package [%s], timeout [%d]\n",
		epkgid, pkgid, toms)

	//change the package's state to "sitting"
	pkg.SetString(ctx, "state", "sitting")
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
		fmt.Printf("PACKAGE [%s] is STTTING\n", pkgid)
	}
}
