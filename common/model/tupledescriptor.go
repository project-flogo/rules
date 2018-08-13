package model

import (
	"encoding/json"
	"sort"
	"sync"
	"bytes"
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"strconv"
)

var (
	typeRegistry sync.Map
)

type TupleDescriptor struct {
	Name         string `json:"name"`
	TTLInSeconds int    `json:"ttl"`
	Props        map[string]TuplePropertyDescriptor `json:"props"`
	keyProps     []string
}

type TuplePropertyDescriptor struct {
	Name     string `json:"-"`
	PropType data.Type `json:"type"`
	KeyIndex int `json:"index"`
}

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

func (tpd TuplePropertyDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"" + "type" +"\"" + ":")
	typestr := "\"" + tpd.PropType.String() + "\""
	buffer.WriteString(typestr +",")
	buffer.WriteString("\"" + "pk-index" +"\"" + ":")
	buffer.WriteString(strconv.Itoa(tpd.KeyIndex))
	s := "}"
	buffer.WriteString(s)
	return buffer.Bytes(), nil
}
func (t *TupleDescriptor) UnmarshalJSON(data []byte) error {
	type alias TupleDescriptor
	ata := &alias{}
	ata.TTLInSeconds = -1

	_ = json.Unmarshal(data, ata)

	*t = TupleDescriptor(*ata)
	return nil
}

//func (tpd *TuplePropertyDescriptor) UnmarshalJSON(b []byte) error {
//	type alias TuplePropertyDescriptor
//	ata := &alias{}
//	ata.KeyIndex = -1
//
//	_ = json.Unmarshal(b, ata)
//
//	*tpd = TuplePropertyDescriptor(*ata)
//	return nil
//}

func (tpd *TuplePropertyDescriptor) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	tpd.PropType, _ = data.ToTypeEnum(v["type"].(string))
	kidx, found := v["pk-index"]
	if !found {
		tpd.KeyIndex = -1
	} else {
		tpd.KeyIndex = int(kidx.(float64))
	}
	return nil
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
