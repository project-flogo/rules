package rete

import (
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/bego/utils"
)

//node a building block of the rete network
type node interface {
	abstractNode
	getIdentifiers() []Identifier
	getID() int
	addNodeLink(nodeLink)
	assertObjects(handles []reteHandle, isRight bool)
}

type nodeImpl struct {
	identifiers []Identifier
	nodeLinkVar nodeLink
	id          int
}

//NewNode ... returns a new node
func newNode(identifiers []Identifier) node {
	nodeImplVar := nodeImpl{}
	nodeImplVar.initNodeImpl(identifiers)
	return &nodeImplVar
}

func (nodeImplVar *nodeImpl) initNodeImpl(identifiers []Identifier) {
	currentNodeID++
	nodeImplVar.id = currentNodeID

	nodeImplVar.identifiers = identifiers
}

func (nodeImplVar *nodeImpl) getIdentifiers() []Identifier {
	return nodeImplVar.identifiers
}

func (nodeImplVar *nodeImpl) getID() int {
	return nodeImplVar.id
}

func (nodeImplVar *nodeImpl) addNodeLink(nl nodeLink) {
	nodeImplVar.nodeLinkVar = nl
}

func (nodeImplVar *nodeImpl) String() string {
	str := "id:" + strconv.Itoa(nodeImplVar.id) + ", idrs:"
	for _, nodeIdentifier := range nodeImplVar.identifiers {
		str += nodeIdentifier.String() + ","
	}
	return str
}

//FindSimilarNodes find similar nodes
func findSimilarNodes(nodeSet utils.ArrayList) []node {
	if nodeSet.Len() < 2 {
		//TODO: Handle error
		return nil
	}
	maxCommon := 0
	similarNodes := make([]node, 2)
	for i := 0; i < nodeSet.Len()-1; i++ {
		node1 := nodeSet.Get(i).(node)
		for j := i + 1; j < nodeSet.Len(); j++ {
			node2 := nodeSet.Get(j).(node)
			common := len(IntersectionIdentifiers(node1.getIdentifiers(), node2.getIdentifiers()))
			if common > maxCommon {
				maxCommon = common
				similarNodes[0] = node1
				similarNodes[1] = node2
			}
		}
	}
	if maxCommon == 0 {
		similarNodes[0] = nodeSet.Get(0).(node)
		similarNodes[1] = nodeSet.Get(1).(node)
	}
	return similarNodes
}

func (nodeImplVar *nodeImpl) assertObjects(handles []reteHandle, isRight bool) {
	fmt.Println("Abstract method here.., see filterNodeImpl and joinNodeImpl")
}
