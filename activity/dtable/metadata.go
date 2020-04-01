package dtable

import (
	"github.com/project-flogo/core/data/coerce"
)

// Input activity input
type Input struct {
	TBDMessage string `md:"message"`
}

// ToMap converts Input struct to map
func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"message": i.TBDMessage,
	}
}

// FromMap fills Input struct from map
func (i *Input) FromMap(values map[string]interface{}) (err error) {
	i.TBDMessage, err = coerce.ToString(values["message"])
	return
}

// Settings activity settings
type Settings struct {
	DTableFile string `md:"dTableFile,required"` // Path to decision table file
}
