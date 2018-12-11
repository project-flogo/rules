package rete

import (
	"github.com/project-flogo/rules/common/model"
)

type joinTable interface {
	nwElemId
	getRule() model.Rule

	addRow(handles []reteHandle) joinTableRow
	removeRow(rowID int) joinTableRow
	getRow(rowID int) joinTableRow
	getRowIterator() rowIterator

	getRowCount() int
}
