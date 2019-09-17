package core

import "erm-tools/src/model"

type TableDiff struct {
}

// Diff 比对表结构差异
func (diff *TableDiff) Diff(oldTable *model.Table, newTable *model.Table) model.DiffTable {

	// 字段比较
	oldCols := oldTable.Columns
	newCols := newTable.Columns

	oldGroupCols := groupColumns(oldCols)

	diffTab := model.NewDiffTable(newTable.PhysicalName)
	for _, newCol := range newCols {
		var diffCol model.DiffColumn
		diffCol.Name = newCol.PhysicalName
		oldCol, ok := oldGroupCols[newCol.PhysicalName]
		delete(oldGroupCols, newCol.PhysicalName)

		if !ok {
			diffCol.NewColumn = &newCol
			diffTab.DiffColumns = append(diffTab.DiffColumns, diffCol)
			continue
		}
		if newCol.Type != oldCol.Type || newCol.Length != oldCol.Length || newCol.Decimal != oldCol.Decimal {
			diffCol.NewColumn = &newCol
			diffCol.OldColumn = &oldCol
			diffTab.DiffColumns = append(diffTab.DiffColumns, diffCol)
			continue
		}
	}
	for _, oldCol := range oldGroupCols {
		var diffCol model.DiffColumn
		diffCol.Name = oldCol.PhysicalName
		diffCol.OldColumn = &oldCol
		diffTab.DiffColumns = append(diffTab.DiffColumns, diffCol)
	}

	// 主键索引比较
	newPks := newTable.PrimaryKeys
	oldPks := oldTable.PrimaryKeys

	oldGroupPks := groupColumns(oldPks)

	for _, newPk := range newPks {
		var diffCol model.DiffColumn
		diffCol.Name = newPk.PhysicalName
		_, ok := oldGroupPks[newPk.PhysicalName]
		delete(oldGroupPks, newPk.PhysicalName)
		if !ok {
			diffCol.NewColumn = &newPk
			diffTab.DiffColumns = append(diffTab.DiffPks, diffCol)
		}
	}
	for _, oldPk := range oldGroupPks {
		var diffCol model.DiffColumn
		diffCol.Name = oldPk.PhysicalName
		diffCol.OldColumn = &oldPk
		diffTab.DiffColumns = append(diffTab.DiffPks, diffCol)
	}

	return diffTab
}

func groupColumns(cols []model.Column) map[string]model.Column {
	var res = map[string]model.Column{}
	for _, col := range cols {
		res[col.PhysicalName] = col
	}
	return res
}
