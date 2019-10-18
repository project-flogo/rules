package redis

import (
	"context"
	"strconv"
	"strings"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
)

type joinTableRowImpl struct {
	types.NwElemIdImpl
	handles []types.ReteHandle
	jtKey   string
	redisutils.RedisHdl
}

func newJoinTableRow(handle redisutils.RedisHdl, jtKey string, handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{
		RedisHdl: handle,
		jtKey:    jtKey,
		handles:  append([]types.ReteHandle{}, handles...),
	}
	jtr.SetID(nw)
	return &jtr
}

func newJoinTableRowLoadedFromStore(handle redisutils.RedisHdl, jtKey string, rowID int, handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{
		handles: handles,
		jtKey:   jtKey,
		NwElemIdImpl: types.NwElemIdImpl{
			ID: rowID,
			Nw: nw,
		},
		RedisHdl: handle,
	}
	return &jtr
}

func (jtr *joinTableRowImpl) Write() {
	handles, row := jtr.handles, make(map[string]interface{})
	end, str := len(handles)-1, ""
	for i, v := range handles {
		str += v.GetTupleKey().String()
		if i < end {
			str += ","
		}
	}
	row[strconv.Itoa(jtr.ID)] = str
	jtr.HSetAll(jtr.jtKey, row)
}

func (jtr *joinTableRowImpl) GetHandles() []types.ReteHandle {
	return jtr.handles
}

func createRow(ctx context.Context, handle redisutils.RedisHdl, jtKey string, rowID string, key string, nw types.Network) types.JoinTableRow {

	values := strings.Split(key, ",")

	handles := []types.ReteHandle{}
	for _, key := range values {
		tupleKey := model.FromStringKey(key)
		var tuple model.Tuple
		if ctx != nil {
			if value := ctx.Value(model.RetecontextKeyType{}); value != nil {
				if value, ok := value.(types.ReteCtx); ok {
					if modified := value.GetRtcModified(); modified != nil {
						if value := modified[tupleKey.String()]; value != nil {
							tuple = value.GetTuple()
						}
					}
					if tuple == nil {
						if added := value.GetRtcAdded(); added != nil {
							tuple = added[tupleKey.String()]
						}
					}
				}
			}
		}
		if tuple == nil {
			tuple = nw.GetTupleStore().GetTupleByKey(tupleKey)
		}
		handle := newReteHandleImpl(nw, handle, tuple, "", types.ReteHandleStatusUnknown, 0)
		handles = append(handles, handle)
	}

	rowId, _ := strconv.Atoi(rowID)
	jtRow := newJoinTableRowLoadedFromStore(handle, jtKey, rowId, handles, nw)

	return jtRow
}
