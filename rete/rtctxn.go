package rete

import "github.com/project-flogo/rules/common/model"

type rtcTxnImpl struct {
	added    map[string]map[string]model.Tuple
	modified map[string]map[string]model.RtcModified
	deleted  map[string]map[string]model.Tuple
}

func newRtcTxn(addedTxn map[string]model.Tuple, modifiedTxn map[string]model.RtcModified, deletedTxn map[string]model.Tuple) model.RtcTxn {
	rtxn := rtcTxnImpl{}
	rtxn.init(addedTxn, modifiedTxn, deletedTxn)
	return &rtxn
}

func (tx *rtcTxnImpl) init(addedTxn map[string]model.Tuple, modifiedTxn map[string]model.RtcModified, deletedTxn map[string]model.Tuple) {
	tx.added = make(map[string]map[string]model.Tuple)
	tx.modified = make(map[string]map[string]model.RtcModified)
	tx.deleted = make(map[string]map[string]model.Tuple)

	tx.groupAddedByType(addedTxn)

	tx.groupModifiedByType(modifiedTxn)

	tx.groupDeletedByType(deletedTxn)
}

func (tx *rtcTxnImpl) groupDeletedByType(deletedTxn map[string]model.Tuple) {
	for key, tuple := range deletedTxn {
		tdType := tuple.GetTupleDescriptor().Name
		tupleMap, found := tx.deleted[tdType]
		if !found {
			tupleMap = make(map[string]model.Tuple)
			tx.deleted[tdType] = tupleMap
		}
		tupleMap[key] = tuple
	}
}

func (tx *rtcTxnImpl) groupModifiedByType(modifiedTxn map[string]model.RtcModified) {
	for key, rtcModified := range modifiedTxn {
		tdType := rtcModified.GetTuple().GetTupleDescriptor().Name
		tupleMap, found := tx.modified[tdType]
		if !found {
			tupleMap = make(map[string]model.RtcModified)
			tx.modified[tdType] = tupleMap
		}
		tupleMap[key] = rtcModified
	}
}

func (tx *rtcTxnImpl) groupAddedByType(addedTxn map[string]model.Tuple) {
	for key, tuple := range addedTxn {
		tdType := tuple.GetTupleDescriptor().Name
		tupleMap, found := tx.added[tdType]
		if !found {
			tupleMap = make(map[string]model.Tuple)
			tx.added[tdType] = tupleMap
		}
		tupleMap[key] = tuple
	}
}

func (tx *rtcTxnImpl) GetRtcAdded() map[string]map[string]model.Tuple {
	return tx.added
}

func (tx *rtcTxnImpl) GetRtcModified() map[string]map[string]model.RtcModified {
	return tx.modified
}

func (tx *rtcTxnImpl) GetRtcDeleted() map[string]map[string]model.Tuple {
	return tx.deleted
}
