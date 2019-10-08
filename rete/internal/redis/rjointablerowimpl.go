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
}

func newJoinTableRow(jtKey string, handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{
		handles: append([]types.ReteHandle{}, handles...),
		jtKey:   jtKey,
	}
	jtr.SetID(nw)
	return &jtr
}

func newJoinTableRowLoadedFromStore(jtKey string, rowID int, handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{}
	jtr.jtKey = jtKey
	jtr.Nw = nw
	jtr.ID = rowID
	jtr.handles = handles

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
	hdl := redisutils.GetRedisHdl()
	hdl.HSetAll(jtr.jtKey, row)
}

func (jtr *joinTableRowImpl) GetHandles() []types.ReteHandle {
	return jtr.handles
}

func createRow(ctx context.Context, jtKey string, rowID string, key string, nw types.Network) types.JoinTableRow {

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
		handle := newReteHandleImpl(nw, tuple, "", types.ReteHandleStatusUnknown)
		handles = append(handles, handle)
	}

	rowId, _ := strconv.Atoi(rowID)
	jtRow := newJoinTableRowLoadedFromStore(jtKey, rowId, handles, nw)

	return jtRow
}
