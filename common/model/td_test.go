package model

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
)

func TestOne(t *testing.T) {

	dat, err := ioutil.ReadFile("/home/bala/go/src/github.com/TIBCOSoftware/bego/common/model/tupledescriptor.json")
	if err != nil {
		log.Fatal(err)
	}

	tds := []TupleDescriptor{}
	json.Unmarshal([]byte(string(dat)), &tds)

	str, _ok := json.Marshal(&tds)
	if _ok == nil {
		fmt.Printf("succes %s\n", str)
	} else {
		fmt.Printf("err:  %s\n", "error")
	}

	for _, td := range tds {
		fmt.Printf("Type [%s], KeyProps [%s]\n", td.Name, td.GetKeyProps())
	}

}
