package ruleaction

import (
	"github.com/TIBCOSoftware/bego/common/model"
	"testing"
	"time"
)

func TestPkgFlowNormalWithDeps(t *testing.T) {

	rs := createRuleSessionAndRules()

	loadPkgRulesWithDeps(rs)

	pkgEvt := model.NewTuple(model.TupleType("packageevent"))
	//ctx := context.TODO()
	pkgEvt.SetString(nil, "packageid", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	rs.Assert(nil, pkgEvt)
	//time.Sleep(time.Second*20)
	scanEv := model.NewTuple(model.TupleType("scanevent"))
	scanEv.SetString(nil, "packageid", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	scanEv.SetString(nil, "eta", "5")

	rs.Assert(nil, scanEv)

	scanEv1 := model.NewTuple(model.TupleType("scanevent"))
	scanEv1.SetString(nil, "packageid", "1")
	scanEv1.SetString(nil, "curr", "ny")
	scanEv1.SetString(nil, "next", "done")
	scanEv1.SetString(nil, "eta", "5")
	rs.Assert(nil, scanEv1)

}

func TestPkgFlowNormalWithDepsWithTimeout(t *testing.T) {

	rs := createRuleSessionAndRules()

	loadPkgRulesWithDeps(rs)

	pkgEvt := model.NewTuple(model.TupleType("packageevent"))
	//ctx := context.TODO()
	pkgEvt.SetString(nil, "packageid", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	rs.Assert(nil, pkgEvt)
	//time.Sleep(time.Second*20)
	scanEv := model.NewTuple(model.TupleType("scanevent"))
	scanEv.SetString(nil, "packageid", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	scanEv.SetString(nil, "eta", "5")

	rs.Assert(nil, scanEv)

	time.Sleep(time.Second * time.Duration(20))

}
