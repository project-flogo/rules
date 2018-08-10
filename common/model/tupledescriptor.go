package model

import (
	"encoding/json"
	"sort"
	"sync"
)

var (
	typeRegistry sync.Map
)

func RegisterTupleDescriptors(jsonRegistry string) {
	tds := []TupleDescriptor{}
	json.Unmarshal([]byte(jsonRegistry), &tds)
	for _, key := range tds {
		typeRegistry.LoadOrStore(TupleType(key.Name), key)
	}
}

func GetTupleDescriptor(tupleType TupleType) *TupleDescriptor {
	tdi, found := typeRegistry.Load(tupleType)
	if found {
		td := tdi.(TupleDescriptor)
		return &td
	} else {
		return nil
	}
}

type TupleDescriptor struct {
	Name         string //`json:"alias"`
	TTLInSeconds int    //`json:"timeout"`
	Props        map[string]TuplePropertyDescriptor
	keyProps     []string
}

type TuplePropertyDescriptor struct {
	Name     string
	PropType string
	KeyIndex int //index position of this property in a compound key
}

func (td *TupleDescriptor) GetProperty(prop string) (TuplePropertyDescriptor, bool) {
	p, ok := td.Props[prop]
	return p, ok
}

func (td *TupleDescriptor) GetKeyProps() []string {
	if td.keyProps == nil {
		keyProps := []string{}
		keysmap := make(map[int]string)
		keys := []int{}
		for propNm, tdp := range td.Props {
			if tdp.KeyIndex != -1 {
				keysmap[tdp.KeyIndex] = propNm
				keys = append(keys, tdp.KeyIndex)
			}
		}
		sort.Ints(keys)
		for k, _ := range keys {
			keyProps = append(keyProps, keysmap[k])
		}
		td.keyProps = keyProps
	}
	return td.keyProps
}
