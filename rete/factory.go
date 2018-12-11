package rete

import (
	"github.com/project-flogo/rules/common/model"
)

type StoreProvider int

const (
	Memory = iota
	Redis
)

type TypeFactory struct {
	config map[string]string
}

func NewFactory(config map[string]string) TypeFactory {
	tf := TypeFactory{config}
	return tf
}

func (f TypeFactory) getJoinTable(nw Network, rule model.Rule, identifiers []model.TupleType) joinTable {
	jtStore := f.config["jtstore"]
	if jtStore == "" || jtStore == "memory" {
		jt := newJoinTable(nw, rule, identifiers)
		return jt
	}
	return nil
}

func (f TypeFactory) getJoinTableRefs() joinTableRefsInHdl {
	jtType := f.config["jtstore"]
	if jtType == "" || jtType == "memory" {
		jtRef := newJoinTableRefsInHdlImpl()
		return jtRef
	}
	return nil
}
