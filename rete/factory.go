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
	jtStore := f.config["jtstore"]
	if jtStore == "" || jtStore == "memory" {
		jt := memimpl.NewJoinTable(f.nw, rule, identifiers)
		return jt
	}
	return nil
}

func (f *TypeFactory) getJoinTableRefs() types.JoinTableRefsInHdl {
	jtType := f.config["jtstore"]
	if jtType == "" || jtType == "memory" {
		jtRef := memimpl.NewJoinTableRefsInHdlImpl()
		return jtRef
	}
	return nil
}
