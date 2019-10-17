package redis

import (
	"context"
	"strconv"

	"github.com/project-flogo/rules/redisutils"
	"github.com/project-flogo/rules/rete/common"
	"github.com/project-flogo/rules/rete/internal/types"
)

type jtRefsServiceImpl struct {
	types.NwServiceImpl
	redisutils.RedisHdl
}

func NewJoinTableRefsInHdlImpl(nw types.Network, config common.Config) types.JtRefsService {
	hdlJt := jtRefsServiceImpl{
		NwServiceImpl: types.NwServiceImpl{
			Nw: nw,
		},
		RedisHdl: redisutils.NewRedisHdl(config.Jts.Redis),
	}
	return &hdlJt
}

func (h *jtRefsServiceImpl) Init() {

}

func (h *jtRefsServiceImpl) AddEntry(handle types.ReteHandle, jtName string, rowID int) {
	// format: prefix:rtbls:tkey ==> {rowID=jtname, ...}
	key := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	valMap := make(map[string]interface{})
	valMap[strconv.Itoa(rowID)] = jtName
	h.HSetAll(key, valMap)
}

func (h *jtRefsServiceImpl) RemoveEntry(handle types.ReteHandle, jtName string, rowID int) {
	key := h.Nw.GetPrefix() + ":rtbls:" + handle.GetTupleKey().String()
	h.HDel(key, strconv.Itoa(rowID))
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
	r.iter = h.GetMapIterator(r.key)
	return &r
}
