package model

import (
	"testing"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

func TestOne (t *testing.T) {


	dat, err := ioutil.ReadFile("/home/bala/go/src/github.com/TIBCOSoftware/bego/common/model/tupledescriptor.json")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Tuple descriptors: [%s]\n", string(dat))

	tds := []TupleDescriptor{}
	json.Unmarshal([]byte(string(dat)),&tds)
	fmt.Printf("xx")

	str, _ok:= json.Marshal(&tds)
	if _ok == nil  {
		fmt.Printf ("succes %s\n", str)
	} else {
		fmt.Printf ("err:  %s\n", "error")
	}
	//td := TupleDescriptor{"n1", 10, make(map[string]TuplePropertyDescriptor)}
	//
	//p1 := TuplePropertyDescriptor{"p1", "t1"}
	//p2 := TuplePropertyDescriptor{"p2", "t2"}
	//
	//
	//td.Props["p1"] = p1
	//td.Props["p2"] = p2
	//
	//
	//
	//str, _ok:= json.Marshal(&td)
	//

	//
	//
	//jsonized1 := "{\"Alias\":\"n1\",\"Timeout\":10,\"Props\":{\"p1\":{\"Name\":\"p1\",\"PropType\":\"t1\"},\"p2\":{\"Name\":\"p2\",\"PropType\":\"t2\"}}}"
	//td1 := TupleDescriptor{}
	//json.Unmarshal([]byte(jsonized1),&td1)
	//jsonized2 := "{\"Alias\":\"n2\",\"Timeout\":10,\"Props\":[\"p1\":{\"Name\":\"p1\",\"PropType\":\"t1\"}]}"
	//
	//json3 := "["+jsonized1 + "," + jsonized2+"]"
	//fmt.Printf ("%s\n", json3)
	//
	// tds :=make([]TupleDescriptor, 0)
	//
	//json.Unmarshal([]byte(json3),&tds)


	//td1 := TupleDescriptor{}

	//json.Unmarshal([]byte(jsonized),&td1)
	//
	//
	//str2, _:= json.Marshal(&td1)
	//fmt.Printf ("again... %s\n", str2)
	//
	//
	//tds := []TupleDescriptor {td, td1}
	//
	//str3 , _ := json.Marshal(tds)
	//fmt.Printf ("ccccc%s\n", str3)
}
