package rete

import (
	"context"

	"github.com/TIBCOSoftware/bego/common/model"
	"github.com/TIBCOSoftware/bego/utils"
)

//classNode holds links to filter and join nodes eventually leading upto the rule node
type classNode interface {
	abstractNode
	getName() string
	addClassNodeLink(classNodeLink)
	removeClassNodeLink(classNodeLink)
	getClassNodeLinks() utils.ArrayList
	assert(ctx context.Context, tuple model.StreamTuple)
}

type classNodeImpl struct {
	classNodeLinks utils.ArrayList
	name           string
}

func newClassNode(name string) classNode {
	cn := classNodeImpl{}
	cn.initClassNodeImpl(name)
	return &cn
}

func (cn *classNodeImpl) initClassNodeImpl(name string) {
	cn.name = name
	cn.classNodeLinks = utils.NewArrayList()
}

func (cn *classNodeImpl) addClassNodeLink(classNodeLinkVar classNodeLink) {
	cn.classNodeLinks.Add(classNodeLinkVar)
}

func (cn *classNodeImpl) removeClassNodeLink(classNodeLinkInList classNodeLink) {

	for i := 0; i < cn.getClassNodeLinks().Len(); i++ {
		classNodeLinkInList := cn.getClassNodeLinks().Get(i)
		if classNodeLinkInList != nil && classNodeLinkInList == classNodeLinkInList {
			cn.getClassNodeLinks().RemoveAt(i)
			break
		}
	}
}

func (cn *classNodeImpl) getClassNodeLinks() utils.ArrayList {
	return cn.classNodeLinks
}

func (cn *classNodeImpl) getName() string {
	return cn.name
}

//Implements Stringer.String
func (cn *classNodeImpl) String() string {
	links := ""
	for i := cn.classNodeLinks.Len() - 1; i >= 0; i-- {
		nl := cn.classNodeLinks.Get(i).(classNodeLink)
		links += "\t" + nl.String()
		if i != cn.classNodeLinks.Len()-1 {
			links += "\n"
		}
	}
	ret := "[ClassNode Class(" + cn.name + ")"
	if len(links) > 0 {
		ret += "\n" + links + "]" + "\n"
	} else {
		ret += "]" + "\n"
	}
	return ret
}

func (cn *classNodeImpl) assert(ctx context.Context, tuple model.StreamTuple) {
	handle := getOrCreateHandle(tuple)

	handles := make([]reteHandle, 1)
	handles[0] = handle

	for i := 0; i < cn.getClassNodeLinks().Len(); i++ {
		classNodeLinkVar := cn.getClassNodeLinks().Get(i).(classNodeLink)
		classNodeLinkVar.propagateObjects(ctx, handles)
	}

}
