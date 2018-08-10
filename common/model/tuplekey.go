package model

import (
	"sort"
	"strconv"
)

type tupleKeyImpl struct {
	td   TupleDescriptor
	keys map[string]interface{}
}

func (tk tupleKeyImpl) String() string {

	s := []string{}
	for k, _ := range tk.keys {
		s = append(s, k)
	}
	sort.Strings(s)
	k := ""
	i := 0
	for i = 0; i < len(s); i++ {
		ky := s[i]
		k = k + ky + ":"
		val := tk.keys[ky]

		switch v := val.(type) {
		case int:
			k += strconv.Itoa(v)
		case string:
			k += v
		case uint64:
			k += strconv.FormatUint(v, 10)
		}
		if i < len(s)-1 {
			k += ","
		}
	}
	return k
}

func (tk tupleKeyImpl) GetTupleDescriptor() TupleDescriptor {
	return tk.td
}

func newTupleKey(tupleType TupleType) TupleKey {
	td := GetTupleDescriptor(tupleType)
	if td == nil {
		return nil
	}
	key := tupleKeyImpl{}
	return &key
}