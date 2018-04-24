package rete

//All package scope variables are easy to llocate when in one place
import (
	"fmt"

	"github.com/TIBCOSoftware/bego/common/model"
)

var (
	currentNodeID int                              //id generator of sorts
	allHandles    map[model.StreamTuple]reteHandle //global handle map
)

func init() {
	fmt.Println("initializing rete")
	currentNodeID = 0
	allHandles = make(map[model.StreamTuple]reteHandle)
}
