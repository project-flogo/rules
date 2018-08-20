package model

import (
	"github.com/TIBCOSoftware/flogo-lib/core/data"
)

// TupleKey primary key of a tuple
type TupleKey interface {
	String() string
	GetTupleDescriptor() TupleDescriptor
}


type tupleKeyImpl struct {
	td   TupleDescriptor
	keys map[string]interface{}
	keyAsStr string
}

func (tk *tupleKeyImpl) String() string {

	if tk.keyAsStr == "" {
		i := 0
		keysLen := len(tk.td.GetKeyProps())
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