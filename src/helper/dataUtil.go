package helper

import (
	"erm-tools/src/model"
	"sort"
	"strings"
)

func ColumnsName(cols []*model.Column) string {
	if cols == nil {
		return ""
	}
	var names []string
	for _, col := range cols {
		names = append(names, col.PhysicalName)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}
