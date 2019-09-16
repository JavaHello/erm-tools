package helper

import (
	"log"
	"strconv"

	"erm-tools/src/model"
)

// ErmToTable erm 转 table
func ErmToTable(erm *model.Diagram, tableMap map[string]*model.Table) {
	if erm == nil || tableMap == nil {
		return
	}
	var wordMap = make(map[string]*model.Word, 16)
	for _, w := range erm.Dictionary.Words {
		wordMap[w.Id] = &w
	}

	for _, t := range erm.Contents.Table {
		var tb = model.Table{PhysicalName: t.PhysicalName,
			LogicalName: t.LogicalName,
			Description: t.Description}
		tb.Columns = make([]model.Column, 8)
		tb.Indexs = make([]model.Index, 1)
		tb.Uniques = make([]model.Index, 1)
		var mapCols = make(map[string]model.Column, 16)
		for _, ermCol := range t.Columns.NormalColumn {
			mapCol := wordMap[ermCol.WordId]
			if mapCol == nil {
				log.Println(t.PhysicalName + "表缺失字段")
				continue
			}
			col := model.Column{
				PhysicalName: mapCol.PhysicalName,
				LogicalName:  mapCol.LogicalName,
				Description:  mapCol.Description,
				Type:         mapCol.Type}
			col.NotNull, _ = strconv.ParseBool(ermCol.NotNull)
			col.AutoIncrement, _ = strconv.ParseBool(ermCol.AutoIncrement)
			col.PrimaryKey, _ = strconv.ParseBool(ermCol.PrimaryKey)
			col.UniqueKey, _ = strconv.ParseBool(ermCol.UniqueKey)
			if mapCol.Length != "" {
				l, _ := strconv.Atoi(mapCol.Length)
				col.Length = int8(l)
			}
			if mapCol.Decimal != "" {
				l, _ := strconv.Atoi(mapCol.Decimal)
				col.Decimal = int8(l)
			}
			if col.UniqueKey {
				createColUniqueKey(&col, &tb)
			}

			tb.Columns = append(tb.Columns, col)
			mapCols[ermCol.Id] = col

		}

		for _, idxes := range t.Indexes.Inidex {
			nonunique, _ := strconv.ParseBool(idxes.NonUnique)
			tbIdx := model.Index{Name: idxes.Name,
				NonUnique: nonunique}
			tbIdx.Columns = make([]model.Column, 1)
			for _, idxCol := range idxes.Columns.Column {
				tbCol := mapCols[idxCol.Id]
				tbCol.Desc, _ = strconv.ParseBool(idxCol.Desc)
				tbIdx.Columns = append(tbIdx.Columns, tbCol)
			}
			if nonunique {
				tb.Indexs = append(tb.Indexs, tbIdx)
			} else {
				tb.Uniques = append(tb.Uniques, tbIdx)
			}
		}
		tableMap[t.PhysicalName] = &tb
	}

}

func createColUniqueKey(column *model.Column, table *model.Table) {
	tbIdx := model.Index{Name: "idxes",
		NonUnique: true}
	tbIdx.Columns = append(tbIdx.Columns, *column)
	table.Uniques = append(table.Uniques, tbIdx)
}
