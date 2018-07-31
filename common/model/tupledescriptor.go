package model

type TupleDescriptor struct {
	Name string //`json:"alias"`
	Expiry int //`json:"timeout"`
	Props map[string]TuplePropertyDescriptor
}

type TuplePropertyDescriptor struct {
	Name string
	PropType string
}

func (td *TupleDescriptor) GetProperty(prop string) (TuplePropertyDescriptor, bool) {
	p, ok := td.Props[prop]
	return p, ok
}






