package model

import (
	"fmt"
	"reflect"

	"github.com/project-flogo/core/data/coerce"
)

// TupleKey primary key of a tuple
type TupleKey interface {
	String() string
	GetTupleDescriptor() TupleDescriptor
	GetProps() []string
	GetValue(string) interface{}
}

type tupleKeyImpl struct {
	td       TupleDescriptor
	keys     map[string]interface{}
	keyAsStr string
}

func (tk *tupleKeyImpl) String() string {
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

	tk := tupleKeyImpl{}
	tk.td = *td
	tk.keys = make(map[string]interface{})

	for _, tdp := range td.Props {
		if tdp.KeyIndex != -1 {
			val, found := values[tdp.Name]
			if found {
				coerced, err := coerce.ToType(val, tdp.PropType)
				if err == nil {
					tk.keys[tdp.Name] = coerced
				} else {
					return nil, fmt.Errorf("Type mismatch for key field [%s] in type [%s] Expecting [%s], got [%v]",
						tdp.Name, td.Name, tdp.PropType.String(), reflect.TypeOf(val))
				}
			} else if tdp.KeyIndex != -1 { //key prop
				return nil, fmt.Errorf("Key property [%s] not found", tdp.Name)
			}
		}
	}
	tk.keyAsStr = tk.keysAsString()
	return &tk, err
}

func NewTupleKeyWithKeyValues(tupleType TupleType, values ...interface{}) (tupleKey TupleKey, err error) {

	td := GetTupleDescriptor(tupleType)
	if td == nil {
		return nil, fmt.Errorf("Tuple descriptor not found [%s]", string(tupleType))
	}

	tk := tupleKeyImpl{}
	tk.td = *td
	tk.keys = make(map[string]interface{})

	if len(values) != len(td.GetKeyProps()) {
		return nil, fmt.Errorf("Wrong number of key values in type [%s]. Expecting [%d], got [%d]",
			td.Name, len(td.GetKeyProps()), len(values))
	}

	i := 0
	for _, keyProp := range td.GetKeyProps() {
		tdp := td.GetProperty(keyProp)
		val := values[i]
		coerced, err := coerce.ToType(val, tdp.PropType)
		if err == nil {
			tk.keys[keyProp] = coerced
		} else {
			return nil, fmt.Errorf("Type mismatch for field [%s] in type [%s] Expecting [%s], got [%v]",
				keyProp, td.Name, tdp.PropType.String(), reflect.TypeOf(val))
		}
		i++
	}
	tk.keyAsStr = tk.keysAsString()
	return &tk, err
}

func (tk *tupleKeyImpl) GetProps() []string {
	td := tk.GetTupleDescriptor()
	return td.GetKeyProps()
}

func (tk *tupleKeyImpl) GetValue(prop string) interface{} {
	val := tk.keys[prop]
	return val
}

func (tk *tupleKeyImpl) keysAsString() string {
	str := ""
	i := 0
	keysLen := len(tk.td.GetKeyProps())
	str += tk.td.Name + ":"
	for i = 0; i < keysLen; i++ {
		ky := tk.td.GetKeyProps()[i]
		str = str + ky + ":"
		val := tk.keys[ky]
		strval, _ := coerce.ToString(val)
		str += strval
		if i < keysLen-1 {
			str += ","
		}
	}
	return str
}
