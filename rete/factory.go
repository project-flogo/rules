package rete

import (
	"encoding/json"
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

func (f *TypeFactory) getJoinTableRefs() types.JtRefsService {
	var jtRefs types.JtRefsService
	if f.parsedJson == nil {
		jtRefs = mem.NewJoinTableRefsInHdlImpl(f.nw, f.parsedJson)

	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {
			idgen := rete["jt"].(string)
			if idgen == "" || idgen == "mem" {
				jtRefs = mem.NewJoinTableRefsInHdlImpl(f.nw, f.parsedJson)
			} else if idgen == "redis" {
				jtRefs = redis.NewJoinTableRefsInHdlImpl(f.nw, f.parsedJson)
			}
		}
	}
	return jtRefs
}

func (f *TypeFactory) getJoinTableCollection() types.JtService {
	var allJt types.JtService
	if f.parsedJson == nil {
		allJt = mem.NewJoinTableCollection(f.nw, f.parsedJson)

	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {
			idgen := rete["jt"].(string)
			if idgen == "" || idgen == "mem" {
				allJt = mem.NewJoinTableCollection(f.nw, f.parsedJson)
			} else if idgen == "redis" {
				allJt = redis.NewJoinTableCollection(f.nw, f.parsedJson)
			}
		}
	}
	return allJt
}

func (f *TypeFactory) getHandleCollection() types.HandleService {
	var hc types.HandleService
	if f.parsedJson == nil {
		hc = mem.NewHandleCollection(f.nw, f.parsedJson)
	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {
			idgen := rete["jt"].(string)
			if idgen == "" || idgen == "mem" {
				hc = mem.NewHandleCollection(f.nw, f.parsedJson)
			} else if idgen == "redis" {
				hc = redis.NewHandleCollection(f.nw, f.parsedJson)
			}
		}
	}
	return hc
}

func (f *TypeFactory) getIdGen() types.IdGen {
	var idg types.IdGen
	if f.parsedJson == nil {
		idg = mem.NewIdGenImpl(f.nw, f.parsedJson)
		return idg
	} else {
		rete := f.parsedJson["rete"].(map[string]interface{})
		if rete != nil {

			idgen := rete["idgen"].(string)
			if idgen == "" || idgen == "mem" {
				idg = mem.NewIdGenImpl(f.nw, f.parsedJson)
			} else if idgen == "redis" {
				idg = redis.NewIdGenImpl(f.nw, f.parsedJson)
			}
		}
	}
	return idg
}
