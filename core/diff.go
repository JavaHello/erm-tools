package core

import (
	"erm-tools/model"
	"strings"
)

var (
	// IgLenCmp 忽略长度比对类型
	IgLenCmp = []string{"date", "datetime"}
)

// TableDiff 差异比较
type TableDiff struct {
}

func (diff *TableDiff) diffTable(diffTable *model.DiffTable, physicalName, logicalName string) *model.DiffTable {
	if diffTable == nil {
		diffTable = model.NewDiffTable(physicalName, logicalName)
	}
	return diffTable
}

// Diff 比对表结构差异
func (diff *TableDiff) Diff(oldTable *model.Table, newTable *model.Table) *model.DiffTable {

	// 字段比较
	oldCols := oldTable.Columns
	newCols := newTable.Columns

	oldGroupCols := groupColumns(oldCols)

	var diffTab *model.DiffTable
	for _, newCol := range newCols {
		var diffCol model.DiffColumn
		diffCol.Name = newCol.PhysicalName
		oldCol, ok := oldGroupCols[newCol.PhysicalName]
		delete(oldGroupCols, newCol.PhysicalName)

		if !ok {
			diffCol.NewColumn = newCol
			if diffTab == nil {
				diffTab = diff.diffTable(diffTab, newTable.PhysicalName, newTable.LogicalName)
			}
			diffTab.DiffColumns = append(diffTab.DiffColumns, diffCol)
			continue
		}
		if !same(newCol.Type, oldCol.Type) || (newCol.Length != oldCol.Length || newCol.Decimal != oldCol.Decimal) && !igLenCmp(newCol.Type) {
			diffCol.NewColumn = newCol
			diffCol.OldColumn = oldCol
			if diffTab == nil {
				diffTab = diff.diffTable(diffTab, newTable.PhysicalName, newTable.LogicalName)
			}
			diffTab.DiffColumns = append(diffTab.DiffColumns, diffCol)
			continue
		}
	}
	for _, oldCol := range oldGroupCols {
		var diffCol model.DiffColumn
		diffCol.Name = oldCol.PhysicalName
		diffCol.OldColumn = oldCol
		if diffTab == nil {
			diffTab = diff.diffTable(diffTab, newTable.PhysicalName, newTable.LogicalName)
		}
		diffTab.DiffColumns = append(diffTab.DiffColumns, diffCol)
	}

	// 主键比较
	newPks := newTable.PrimaryKeys
	oldPks := oldTable.PrimaryKeys

	oldGroupPks := groupColumns(oldPks)

	for _, newPk := range newPks {
		var diffPk model.DiffColumn
		diffPk.Name = newPk.PhysicalName
		_, ok := oldGroupPks[newPk.PhysicalName]
		delete(oldGroupPks, newPk.PhysicalName)
		if !ok {
			diffPk.NewColumn = newPk
			if diffTab == nil {
				diffTab = diff.diffTable(diffTab, newTable.PhysicalName, newTable.LogicalName)
			}
			diffTab.DiffPks = append(diffTab.DiffPks, diffPk)
		}
	}
	for _, oldPk := range oldGroupPks {
		var diffPk model.DiffColumn
		diffPk.Name = oldPk.PhysicalName
		diffPk.OldColumn = oldPk
		if diffTab == nil {
			diffTab = diff.diffTable(diffTab, newTable.PhysicalName, newTable.LogicalName)
		}
		diffTab.DiffPks = append(diffTab.DiffPks, diffPk)
	}

	// 索引比较
	diffIndexes(newTable.Indices, oldTable.Indices, diffTab)

	if diffTab != nil {
		diffTab.IsNew = len(oldTable.Columns) == 0
	}
	return diffTab
}

func same(t1, t2 string) bool {
	if t1 == t2 {
		return true
	}
	if isInt(t1) && isInt(t2) {
		return true
	}
	return false
}
func isInt(t string) bool {
	t = strings.ToLower(t)
	return "int" == t || "integer" == t
}

func igLenCmp(ct string) bool {
	for _, t := range IgLenCmp {
		if strings.ToLower(t) == strings.ToLower(ct) {
			return true
		}
	}
	return false
}

// 索引比较
func diffIndexes(newIdxes, oldIdxes []*model.Index, diffTab *model.DiffTable) {
	oldGroupIdxes := groupIndex(oldIdxes)
	for _, newIdx := range newIdxes {
		var diffIdx model.DiffIndex
		diffIdx.Name = newIdx.Name
		newKey := model.IndexKeyName(newIdx)
		_, ok := oldGroupIdxes[newKey]
		delete(oldGroupIdxes, newKey)
		if !ok {
			diffIdx.NewIndex = newIdx
			diffTab.DiffIndexes = append(diffTab.DiffIndexes, diffIdx)
		}
	}
	for _, oldIdx := range oldGroupIdxes {
		var diffIdx model.DiffIndex
		diffIdx.Name = oldIdx.Name
		diffIdx.OldIndex = oldIdx
		diffTab.DiffIndexes = append(diffTab.DiffIndexes, diffIdx)
	}
}

func groupIndex(idxes []*model.Index) (res map[string]*model.Index) {
	res = map[string]*model.Index{}
	for _, idx := range idxes {
		res[model.IndexKeyName(idx)] = idx
	}
	return res
}

func groupColumns(cols []*model.Column) (res map[string]*model.Column) {
	res = map[string]*model.Column{}
	for _, col := range cols {
		res[col.PhysicalName] = col
	}
	return res
}
