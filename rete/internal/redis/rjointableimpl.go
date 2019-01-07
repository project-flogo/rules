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

func (jt *joinTableImpl) RemoveAllRows() {
	rowIter := jt.GetRowIterator()
	for rowIter.HasNext() {
		row := rowIter.Next()
		//first, from jTable, remove row
		jt.RemoveRow(row.GetID())
		for _, hdl := range row.GetHandles() {
			jt.Nw.GetJtRefService().RemoveTableEntry(hdl, jt.GetName())
		}
		//Delete the rowRef itself
		rowIter.Remove()
	}
}

func (jt *joinTableImpl) GetRowCount() int {
	hdl := redisutils.GetRedisHdl()
	return hdl.HLen(jt.name)
}

func (jt *joinTableImpl) GetRule() model.Rule {
	return jt.rule
}

func (jt *joinTableImpl) GetRowIterator() types.JointableRowIterator {
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

type rowIteratorImpl struct {
	iter   *redisutils.MapIterator
	jtName string
	nw     types.Network
	curr   types.JoinTableRow
}

func newRowIterator(jTable types.JoinTable) types.JointableRowIterator {
	key := jTable.GetNw().GetPrefix() + ":jt:" + jTable.GetName()
	ri := rowIteratorImpl{}
	ri.iter = redisutils.GetRedisHdl().GetMapIterator(key)
	ri.nw = jTable.GetNw()
	ri.jtName = jTable.GetName()
	return &ri
}

func (ri *rowIteratorImpl) HasNext() bool {
	return ri.iter.HasNext()
}

func (ri *rowIteratorImpl) Next() types.JoinTableRow {
	rowId, key := ri.iter.Next()
	tupleKeyStr := key.(string)
	ri.curr = createRow(ri.jtName, rowId, tupleKeyStr, ri.nw)
	return ri.curr
}

func (ri *rowIteratorImpl) Remove() {
	ri.iter.Remove()
}
