package dtable

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/project-flogo/core/activity"
	"github.com/project-flogo/core/data"
	"github.com/project-flogo/core/data/metadata"
	"github.com/project-flogo/rules/common/model"
	"github.com/project-flogo/rules/ruleapi"
)

// Decision table column types
type dtColType int8

const (
	ctID          dtColType = iota // ID
	ctCondition                    // condition
	ctAction                       // action
	ctDescription                  // Description
	ctPriority                     // priority
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity decision table based rule action
type Activity struct {
	dtable *dTable
}

type dTable struct {
	titleRow1 []genCell
	titleRow2 []genCell
	metaRow   []metaCell
	rows      [][]*genCell
}

type metaCell struct {
	colType   dtColType
	tupleDesc *model.TupleDescriptor
	propDesc  *model.TuplePropertyDescriptor
}

type genCell struct {
	*metaCell
	rawValue string
	cdExpr   string
}

// New creates new decision table activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	// Read settings
	settings := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), settings, true)
	if err != nil {
		return nil, err
	}

	// Read decision table from file
	dtable, err := loadFromFile(settings.DTableFile)
	if err != nil {
		return nil, err
	}
	// dtable.print()
	err = dtable.compile()
	if err != nil {
		return nil, err
	}
	// dtable.print()

	// Read setting from init context
	act := &Activity{
		dtable: dtable,
	}
	return act, nil
}

// Metadata activity metadata
func (a *Activity) Metadata() *activity.Metadata {
	return activity.ToMetadata(&Input{})
}

// Eval implements decision table action
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {

	context := ctx.GetInput("ctx").(context.Context)
	tuples := ctx.GetInput("tuples").(map[model.TupleType]model.Tuple)

	// evaluate decision table
	a.dtable.apply(context, tuples)

	return true, nil
}

func loadFromFile(fileName string) (*dTable, error) {
	if fileName == "" {
		return nil, fmt.Errorf("file name can't be empty")
	}
	tokens := strings.Split(fileName, ".")
	fileExtension := tokens[len(tokens)-1]

	if fileExtension == "csv" || fileExtension == "CSV" {
		return loadFromCSVFile(fileName)
	} else if fileExtension == "xls" || fileExtension == "XLS" || fileExtension == "xlsx" || fileExtension == "XLSX" {
		return loadFromXLSFile(fileName)
	}

	return nil, fmt.Errorf("file[%s] extension not supported", fileName)
}

// loadFromCSVFile loads decision table from CSV file
func loadFromCSVFile(fileName string) (*dTable, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("not able open the file [%s] - %s", fileName, err)
	}
	defer file.Close()

	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("not able read the file [%s] - %s", fileName, err)
	}

	dtable := &dTable{}
	dtable.rows = make([][]*genCell, len(lines)-2)
	for i, line := range lines {
		if i == 0 {
			// title row 1
			dtable.titleRow1 = make([]genCell, len(line))
			for j, val := range line {
				dtable.titleRow1[j].rawValue = val
			}
			continue
		}
		if i == 1 {
			// title row 2
			dtable.titleRow2 = make([]genCell, len(line))
			for j, val := range line {
				dtable.titleRow2[j].rawValue = val
			}
			continue
		}
		// other rows
		row := make([]*genCell, len(line))
		for j, val := range line {
			row[j] = &genCell{
				rawValue: val,
			}
		}
		dtable.rows[i-2] = row
	}
	return dtable, nil
}

// loadFromXLSFile loads decision table from Excel file
func loadFromXLSFile(fileName string) (*dTable, error) {
	file, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("not able open the file [%s] - %s", fileName, err)
	}
	rows, err := file.GetRows("DecisionTable")
	if err != nil {
		return nil, fmt.Errorf("DecisionTable worksheet not available in %s", fileName)
	}
	// find titleRowIndex
	titleRowIndex := 0
	for i, r := range rows {
		if len(r) > 0 {
			if r[0] == "DecisionTable" {
				titleRowIndex = i + 2
			}
		}
	}
	titleRowSize := len(rows[titleRowIndex])

	dtable := &dTable{
		titleRow1: make([]genCell, titleRowSize),
		titleRow2: make([]genCell, titleRowSize),
		rows:      make([][]*genCell, 1),
	}
	// title row 1
	for i, val := range rows[titleRowIndex] {
		dtable.titleRow1[i].rawValue = val
	}
	// title row 2
	for i, val := range rows[titleRowIndex+1] {
		dtable.titleRow2[i].rawValue = val
	}
	// other rows
	for _, r := range rows[titleRowIndex+2:] {
		if len(r) == 0 {
			break
		}
		dtrow := make([]*genCell, titleRowSize)
		for i, cell := range r {
			dtrow[i] = &genCell{
				rawValue: cell,
			}
		}
		dtable.rows = append(dtable.rows, dtrow)
	}

	return dtable, nil
}

