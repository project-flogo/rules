package model

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"fmt"
	"reflect"
)

// TupleKey primary key of a tuple
type TupleKey interface {
	String() string
	GetTupleDescriptor() TupleDescriptor
}

type tupleKeyImpl struct {
	td       TupleDescriptor
	keys     map[string]interface{}
	keyAsStr string
}

func (tk *tupleKeyImpl) String() string {

	if tk.keyAsStr == "" {
		i := 0
		keysLen := len(tk.td.GetKeyProps())
		tk.keyAsStr += tk.td.Name + ":"
		for i = 0; i < keysLen; i++ {
			ky := tk.td.GetKeyProps()[i]
			tk.keyAsStr = tk.keyAsStr + ky + ":"
			val := tk.keys[ky]
			strval, _ := data.CoerceToString(val)
			tk.keyAsStr += strval
			if i < keysLen-1 {
				tk.keyAsStr += ","
			}
		}
	}
	return tk.keyAsStr
}

func (tk *tupleKeyImpl) GetTupleDescriptor() TupleDescriptor {
	return tk.td
}

func NewTupleKey(tupleType TupleType, values map[string]interface{}) (tupleKey TupleKey, err error) {
	td := GetTupleDescriptor(tupleType)
	if td == nil {
		return nil, fmt.Errorf("Tuple descriptor not found [%s]", string(tupleType))
	}
	t := tupleKeyImpl{}
	err = t.initTupleKey(td, values)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func NewTupleKeyWithKeyValues(tupleType TupleType, values ...interface{}) (mtuple MutableTuple, err error) {

	td := GetTupleDescriptor(tupleType)
	if td == nil {
		return nil, fmt.Errorf("Tuple descriptor not found [%s]", string(tupleType))
	}
	t := tupleImpl{}
	err = t.initTupleWithKeyValues(td, values...)
	if err != nil {
		return nil, err
	}
	return &t, err
}

func (tk *tupleKeyImpl) initTupleKey(td *TupleDescriptor, values map[string]interface{}) (err error) {

	tk.keys = make(map[string]interface{})
	tk.td = *td
	for _, tdp := range td.Props {
		if tdp.KeyIndex != -1 {
			val, found := values[tdp.Name]
			if found {
				coerced, err := data.CoerceToValue(val, tdp.PropType)
				if err == nil {
					tk.keys[tdp.Name] = coerced
				} else {
					return err
				}
			} else if tdp.KeyIndex != -1 { //key prop
				return fmt.Errorf("Key property [%s] not found", tdp.Name)
			}
		}
	}
	return err
}

func (tk *tupleKeyImpl) initTupleWithKeyValues(td *TupleDescriptor, values ...interface{}) (err error) {

	tk.keys = make(map[string]interface{})
	tk.td = *td
	if len(values) != len(td.GetKeyProps()) {
		return fmt.Errorf("Wrong number of key values in type [%s]. Expecting [%d], got [%d]",
			td.Name, len(td.GetKeyProps()), len(values))
	}

	i := 0
	for _, keyProp := range td.GetKeyProps() {
		tdp := td.GetProperty(keyProp)
		val := values[i]
		coerced, err := data.CoerceToValue(val, tdp.PropType)
		if err == nil {
			tk.keys[keyProp] = coerced
		} else {
			return fmt.Errorf("Type mismatch for field [%s] in type [%s] Expecting [%s], got [%v]",
				keyProp, td.Name, tdp.PropType.String(), reflect.TypeOf(val))
		}
		i++
	}
	return err
}
