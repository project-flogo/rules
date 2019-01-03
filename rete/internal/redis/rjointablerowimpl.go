package redis

import (
	"fmt"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/internal/types"
	"strconv"
	"strings"
)

type joinTableRowImpl struct {
	types.NwElemIdImpl
	handles []types.ReteHandle
	jtKey   string
}

func newJoinTableRow(jtKey string, handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{}
	jtr.SetID(nw)
	jtr.initJoinTableRow(jtKey, handles, nw)
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

func (jtr *joinTableRowImpl) initJoinTableRow(jtKey string, handles []types.ReteHandle, nw types.Network) {
	jtr.jtKey = jtKey
	jtr.handles = append([]types.ReteHandle{}, handles...)

	rowEntry := make(map[string]interface{})

	str := ""

	for i, v := range handles {
		str += v.GetTupleKey().String()
		if i < len(handles)-1 {
			str += ","
		}
	}

	rowEntry[strconv.Itoa(jtr.ID)] = str

	hdl := redisutils.GetRedisHdl()
	hdl.HSetAll(jtKey, rowEntry)

}

func (jtr *joinTableRowImpl) GetHandles() []types.ReteHandle {
	return jtr.handles
}

func createRow(jtKey string, rowID string, key string, nw types.Network) types.JoinTableRow {

	values := strings.Split(key, ",")

	handles := []types.ReteHandle{}
	for _, key := range values {
		tupleKey := model.FromStringKey(key)
		tuple := nw.GetTupleStore().GetTupleByKey(tupleKey)
		ks := tupleKey.String()
		fmt.Printf(ks)
		handle := newReteHandleImpl(nw, tuple)
		handles = append(handles, handle)
	}

	rowId, _ := strconv.Atoi(rowID)
	jtRow := newJoinTableRowLoadedFromStore(jtKey, rowId, handles, nw)

	return jtRow
}
