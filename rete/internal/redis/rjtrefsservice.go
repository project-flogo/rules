package redis

import (
	"context"
	"strconv"

	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type jtRefsServiceImpl struct {
	//keys are jointable-ids and values are lists of row-ids in the corresponding join table
	types.NwServiceImpl

	//tablesAndRows map[string]map[string]map[int]int
}

func NewJoinTableRefsInHdlImpl(nw types.Network, config common.Config) types.JtRefsService {
	hdlJt := jtRefsServiceImpl{}
	hdlJt.Nw = nw
	//hdlJt.tablesAndRows = make(map[string]map[string]map[int]int)
	return &hdlJt
}

func (h *jtRefsServiceImpl) Init() {

}

func (h *jtRefsServiceImpl) AddEntry(handle types.ReteHandle, jtName string, rowID int) {
	// format: prefix:rtbls:tkey ==> {rowID=jtname, ...}
	key := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	hdl := redisutils.GetRedisHdl()
	valMap := make(map[string]interface{})
	valMap[strconv.Itoa(rowID)] = jtName
	hdl.HSetAll(key, valMap)
}

func (h *jtRefsServiceImpl) RemoveEntry(handle types.ReteHandle, jtName string, rowID int) {
	key := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	hdl := redisutils.GetRedisHdl()
	hdl.HDel(key, strconv.Itoa(rowID))
}

type hdlRefsRowIteratorImpl struct {
	ctx  context.Context
	key  string
	iter *redisutils.MapIterator
	nw   types.Network
}

func (r *hdlRefsRowIteratorImpl) HasNext() bool {
	return r.iter.HasNext()
}

func (r *hdlRefsRowIteratorImpl) Next() (types.JoinTableRow, types.JoinTable) {
	rowIDStr, jtName := r.iter.Next()
	rowID, _ := strconv.Atoi(rowIDStr)
	jT := r.nw.GetJtService().GetJoinTable(jtName.(string))
	if jT != nil {
		return jT.GetRow(r.ctx, rowID), jT
	}
	return nil, jT
}

func (r *hdlRefsRowIteratorImpl) Remove() {
	r.iter.Remove()
}

func (h *jtRefsServiceImpl) GetRowIterator(ctx context.Context, handle types.ReteHandle) types.JointableIterator {
	r := hdlRefsRowIteratorImpl{}
	r.ctx = ctx
	r.nw = h.Nw
	r.key = h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	hdl := redisutils.GetRedisHdl()
	r.iter = hdl.GetMapIterator(r.key)
	return &r
}
