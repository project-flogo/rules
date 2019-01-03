package redis

import (
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
	"strconv"
)

type joinTableImpl struct {
	types.NwElemIdImpl
	//table map[int]types.JoinTableRow
	idr   []model.TupleType
	rule  model.Rule
	name  string
	jtKey string
}

func newJoinTableImpl(nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) types.JoinTable {
	jt := joinTableImpl{}
	jt.initJoinTableImpl(nw, rule, identifiers, name)
	return &jt
}

func (jt *joinTableImpl) initJoinTableImpl(nw types.Network, rule model.Rule, identifiers []model.TupleType, name string) {
	jt.SetID(nw)
	jt.idr = identifiers
	jt.rule = rule
	jt.name = name
	jt.jtKey = nw.GetPrefix() + ":" + "jt:" + name
}

func (jt *joinTableImpl) AddRow(handles []types.ReteHandle) types.JoinTableRow {
	row := newJoinTableRow(jt.jtKey, handles, jt.Nw)
	for i := 0; i < len(row.GetHandles()); i++ {
		handle := row.GetHandles()[i]
		jt.Nw.GetJtRefService().AddEntry(handle, jt.name, row.GetID())
	}
	return row
}

func (jt *joinTableImpl) RemoveRow(rowID int) types.JoinTableRow {
	row := jt.GetRow(rowID)
	hdl := redisutils.GetRedisHdl()
	rowId := strconv.Itoa(rowID)
	hdl.HDel(jt.jtKey, rowId)
	return row
}

func (jt *joinTableImpl) GetRowCount() int {
	hdl := redisutils.GetRedisHdl()
	return hdl.HLen(jt.name)
}

func (jt *joinTableImpl) GetRule() model.Rule {
	return jt.rule
}

func (jt *joinTableImpl) GetRowIterator() types.RowIterator {
	return newRowIterator(jt)
}

func (jt *joinTableImpl) GetRow(rowID int) types.JoinTableRow {
	hdl := redisutils.GetRedisHdl()
	key := hdl.HGet(jt.jtKey, strconv.Itoa(rowID))
	rowId := strconv.Itoa(rowID)
	return createRow(jt.name, rowId, key.(string), jt.Nw)
}

func (jt *joinTableImpl) GetName() string {
	return jt.name
}
