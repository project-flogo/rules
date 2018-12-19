package rete

import (
	"encoding/json"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/mem"
	"github.com/project-flogo/rules/rete/internal/redis"
	"github.com/project-flogo/rules/rete/internal/types"
)

type TypeFactory struct {
	nw         *reteNetworkImpl
	config     string
	parsedJson map[string]interface{}
}

func NewFactory(nw *reteNetworkImpl, config string) *TypeFactory {
	tf := TypeFactory{}
	tf.config = config
	json.Unmarshal([]byte(config), &tf.parsedJson)
	tf.nw = nw

	return &tf
}

func (f *TypeFactory) getJoinTable(rule model.Rule, identifiers []model.TupleType) types.JoinTable {
	var jt types.JoinTable
	if f.parsedJson == nil {
		jt = mem.NewJoinTable(f.nw, rule, identifiers)

	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {
			idgen := rete["jt"].(string)
			if idgen == "" || idgen == "mem" {
				jt = mem.NewJoinTable(f.nw, rule, identifiers)
			} else if idgen == "redis" {
				jt = mem.NewJoinTable(f.nw, rule, identifiers)
			}
		}
	}
	return jt
}

func (f *TypeFactory) getJoinTableRefs() types.JoinTableRefsInHdl {
	var jtRefs types.JoinTableRefsInHdl
	if f.parsedJson == nil {
		jtRefs = mem.NewJoinTableRefsInHdlImpl()

	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {
			idgen := rete["jt"].(string)
			if idgen == "" || idgen == "mem" {
				jtRefs = mem.NewJoinTableRefsInHdlImpl()
			} else if idgen == "redis" {
				jtRefs = mem.NewJoinTableRefsInHdlImpl()
			}
		}
	}
	return jtRefs
}

func (f *TypeFactory) getJoinTableCollection() types.JoinTableCollection {
	var allJt types.JoinTableCollection
	if f.parsedJson == nil {
		allJt = mem.NewJoinTableCollection()

	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {
			idgen := rete["jt"].(string)
			if idgen == "" || idgen == "mem" {
				allJt = mem.NewJoinTableCollection()
			} else if idgen == "redis" {
				allJt = mem.NewJoinTableCollection()
			}
		}
	}
	return allJt
}

func (f *TypeFactory) getHandleCollection() types.HandleCollection {
	var hc types.HandleCollection
	if f.parsedJson == nil {
		hc = mem.NewHandleCollection()
	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {
			idgen := rete["jt"].(string)
			if idgen == "" || idgen == "mem" {
				hc = mem.NewHandleCollection()
			} else if idgen == "redis" {
				hc = mem.NewHandleCollection()
			}
		}
	}
	return hc
}

func (f *TypeFactory) getIdGen() types.IdGen {
	var idg types.IdGen
	if f.parsedJson == nil {
		idg = mem.NewIdImpl(f.config)
		return idg
	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {

			idgen := rete["idgen"].(string)
			if idgen == "" || idgen == "mem" {
				idg = mem.NewIdImpl(f.config)
			} else if idgen == "redis" {
				idg = redis.NewIdImpl(f.config)
			}
		}
	}
	return idg
}
