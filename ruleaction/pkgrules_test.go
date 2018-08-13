package ruleaction

import (
	"github.com/TIBCOSoftware/bego/common/model"
	"testing"
	"time"
	"fmt"
)

func TestPkgFlowNormalWithDeps(t *testing.T) {

	rs := createRuleSessionAndRules()

	loadPkgRulesWithDeps(rs)

	pkgEvt, _ := model.NewTuple(model.TupleType("packageevent"))
	//ctx := context.TODO()
	pkgEvt.SetString(nil, "packageid", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	rs.Assert(nil, pkgEvt)
	//time.Sleep(time.Second*20)
	scanEv, _ := model.NewTuple(model.TupleType("scanevent"))
	scanEv.SetString(nil, "packageid", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	err := scanEv.SetValue(nil, "eta", "10")
	if err != nil {
		fmt.Printf("[%s]\n", err)
	}

	rs.Assert(nil, scanEv)

	scanEv1, _ := model.NewTuple(model.TupleType("scanevent"))
	scanEv1.SetString(nil, "packageid", "1")
	scanEv1.SetString(nil, "curr", "ny")
	scanEv1.SetString(nil, "next", "done")
	scanEv.SetString(nil, "next", "ny")

	rs.Assert(nil, scanEv1)

}

func TestPkgFlowNormalWithDepsWithTimeout(t *testing.T) {

	rs := createRuleSessionAndRules()

	loadPkgRulesWithDeps(rs)

	pkgEvt, _ := model.NewTuple(model.TupleType("packageevent"))
	//ctx := context.TODO()
	pkgEvt.SetString(nil, "packageid", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	rs.Assert(nil, pkgEvt)
	//time.Sleep(time.Second*20)
	scanEv, _ := model.NewTuple(model.TupleType("scanevent"))
	scanEv.SetString(nil, "packageid", "1")
	scanEv.SetString(nil, "curr", "sfo")
	scanEv.SetString(nil, "next", "ny")
	err := scanEv.SetValue(nil, "eta", "3")
	if err != nil {
		fmt.Printf("[%s]\n", err)
	}

	rs.Assert(nil, scanEv)

	time.Sleep(time.Second * time.Duration(20))

}

func TestPkgFlowNormalWithDepsMapValues(t *testing.T) {

	rs := createRuleSessionAndRules()

	loadPkgRulesWithDeps(rs)

	pkgEvt, _:= model.NewTuple(model.TupleType("packageevent"))
	pkgEvt.SetString(nil, "packageid", "1")
	pkgEvt.SetString(nil, "next", "sfo")
	pkgEvt.SetString(nil, "status", "normal")
	pkgEvt.SetString(nil, "isnew", "true")

	rs.Assert(nil, pkgEvt)

	values := make (map[string]string)
	values["packageid"] = "1"
	values["curr"] = "sfo"
	values["next"] = "ny"
	values["eta"] = "5"

	scanEv, _ := model.NewTupleFromStringMap("scanevent", values)

	rs.Assert(nil, scanEv)

	time.Sleep(time.Second * time.Duration(20))
}
