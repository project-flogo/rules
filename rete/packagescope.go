package rete

//All package scope variables are easy to llocate when in one place

var (
	currentNodeID int //id generator of sorts
	// allHandles    map[model.StreamTuple]reteHandle //global handle map
)

func init() {
	currentNodeID = 0
}
