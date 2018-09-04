package model

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/core/data"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func TestOne(t *testing.T) {

	td1 := TupleDescriptor{}
	td1.Name = "a"
	td1.TTLInSeconds = 10
	td1p1 := TuplePropertyDescriptor{}
	td1p1.Name = "p1"
	td1p1.KeyIndex = 3
	td1p1.PropType = data.TypeDouble
	td1p2 := TuplePropertyDescriptor{}
	td1p2.Name = "p2"
	td1p2.KeyIndex = 31
	td1p2.PropType = data.TypeString

	td1.Props = []TuplePropertyDescriptor{td1p1, td1p2}

	str, _ := json.Marshal(&td1)
	fmt.Printf("succes %s\n", str)

	tpdx := TupleDescriptor{}
	tpdx.TTLInSeconds = -1
	json.Unmarshal([]byte(str), &tpdx)

	str1, _ := json.Marshal(&tpdx)
	fmt.Printf("succes %s\n", str1)

}

func TestTwo(t *testing.T) {
	tupleDescAbsFileNm := getAbsPathForResource("src/github.com/project-flogo/rules/rulesapp/rulesapp.json")
	tupleDescriptor := fileToString(tupleDescAbsFileNm)

	fmt.Printf("Loaded tuple descriptor: \n%s\n", tupleDescriptor)
	//First register the tuple descriptors
	RegisterTupleDescriptors(tupleDescriptor)

}

func getAbsPathForResource(resourcepath string) string {
	GOPATH := os.Getenv("GOPATH")
	regex, err := regexp.Compile(":|;")
	if err != nil {
		return ""
	}
	paths := regex.Split(GOPATH, -1)
	if os.PathListSeparator == ';' {
		//windows
		resourcepath = strings.Replace(resourcepath, "/", string(os.PathSeparator), -1)
	}
	for _, path := range paths {
		absPath := path + string(os.PathSeparator) + resourcepath
		_, err := os.Stat(absPath)
		if err == nil {
			return absPath
		}
	}
	return ""
}

func fileToString(fileName string) string {
	dat, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	return string(dat)
}