func (dtable *dTable) compile() error {
	// compute meta row from titleRow1 & titleRow2
	metaRow := make([]metaCell, len(dtable.titleRow1))
	dtable.metaRow = metaRow
	// titleRow1 determines column type
	for colIndex, cell := range dtable.titleRow1 {
		if strings.Contains(cell.rawValue, "Id") {
			metaRow[colIndex].colType = ctID
		} else if strings.Contains(cell.rawValue, "Condition") {
			metaRow[colIndex].colType = ctCondition
		} else if strings.Contains(cell.rawValue, "Action") {
			metaRow[colIndex].colType = ctAction
		} else if strings.Contains(cell.rawValue, "Description") {
			metaRow[colIndex].colType = ctDescription
		} else if strings.Contains(cell.rawValue, "Priority") {
			metaRow[colIndex].colType = ctPriority
		} else {
			return fmt.Errorf("unknown column type - %s", cell.rawValue)
		}
	}
	// titleRow2 determines tuple type & property
	for colIndex, cell := range dtable.titleRow2 {
		if cell.rawValue == "" {
			continue
		}
		tokens := strings.Split(cell.rawValue, ".")
		if len(tokens) != 2 {
			return fmt.Errorf("[%s] is not a valid tuple property representation", cell.rawValue)
		}
		tupleType := tokens[0]
		propName := tokens[1]
		tupleDesc := model.GetTupleDescriptor(model.TupleType(tupleType))
		if tupleDesc == nil {
			return fmt.Errorf("tuple type[%s] is not recognised", tupleType)
		}
		propDesc := tupleDesc.GetProperty(propName)
		if propDesc == nil {
			return fmt.Errorf("property[%s] is not a valid property for the tuple type[%s]", propName, tupleType)
		}
		metaRow[colIndex].tupleDesc = tupleDesc
		metaRow[colIndex].propDesc = propDesc
	}
	// process all rows
	for _, row := range dtable.rows {
		for colIndex, cell := range row {
			if cell == nil {
				continue
			}
			cell.metaCell = &metaRow[colIndex]
			if cell.colType == ctCondition {
				value := cell.rawValue
				if !strings.HasPrefix(value, "==") && !strings.HasPrefix(value, ">") && !strings.HasPrefix(value, "<") && !strings.HasPrefix(value, "!") {
					value = "== " + value
				}
				cell.cdExpr = "$." + cell.tupleDesc.Name + "." + cell.propDesc.Name + " " + value
			}
		}
	}
	return nil
}

func (dtable *dTable) apply(ctx context.Context, tuples map[model.TupleType]model.Tuple) {
	// process all rows
	for _, row := range dtable.rows {
		// process row conditions
		rowTruthiness := true
		for _, cell := range row {
			if cell == nil {
				continue
			}
			if cell.colType == ctCondition {
				cellTruthiness := evaluateExpression(cell.cdExpr, tuples)
				rowTruthiness = rowTruthiness && cellTruthiness
				if !rowTruthiness {
					break
				}
			}
		}
		// process row actions if all row conditions are evaluated to true
		if rowTruthiness {
			for _, cell := range row {
				if cell == nil {
					continue
				}
				if cell.colType == ctAction {
					updateTuple(ctx, tuples, cell.tupleDesc.Name, cell.propDesc.Name, cell.rawValue)
				}
			}
		}
	}
}

// print prints decision table into stdout - TO BE REMOVED
func (dtable *dTable) print() {
	// title
	for _, v := range dtable.titleRow1 {
		fmt.Printf("|  %v  |", v.rawValue)
	}
	fmt.Println()
	// meta title
	for _, v := range dtable.titleRow2 {
		fmt.Printf("|  %v  |", v.rawValue)
	}
	fmt.Println()
	// data
	for _, row := range dtable.rows {
		for _, rv := range row {
			// fmt.Printf("|  %v--%v  |", rv.cdExpr, rv.metaCell)
			fmt.Print(rv)
		}
		fmt.Println()
	}
}

// evaluateExpression evaluates expr into boolean value in tuples scope
func evaluateExpression(expr string, tuples map[model.TupleType]model.Tuple) bool {
	condExpr := ruleapi.NewExprCondition(expr)
	result, err := condExpr.Evaluate("", "", tuples, "")
	if err != nil {
		return false
	}
	return result
}

// updateTuple updates tuple's prop with toValue
func updateTuple(context context.Context, tuples map[model.TupleType]model.Tuple, tupleType string, prop string, toVlaue interface{}) {
	tuple := tuples[model.TupleType(tupleType)]
	if tuple == nil {
		return
	}
	mutableTuple := tuple.(model.MutableTuple)
	tds := mutableTuple.GetTupleDescriptor()
	strVal := fmt.Sprintf("%v", toVlaue)
	switch tds.GetProperty(prop).PropType {
	case data.TypeString:
		if strings.Compare(strVal, "<nil>") == 0 {
			strVal = ""
		}
		mutableTuple.SetString(context, prop, strVal)
	case data.TypeBool:
		if strings.Compare(strVal, "<nil>") == 0 {
			strVal = "false"
		}
		b, err := strconv.ParseBool(strVal)
		if err == nil {
			mutableTuple.SetBool(context, prop, b)
		}
	case data.TypeInt:
		if strings.Compare(strVal, "<nil>") == 0 {
			strVal = "0"
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err == nil {
			mutableTuple.SetInt(context, prop, int(i))
		}
	case data.TypeInt32:
		if strings.Compare(strVal, "<nil>") == 0 {
			strVal = "0"
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err == nil {
			mutableTuple.SetInt(context, prop, int(i))
		}
	case data.TypeInt64:
		if strings.Compare(strVal, "<nil>") == 0 {
			strVal = "0"
		}
		i, err := strconv.ParseInt(strVal, 10, 64)
		if err == nil {
			mutableTuple.SetLong(context, prop, i)
		}
	case data.TypeFloat32:
		if strings.Compare(strVal, "<nil>") == 0 {
			strVal = "0.0"
		}
		f, err := strconv.ParseFloat(strVal, 32)
		if err == nil {
			mutableTuple.SetDouble(context, prop, f)
		}
	case data.TypeFloat64:
		if strings.Compare(strVal, "<nil>") == 0 {
			strVal = "0.0"
		}
		f, err := strconv.ParseFloat(strVal, 64)
		if err == nil {
			mutableTuple.SetDouble(context, prop, f)
		}
	default:
		mutableTuple.SetValue(context, prop, toVlaue)

	}
}
