package redis

import (
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

type rowIteratorImpl struct {
	iter   *redisutils.MapIterator
	jtName string
	nw     types.Network
	curr   types.JoinTableRow
}

func newRowIterator(jTable types.JoinTable) types.RowIterator {
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
