package test

import (
	"encoding/json"
	"testing"
	"bytes"

    "github.com/project-flogo/rules/common/model"
)

func TestTupleObjectValue(t *testing.T) {

	if err := model.RegisterTupleDescriptors(`
        [
          {
            "name": "td1",
            "properties": [
              {
                "name": "key",
                "pk-index": 0,
                "type": "string"
              },
              {
                "name": "thing",
                "type": "object"
              }
            ]
		}]
	`); err != nil {
		t.Errorf("Failed to register tuple descriptor: %s", err)
	}

	tuple, err := model.NewTupleWithKeyValues("td1", "Bob")
	if err != nil {
		t.Errorf("Failed to create tuple: %s", err)
	}

	var thing map[string]interface{}
	if err = json.Unmarshal([]byte(`{"foo": {"bar":"baz"}}`), &thing); err != nil {
		t.Errorf("Failed to unmarshal object: %s", err);
	}

	if err = tuple.SetObject(nil, "thing", thing); err != nil {
		t.Errorf("Failed to SetObject: %s", err)
	}

	thing2, err := tuple.GetObject("thing")
	if err != nil {
		t.Errorf("Failed to GetObject: %s", err)
	}

	thingJson, err := json.Marshal(thing)
	thing2Json, err := json.Marshal(thing2)

	if !bytes.Equal(thingJson, thing2Json) {
		t.Errorf("Thing and thing2 are not equal: [%s] [%s]", thingJson, thing2Json);
	}
}

func TestTupleObjectValueAsKey(t *testing.T) {

	expect := "Property [key] is a mutable object and cannot be defined as a key for type [td1]"
	err := model.RegisterTupleDescriptors(`
        [
          {
            "name": "td1",
            "properties": [
              {
                "name": "key",
                "pk-index": 0,
                "type": "object"
              }
            ]
		}]
	`)
	if err == nil {
		t.Errorf("Object was allowed as tuple key")
	} else
	if err.Error() != expect {
		t.Errorf("Unexpected error registering td, wanted \"%s\" got \"%s\"", expect, err.Error())
	}
}
