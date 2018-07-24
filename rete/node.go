package rete

import (
	"context"
	"fmt"
	"strconv"

	"github.com/TIBCOSoftware/bego/utils"
	"github.com/TIBCOSoftware/bego/common/model"
)

//node a building block of the rete network
type node interface {
	abstractNode
	getIdentifiers() []model.TupleTypeAlias
	getID() int
	addNodeLink(nodeLink)
	assertObjects(ctx context.Context, handles []reteHandle, isRight bool)
}

type nodeImpl struct {
	identifiers []model.TupleTypeAlias
	nodeLinkVar nodeLink
	id          int
}

//NewNode ... returns a new node
func newNode(identifiers []model.TupleTypeAlias) node {
	n := nodeImpl{}
	n.initNodeImpl(identifiers)
	return &n
}

func (n *nodeImpl) initNodeImpl(identifiers []model.TupleTypeAlias) {
	currentNodeID++
	n.id = currentNodeID

	n.identifiers = identifiers
}

func (n *nodeImpl) getIdentifiers() []model.TupleTypeAlias {
	return n.identifiers
}

func (n *nodeImpl) getID() int {
	return n.id
}

func (n *nodeImpl) addNodeLink(nl nodeLink) {
	n.nodeLinkVar = nl
}

func (n *nodeImpl) String() string {
	str := "id:" + strconv.Itoa(n.id) + ", idrs:"
	for _, nodeIdentifier := range n.identifiers {
		str += string(nodeIdentifier) + ","
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

func (n *nodeImpl) assertObjects(ctx context.Context, handles []reteHandle, isRight bool) {
	fmt.Println("Abstract method here.., see filterNodeImpl and joinNodeImpl")
}
