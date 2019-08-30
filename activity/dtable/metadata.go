package dtable

import (
	"strings"

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
	Make []*DecisionTable `md:"make"`
}

// DecisionTable decision table rows
type DecisionTable struct {
	DtConditions []*DtCondition `md:"condition"`
	DtActions    []*DtAction    `md:"action"`
}

// DtCondition decision row condition
type DtCondition struct {
	Tuple string `md:"tuple"`
	Field string `md:"field"`
	Expr  string `md:"expr"`
}

// DtAction decision row action
type DtAction struct {
	Tuple string `md:"tuple"`
	Field string `md:"field"`
	Value string `md:"value"`
}

// FromMap fills Input struct from map
func (s *Settings) FromMap(values map[string]interface{}) (err error) {
	tasks, err := coerce.ToArray(values["make"])
	if err != nil {
		return
	}
	s.Make = make([]*DecisionTable, len(tasks))
	for i, t := range tasks {
		task := &DecisionTable{}

		condArr := make([]*DtCondition, 0)
		actArr := make([]*DtAction, 0)

		taskMap := t.(map[string]interface{})

		for k, v := range taskMap {
			if strings.Compare(k, "condition") == 0 {
				tempCondArr := v.([]interface{})
				for _, cond := range tempCondArr {
					tempMap := cond.(map[string]interface{})
					dtCondtion := &DtCondition{}
					err = metadata.MapToStruct(tempMap, dtCondtion, true)
					condArr = append(condArr, dtCondtion)
				}
			} else {
				tempActArr := v.([]interface{})
				for _, act := range tempActArr {
					tempMap := act.(map[string]interface{})
					dtAction := &DtAction{}
					err = metadata.MapToStruct(tempMap, dtAction, true)
					actArr = append(actArr, dtAction)
				}
			}
		}
		task.DtConditions = condArr
		task.DtActions = actArr
		s.Make[i] = task
	}
	return
}
