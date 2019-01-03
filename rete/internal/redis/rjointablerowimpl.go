package redis

import (
	"github.com/project-flogo/rules/rete/internal/types"
	"strconv"
	"github.com/project-flogo/rules/redisutils"
)

type joinTableRowImpl struct {
	types.NwElemIdImpl
	handles []types.ReteHandle
	jtName string
}

func newJoinTableRow(jtName string, handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{}
	jtr.SetID(nw)
	jtr.initJoinTableRow(jtName, handles, nw)
	return &jtr
}

func newJoinTableRowLoadedFromStore(jtName string, rowID int, handles []types.ReteHandle, nw types.Network) types.JoinTableRow {
	jtr := joinTableRowImpl{}
	jtr.jtName = jtName
	jtr.Nw = nw
	jtr.ID = rowID
	jtr.handles = handles

	return &jtr
}

func (jtr *joinTableRowImpl) initJoinTableRow(jtName string, handles []types.ReteHandle, nw types.Network) {
	jtr.jtName = jtName
	jtr.handles = append([]types.ReteHandle{}, handles...)

	rowEntry := make (map[string]interface{})

	str := ""

	for i, v := range handles {
		str += v.GetTupleKey().String()
		if i < len(handles) - 1 {
			str += ","
		}
	}

	rowEntry[strconv.Itoa(jtr.ID)] = str


	key := nw.GetPrefix() + ":jt:" + jtName

	hdl := redisutils.GetRedisHdl()
	hdl.HSetAll(key, rowEntry)

}

func (jtr *joinTableRowImpl) GetHandles() []types.ReteHandle {
	return jtr.handles
}

