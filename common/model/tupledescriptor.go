package model

import (
	"bytes"
	"encoding/json"
	"sort"
	"strconv"
	"sync"

	"fmt"

	"github.com/project-flogo/core/data"
)

var (
	typeRegistry sync.Map
)

//TupleType Each tuple is of a certain type, described by TypeDescriptor
type TupleType string

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
func RegisterTupleDescriptors(jsonRegistry string) (err error) {
	tds := []TupleDescriptor{}
	err = json.Unmarshal([]byte(jsonRegistry), &tds)
	if err != nil {
		return err
	}
	for _, key := range tds {
		typeRegistry.LoadOrStore(TupleType(key.Name), key)
	}
	return nil
}

// RegisterTupleDescriptors registers the TupleDescriptors
func RegisterTupleDescriptorsFromTds(tds []TupleDescriptor) (err error) {
	if err != nil {
		return err
	}
	for _, key := range tds {
		typeRegistry.LoadOrStore(TupleType(key.Name), key)
	}
	return nil
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

	jsonProps := val["properties"].([]interface{})

	idxProp := make(map[int]string)
	for _, v := range jsonProps {
		tdp := TuplePropertyDescriptor{}
		tdp.KeyIndex = -1
		pm := v.(map[string]interface{})

		//ensure you get the name first
		for pn, pv := range pm {
			if pn == "name" {
				tdp.Name = pv.(string)
				break
			}
		}
		//duplicate pk-index validation
		for pn, pv := range pm {
			if pn == "type" {
				tdp.PropType, _ = data.ToTypeEnum(pv.(string))
			} else if pn == "pk-index" {
				idx := int(pv.(float64))
				if idx != -1 {
					prop, exists := idxProp[idx]
					if exists {
						return fmt.Errorf("Property [%s] already defined as key at index [%d] for type [%s]",
							prop, idx, nm)
					}
					idxProp[idx] = tdp.Name
				}
				tdp.KeyIndex = idx
			}
		}
		td.Props = append(td.Props, tdp)
	}

	//index validation
	idsx := make([]int, 0)
	for k := range idxProp {
		idsx = append(idsx, k)
	}
	sort.Ints(idsx)
	for i := 0; i < len(idsx); i++ {
		if idsx[i] != i {
			return fmt.Errorf("Missing key at index [%d]", i)
		}
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
