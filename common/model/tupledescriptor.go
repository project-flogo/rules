package model

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"sync"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

var (
	typeRegistry sync.Map
)

// TupleDescriptor defines the type of the structure, its properties, types
type TupleDescriptor struct {
	Name         string                    `json:"name"`
	TTLInSeconds int                       `json:"ttl"`
	Props        []TuplePropertyDescriptor `json:"properties"`
	keyProps     []string
}

// TuplePropertyDescriptor defines the actual property, its type, key index
type TuplePropertyDescriptor struct {
	Name     string    `json:"name"`
	PropType data.Type `json:"type"`
	KeyIndex int       `json:"pk-index"`
}

// RegisterTupleDescriptors registers the TupleDescriptors
func RegisterTupleDescriptors(jsonRegistry string) {
	tds := []TupleDescriptor{}
	json.Unmarshal([]byte(jsonRegistry), &tds)
	for _, key := range tds {
		typeRegistry.LoadOrStore(TupleType(key.Name), key)
	}
}

// GetTupleDescriptor gets the TupleDescriptor based on the TupleType
func GetTupleDescriptor(tupleType TupleType) *TupleDescriptor {
	tdi, found := typeRegistry.Load(tupleType)
	if found {
		td := tdi.(TupleDescriptor)
		return &td
	}

	return nil
}

// MarshalJSON allows to hook & customize TupleDescriptor to JSON conversion
func (tpd TuplePropertyDescriptor) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString("{")
	buffer.WriteString("\"" + "name" + "\"" + ":")
	namestr := "\"" + tpd.Name + "\""
	buffer.WriteString(namestr + ",")
	buffer.WriteString("\"" + "type" + "\"" + ":")
	typestr := "\"" + tpd.PropType.String() + "\""
	buffer.WriteString(typestr + ",")
	buffer.WriteString("\"" + "pk-index" + "\"" + ":")
	buffer.WriteString(strconv.Itoa(tpd.KeyIndex))
	s := "}"
	buffer.WriteString(s)
	return buffer.Bytes(), nil
}

func (td *TupleDescriptor) UnmarshalJSON(b []byte) error {

	val := map[string]interface{}{}

	json.Unmarshal(b, &val)


	nm := val["name"]
	//fmt.Printf("%s", nm)

	td.Name = nm.(string)
	td.TTLInSeconds = -1

	ttl, ok := val["ttl"]
	if ok {
		td.TTLInSeconds = int(ttl.(float64))
	}

	pm1 := val["properties"].([]interface{})
	//pm := pm1.(map[string]string)
	//fmt.Printf("%v",  pm1)

	for _, v:= range pm1 {
		tdp := TuplePropertyDescriptor{}
		tdp.KeyIndex = -1
		pm := v.(map[string]interface{})
		for pn, pv := range pm {
			if pn == "name" {
				tdp.Name = pv.(string)
			} else if pn == "type" {
				tdp.PropType, _ = data.ToTypeEnum(pv.(string))
			} else if pn == "pk-index" {
				tdp.KeyIndex = int(pv.(float64))
			}
		}
		td.Props = append (td.Props, tdp)
	}
	return nil
}

// UnmarshalJSON allows to hook & customize JSON to TuplePropertyDescriptor conversion
func (tpd *TuplePropertyDescriptor) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	tpd.Name = v["name"].(string)
	tpd.PropType, _ = data.ToTypeEnum(v["type"].(string))
	kidx, found := v["pk-index"]
	if !found {
		tpd.KeyIndex = -1
	} else {
		tpd.KeyIndex = int(kidx.(float64))
	}
	return nil
}

// GetProperty fetches the property by name
func (td *TupleDescriptor) GetProperty(prop string) *TuplePropertyDescriptor {
	for idx := range td.Props {
		p := td.Props[idx]
		if p.Name == prop {
			return &p
		}
	}
	return nil
}

// GetKeyProps returns all the key properties
func (td *TupleDescriptor) GetKeyProps() []string {
	if td.keyProps == nil {
		keyProps := []string{}
		keysmap := make(map[int]string)
		keys := []int{}
		for idx := range td.Props {
			p := td.Props[idx]
			if p.KeyIndex != -1 {
				keysmap[p.KeyIndex] = p.Name
				keys = append(keys, p.KeyIndex)
			}
		}
		sort.Ints(keys)
		for k := range keys {
			keyProps = append(keyProps, keysmap[k])
		}
		td.keyProps = keyProps
	}
	return td.keyProps
}
