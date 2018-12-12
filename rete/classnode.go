package rete

import (
	"container/list"
	"context"

	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/rete/internal/types"
)

//classNode holds links to filter and join nodes eventually leading upto the rule node
type classNode interface {
	abstractNode
	getName() string
	addClassNodeLink(classNodeLink)
	removeClassNodeLink(classNodeLink)
	getClassNodeLinks() *list.List
	assert(ctx context.Context, tuple model.Tuple, changedProps map[string]bool)
}

type classNodeImpl struct {
	classNodeLinks *list.List
	name           string
}

func newClassNode(name string) classNode {
	cn := classNodeImpl{}
	cn.initClassNodeImpl(name)
	return &cn
}

func (cn *classNodeImpl) initClassNodeImpl(name string) {
	cn.name = name
	cn.classNodeLinks = list.New()
}

func (cn *classNodeImpl) addClassNodeLink(classNodeLinkVar classNodeLink) {
	cn.classNodeLinks.PushBack(classNodeLinkVar)
}

func (cn *classNodeImpl) removeClassNodeLink(classNodeLinkVar classNodeLink) {

	for e := cn.getClassNodeLinks().Front(); e != nil; e = e.Next() {
		classNodeLinkInList := e.Value
		if classNodeLinkInList != nil && classNodeLinkVar == classNodeLinkInList {
			cn.getClassNodeLinks().Remove(e)
			break
		}
	}
}

func (cn *classNodeImpl) getClassNodeLinks() *list.List {
	return cn.classNodeLinks
}

func (cn *classNodeImpl) getName() string {
	return cn.name
}

//Implements Stringer.String
func (cn *classNodeImpl) String() string {
	links := "\n"

	for e := cn.classNodeLinks.Back(); e != nil; e = e.Prev() {
		nl := e.Value.(classNodeLink)
		links += "\t" + nl.String()
	}

	ret := "[ClassNode Class(" + cn.name + ")"
	if len(links) > 0 {
		ret += "\n" + links + "]" + "\n"
	} else {
		ret += "]" + "\n"
	}
	return ret
}

func (cn *classNodeImpl) assert(ctx context.Context, tuple model.Tuple, changedProps map[string]bool) {
	handle := getOrCreateHandle(ctx, tuple)
	handles := make([]types.ReteHandle, 1)
	handles[0] = handle
	propagate := false
	for e := cn.getClassNodeLinks().Front(); e != nil; e = e.Next() {
		classNodeLinkVar := e.Value.(classNodeLink)
		if changedProps != nil {
			depProps, found := classNodeLinkVar.getRule().GetDeps()[model.TupleType(cn.name)]
			if found { // rule depends on this type
				for changedProp := range changedProps {
					_, foundProp := depProps[changedProp]
					if foundProp {
						propagate = true
						break
					}
				}
			}
		} else {
			propagate = true
		}
		if propagate {
			classNodeLinkVar.propagateObjects(ctx, handles)
		}
	}
}
