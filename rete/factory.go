package rete

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/memimpl"
	"github.com/project-flogo/rules/rete/internal/types"
)

type TypeFactory struct {
	nw     *reteNetworkImpl
	config map[string]string
}

func NewFactory(nw *reteNetworkImpl, config map[string]string) *TypeFactory {
	tf := TypeFactory{nw, config}
	return &tf
}

func (f *TypeFactory) getJoinTable(rule model.Rule, identifiers []model.TupleType) types.JoinTable {
	if f.config == nil {
		jt := memimpl.NewJoinTable(f.nw, rule, identifiers)
		return jt
	} else {
		jtStore := f.config["jtstore"]
		if jtStore == "" || jtStore == "memory" {
			jt := memimpl.NewJoinTable(f.nw, rule, identifiers)
			return jt
		}
	}
	return nil
}

func (f *TypeFactory) getJoinTableRefs() types.JoinTableRefsInHdl {
	if f.config == nil {
		jtRef := memimpl.NewJoinTableRefsInHdlImpl()
		return jtRef
	} else {
		jtType := f.config["jtstore"]
		if jtType == "" || jtType == "memory" {
			jtRef := memimpl.NewJoinTableRefsInHdlImpl()
			return jtRef
		}
	}
	return nil
}

func (f *TypeFactory) getJoinTableCollection() types.JoinTableCollection {
	if f.config == nil {
		jtRef := memimpl.NewJoinTableCollection()
		return jtRef
	} else {
		jtType := f.config["jtstore"]
		if jtType == "" || jtType == "memory" {
			jtRef := memimpl.NewJoinTableCollection()
			return jtRef
		}
	}
	return nil
}

func (f *TypeFactory) getHandleCollection() types.HandleCollection {
	if f.config == nil {
		hc := memimpl.NewHandleCollection()
		return hc
	} else {
		jtType := f.config["jtstore"]
		if jtType == "" || jtType == "memory" {
			hc := memimpl.NewHandleCollection()
			return hc
		}
	}
	return nil
}
