package rete

import "github.com/project-flogo/rules/common/model"

type rtcModifiedImpl struct {
	tuple model.Tuple
	props map[string]bool
}

func NewRtcModified(tuple model.Tuple) model.RtcModified {
	rm := rtcModifiedImpl{}
	rm.tuple = tuple
	rm.props = make(map[string]bool)
	return &rm
}

func (rm *rtcModifiedImpl) GetTuple() model.Tuple {
	return rm.tuple
}

func (rm *rtcModifiedImpl) GetModifiedProps() map[string]bool {
	return rm.props
}

func (rm *rtcModifiedImpl) addProp(prop string) {
	rm.props[prop] = true
}
