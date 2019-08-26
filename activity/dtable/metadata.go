package dtable

import (
	"github.com/project-flogo/core/data/coerce"
	"github.com/project-flogo/core/data/metadata"
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
	Make []*Task `md:"make"`
}

// Task task
type Task struct {
	Tuple string `md:"tuple"`
	Field string `md:"field"`
	To    string `md:"to"`
}

// FromMap fills Input struct from map
func (s *Settings) FromMap(values map[string]interface{}) (err error) {
	tasks, err := coerce.ToArray(values["make"])
	if err != nil {
		return
	}
	s.Make = make([]*Task, len(tasks))
	for i, t := range tasks {
		task := &Task{}
		taskMap := t.(map[string]interface{})
		err = metadata.MapToStruct(taskMap, task, true)
		s.Make[i] = task
	}
	return
}
