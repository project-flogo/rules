package redis

import (
	"github.com/project-flogo/rules/rete/internal/types"
	"github.com/project-flogo/rules/redisutils"
	"strings"
	"github.com/project-flogo/rules/common/model"
	"strconv"
)

type rowIteratorImpl struct {
	iter *redisutils.MapIterator
	jtName string
	nw types.Network
}

func newRowIterator(jTable types.JoinTable) types.RowIterator {
	key := jTable.GetNw().GetPrefix() + "jt:" + jTable.GetName()
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

	rowId, value := ri.iter.Next()
	rowID, _ := strconv.Atoi(rowId)

	strval := value.(string)
	values := strings.Split(strval, ",")


	handles := []types.ReteHandle{}
	for _, key := range values {
		tupleKey := model.FromStringKey(key)
		tuple := ri.nw.GetTupleStore().GetTupleByKey(tupleKey)
		handle := newReteHandleImpl(ri.nw, tuple)
		handles = append (handles, handle)
	}

	jtRow := newJoinTableRowLoadedFromStore(ri.jtName, rowID, handles, ri.nw)

	return jtRow

}
