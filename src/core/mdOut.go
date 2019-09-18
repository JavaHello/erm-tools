package core

import (
	"erm-tools/src/helper"
	"erm-tools/src/model"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	colTitle = `
|old名称|old类型|old长度|old精度||new名称|new类型|new长度|new精度|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
`
	idxTitle = `
|old名称|old字段||new名称|new字段|索引类型|
|:-:|:-:|:-:|:-:|:-:|:-:|
`
)

type MdOut struct {
	OutPath  string
	diffFile *os.File
}

func (out *MdOut) Writer(diffTables []model.DiffTable) {
	if fp, err := os.Create(out.OutPath); err != nil {
		log.Println("创建文件失败", out.OutPath, err)
		return
	} else {
		out.diffFile = fp
	}
	defer func() {
		if out.diffFile != nil {
			out.diffFile.Close()
		}
	}()
	for _, diffTab := range diffTables {
		out.colDiff(diffTab)
		out.idxDiff(diffTab)
	}
}

func (out *MdOut) colDiff(diffTable model.DiffTable) {
	out.diffFile.WriteString("# " + diffTable.Name + "\n")
	out.diffFile.WriteString(colTitle)
	for _, diffCol := range diffTable.DiffColumns {
		oldCol := diffCol.OldColumn
		newCol := diffCol.NewColumn
		var colMd string
		var oldMd string
		var newMd string
		if oldCol != nil {
			oldMd = "|" + oldCol.PhysicalName + "|" + oldCol.Type + "|" + strconv.Itoa(oldCol.Length) + "|" + strconv.Itoa(oldCol.Decimal)
		}
		if newCol != nil {
			newMd = "|" + newCol.PhysicalName + "|" + newCol.Type + "|" + strconv.Itoa(newCol.Length) + "|" + strconv.Itoa(newCol.Decimal)
		}
		if len(oldMd) == 0 {
			oldMd = "||||"
		}
		if len(newMd) == 0 {
			newMd = "||||"
		}
		colMd = oldMd + "|" + newMd + "\n"
		out.diffFile.WriteString(colMd)
	}
}

func (out *MdOut) idxDiff(diffTab model.DiffTable) {
	idxFlag := len(diffTab.DiffIndexes) > 0 || len(diffTab.DiffUniques) > 0 || len(diffTab.DiffPks) > 0
	if idxFlag {
		out.diffFile.WriteString("# " + diffTab.Name + " 索引差异\n")
		out.diffFile.WriteString(idxTitle)
	}
	if len(diffTab.DiffPks) > 0 {
		var oldPks []string
		var newPks []string
		for _, pk := range diffTab.DiffPks {
			oldPks = appendColName(pk.OldColumn, oldPks)
			newPks = appendColName(pk.NewColumn, newPks)
		}
		out.diffFile.WriteString("|PRIMARY|" + strings.Join(oldPks, ", ||") + "|PRIMARY|" + strings.Join(newPks, ", |主键|\n"))
	}
	out.indexDiff(diffTab.DiffUniques, true)
	out.indexDiff(diffTab.DiffIndexes, false)

}
func (out *MdOut) indexDiff(diffIdx []model.DiffIndex, uk bool) {
	var idxType string
	if uk {
		idxType = "唯一索引"
	} else {
		idxType = "索引"
	}
	for _, idx := range diffIdx {

		var oldName string
		var newName string
		var oldColName string
		var newColName string

		if idx.OldIndex != nil {
			oldName = idx.OldIndex.Name
			oldColName = helper.ColumnsName(idx.OldIndex.Columns)
		}
		if idx.NewIndex != nil {
			newName = idx.NewIndex.Name
			newColName = helper.ColumnsName(idx.NewIndex.Columns)
		}

		out.diffFile.WriteString("|" + oldName + "|" + oldColName + "||" + newName + "|" + newColName + "|" + idxType + "|\n")
	}
}

func appendColName(col *model.Column, nameList []string) []string {
	if col != nil {
		return append(nameList, col.PhysicalName)
	}
	return nameList
}
