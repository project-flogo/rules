package rete

import (
	"strconv"
)

//nodelink connects 2 nodes, a rete building block
type nodeLink interface {
	String() string
	getChild() node
	isRightNode() bool

	setChild(child node)
	setIsRightChild(isRight bool)
	propagateObjects(handles []reteHandle)
}

type nodeLinkImpl struct {
	convert        []int
	numIdentifiers int

	parent    node
	parentIds []identifier

	child    node
	childIds []identifier

	isRight bool
	id      int
}

func newNodeLink(parent node, child node, isRight bool) nodeLink {
	nodeLinkImplVar := nodeLinkImpl{}
	nodeLinkImplVar.initNodeLink(parent, child, isRight)
	return &nodeLinkImplVar
}

func (nodeLinkImplVar *nodeLinkImpl) initNodeLink(parent node, child node, isRight bool) {
	nodeLinkImplVar.id = currentNodeID
	nodeLinkImplVar.child = child
	nodeLinkImplVar.isRight = isRight

	switch v := child.(type) {

	case *joinNodeImpl:
		if isRight {
			nodeLinkImplVar.childIds = v.rightIdrs
		} else {
			nodeLinkImplVar.childIds = v.leftIdrs
		}
	case *nodeImpl:
		nodeLinkImplVar.childIds = v.identifiers
	}
	nodeLinkImplVar.parent = parent
	nodeLinkImplVar.setConvert()
	parent.addNodeLink(nodeLinkImplVar)
}

//initialize node link : for use with ClassNodeLink
func initClassNodeLink(nodeLinkImplVar *nodeLinkImpl, child node) {
	currentNodeID++
	nodeLinkImplVar.id = currentNodeID
	nodeLinkImplVar.child = child
	nodeLinkImplVar.childIds = child.getIdentifiers()
}

func (nodeLinkImplVar *nodeLinkImpl) getChild() node {
	return nodeLinkImplVar.child
}

func (nodeLinkImplVar *nodeLinkImpl) setConvert() {

	if len(nodeLinkImplVar.parentIds) != len(nodeLinkImplVar.childIds) {
		//TODO: ERROR handling
	}
	nodeLinkImplVar.numIdentifiers = len(nodeLinkImplVar.parentIds)
	nodeLinkImplVar.convert = make([]int, nodeLinkImplVar.numIdentifiers)

	for i := 0; i < nodeLinkImplVar.numIdentifiers; i++ {
		found := false
		for j := 0; j < nodeLinkImplVar.numIdentifiers; j++ {
			if nodeLinkImplVar.parentIds[i].equals(nodeLinkImplVar.childIds[j]) {
				found = true
				nodeLinkImplVar.convert[i] = j
				break
			}
		}
		if !found {
			//TODO: ERROR handling
		}
	}

	need := false
	for i := 0; i < nodeLinkImplVar.numIdentifiers; i++ {
		if nodeLinkImplVar.convert[i] != i {
			need = true
			break
		}
	}
	if !need {
		nodeLinkImplVar.convert = nil
	}
}

func (nodeLinkImplVar *nodeLinkImpl) String() string {
	nextNode := ""
	switch nodeLinkImplVar.child.(type) {
	case *joinNodeImpl:
		if nodeLinkImplVar.isRight {
			nextNode += "j" + strconv.Itoa(nodeLinkImplVar.child.getID()) + "R"
		} else {
			nextNode += "j" + strconv.Itoa(nodeLinkImplVar.child.getID()) + "L"
		}
	case *filterNodeImpl:
		nextNode += "f" + strconv.Itoa(nodeLinkImplVar.child.getID())
	}
	return "link (" + nextNode + ")"
}

func (nodeLinkImplVar *nodeLinkImpl) isRightNode() bool {
	return nodeLinkImplVar.isRight
}

func (nodeLinkImplVar *nodeLinkImpl) setChild(child node) {
	nodeLinkImplVar.child = child
}
func (nodeLinkImplVar *nodeLinkImpl) setIsRightChild(isRight bool) {
	nodeLinkImplVar.isRight = isRight
}

func (nodeLinkImplVar *nodeLinkImpl) propagateObjects(handles []reteHandle) {
	if nodeLinkImplVar.convert != nil {
		convertedHandles := make([]reteHandle, nodeLinkImplVar.numIdentifiers)
		for i := 0; i < nodeLinkImplVar.numIdentifiers; i++ {
			convertedHandles[nodeLinkImplVar.convert[i]] = handles[i]
		}
		handles = convertedHandles
	}
	nodeLinkImplVar.child.assertObjects(handles, nodeLinkImplVar.isRightNode())
}
