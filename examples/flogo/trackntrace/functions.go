package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/project-flogo/rules/common"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/config"
)

var (
	lastEventType    string
	currentEventType string
)

func init() {

	config.RegisterStartupRSFunction("res://rulesession:simple", AssertThisPackage)

	//rule printPackage
	config.RegisterConditionEvaluator("cPackageEvent", cPackageEvent)
	config.RegisterActionFunction("aPrintPackage", aPrintPackage)

	//rule printMoveEvent
	config.RegisterConditionEvaluator("cMoveEvent", cMoveEvent)
	config.RegisterActionFunction("aPrintMoveEvent", aPrintMoveEvent)

	//rule joinMoveEventAndPackage
	config.RegisterConditionEvaluator("cJoinMoveEventAndPackage", cMoveEventPkg)
	config.RegisterActionFunction("aJoinMoveEventAndPackage", aJoinMoveEventAndPackage)

	//rule timeoutDelayed
	config.RegisterConditionEvaluator("cMoveTimeoutEvent", cMoveTimeoutEvent)
	config.RegisterActionFunction("aMoveTimeoutEvent", aMoveTimeoutEvent)

	//rule joinMoveTimeoutEventAndPackage
	config.RegisterConditionEvaluator("cJoinMoveTimeoutEventAndPackage", cMoveTimeoutEventPkg)
	config.RegisterActionFunction("aJoinMoveTimeoutEventAndPackage", aJoinMoveTimeoutEventAndPackage)

	//rule packageInSitting
	config.RegisterConditionEvaluator("cPackageInSitting", cPackageInSitting)
	config.RegisterActionFunction("aPackageInSitting", aPackageInSitting)

	//rule packageInDelayed
	config.RegisterConditionEvaluator("cPackageInDelayed", cPackageInDelayed)
	config.RegisterActionFunction("aPackageInDelayed", aPackageInDelayed)

	//rule packageInMoving
	config.RegisterConditionEvaluator("cPackageInMoving", cPackageInMoving)
	config.RegisterActionFunction("aPackageInMoving", aPackageInMoving)

	//rule packageInDropped
	config.RegisterConditionEvaluator("cPackageInDropped", cPackageInDropped)
	config.RegisterActionFunction("aPackageInDropped", aPackageInDropped)

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

func cPackageEvent(ruleName string, condName string, tuples map[model.TupleType]model.Tuple, ctx model.RuleContext) bool {
	pkg := tuples["package"]
	return pkg != nil
}

func aPrintPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	pkg := tuples["package"]
	pkgid, _ := pkg.GetString("id")
	fmt.Printf("Received package [%s]\n", pkgid)
}

func aPrintMoveEvent(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["moveevent"]
	meid, _ := me.GetString("id")
	s, _ := me.GetString("changeStateTo")

	fmt.Printf("Received a 'moveevent' [%s] change state to [%s]\n", meid, s)
}

func aJoinMoveEventAndPackage(ctx context.Context, rs model.RuleSession, ruleName string, tuples map[model.TupleType]model.Tuple, ruleCtx model.RuleContext) {
	me := tuples["moveevent"]
	mepkgid, _ := me.GetString("packageid")
	s, _ := me.GetString("changeStateTo")

	pkg := tuples["package"]
	pkgid, _ := pkg.GetString("id")

	fmt.Printf("Joining a 'moveevent' with packageid [%s] to package [%s], change state to [%s]\n", mepkgid, pkgid, s)

	if strings.Compare("sitting", s) == 0 {
		currentEventType = "sitting"
	}

	if currentEventType == "sitting" {
		if lastEventType != "sitting" {

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
		//a non-sitting event, cancel a previous timer
		rs.CancelScheduledAssert(ctx, pkgid)

		if strings.Compare("moving", s) == 0 {
			pkgMutable := pkg.(model.MutableTuple)
			pkgMutable.SetString(ctx, "state", "moving")
		} else if strings.Compare("dropped", s) == 0 {
			pkgMutable := pkg.(model.MutableTuple)
			pkgMutable.SetString(ctx, "state", "dropped")
		}

	}
	lastEventType = currentEventType
	currentEventType = "notsitting"
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
			return strings.Compare(pkgid, mpkgid) == 0
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
			return strings.Compare(pkgid, mpkgid) == 0
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
		rs.Delete(ctx, pkg)
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
		rs.Delete(ctx, pkg)
	}
}
