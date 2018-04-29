package rete

import (
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
	assert(tuple model.StreamTuple)
}

type classNodeImpl struct {
	classNodeLinks utils.ArrayList
	name           string
}

func newClassNode(name string) classNode {
	classNodeImplVar := classNodeImpl{}
	classNodeImplVar.initClassNodeImpl(name)
	return &classNodeImplVar
}

func (classNodeImplVar *classNodeImpl) initClassNodeImpl(name string) {
	classNodeImplVar.name = name
	classNodeImplVar.classNodeLinks = utils.NewArrayList()
}

func (classNodeImplVar *classNodeImpl) addClassNodeLink(classNodeLinkVar classNodeLink) {
	classNodeImplVar.classNodeLinks.Add(classNodeLinkVar)
}

func (classNodeImplVar *classNodeImpl) removeClassNodeLink(classNodeLinkInList classNodeLink) {

	for i := 0; i < classNodeImplVar.getClassNodeLinks().Len(); i++ {
		classNodeLinkInList := classNodeImplVar.getClassNodeLinks().Get(i)
		if classNodeLinkInList != nil && classNodeLinkInList == classNodeLinkInList {
			classNodeImplVar.getClassNodeLinks().RemoveAt(i)
			break
		}
	}
}

func (classNodeImplVar *classNodeImpl) getClassNodeLinks() utils.ArrayList {
	return classNodeImplVar.classNodeLinks
}

func (classNodeImplVar *classNodeImpl) getName() string {
	return classNodeImplVar.name
}

//Implements Stringer.String
func (classNodeImplVar *classNodeImpl) String() string {
	links := ""
	for i := classNodeImplVar.classNodeLinks.Len() - 1; i >= 0; i-- {
		nl := classNodeImplVar.classNodeLinks.Get(i).(classNodeLink)
		links += "\t" + nl.String()
		if i != classNodeImplVar.classNodeLinks.Len()-1 {
			links += "\n"
		}
	}
	ret := "[ClassNode Class(" + classNodeImplVar.name + ")"
	if len(links) > 0 {
		ret += "\n" + links + "]" + "\n"
	} else {
		ret += "]" + "\n"
	}
	return ret
}

func (classNodeImplVar *classNodeImpl) assert(tuple model.StreamTuple) {
	handle := getOrCreateHandle(tuple)

	handles := make([]reteHandle, 1)
	handles[0] = handle

	for i := 0; i < classNodeImplVar.getClassNodeLinks().Len(); i++ {
		classNodeLinkVar := classNodeImplVar.getClassNodeLinks().Get(i).(classNodeLink)
		classNodeLinkVar.propagateObjects(handles)
	}

}
