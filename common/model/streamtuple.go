package model

import (
	"context"
	"time"
)

type streamTupleImpl struct {
	dataSource TupleTypeAlias
	tuples     map[string]interface{}
}

func NewStreamTuple(dataSource TupleTypeAlias) MutableStreamTuple {
	st := streamTupleImpl{}
	st.initStreamTuple(dataSource)
	return &st
}

func (st *streamTupleImpl) initStreamTuple(dataSource TupleTypeAlias) {
	st.tuples = make(map[string]interface{})
	st.dataSource = dataSource
}

func (st *streamTupleImpl) GetTypeAlias() TupleTypeAlias {
	return st.dataSource
}

func (st *streamTupleImpl) GetString(name string) string {
	v := st.tuples[name]
	return v.(string)
}
func (st *streamTupleImpl) GetInt(name string) int {
	v := st.tuples[name]
	return v.(int)
}
func (st *streamTupleImpl) GetFloat(name string) float64 {
	v := st.tuples[name]
	return v.(float64)
}
func (st *streamTupleImpl) GetDateTime(name string) time.Time {
	v := st.tuples[name]
	return v.(time.Time)
}

func (st *streamTupleImpl) SetString(ctx context.Context, rs RuleSession, name string, value string) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.dataSource, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}

}
func (st *streamTupleImpl) SetInt(ctx context.Context, rs RuleSession, name string, value int) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.dataSource, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}
func (st *streamTupleImpl) SetFloat(ctx context.Context, rs RuleSession, name string, value float64) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.dataSource, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}
func (st *streamTupleImpl) SetDatetime(ctx context.Context, rs RuleSession, name string, value time.Time) {
	if rs == nil || rs != nil && !rs.ValidateUpdate(st.dataSource, name, value) {
		return
	}
	if st.tuples[name] != value {
		st.tuples[name] = value
		callChangeListener(ctx, st, name)
	}
}

func callChangeListener(ctx context.Context, tuple StreamTuple, prop string) {
	if ctx != nil {
		ctxR := ctx.Value(reteCTXKEY)
		if ctxR != nil {
			valChangeLister := ctxR.(ValueChangeHandler)
			valChangeLister.OnValueChange(tuple, prop)
		}
	}
}

func (st *streamTupleImpl) GetProperties() []string {
	keys := []string{}
	for k := range st.tuples {
		keys = append(keys, k)
	}
	return keys
}