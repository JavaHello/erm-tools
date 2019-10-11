package core

import (
	"erm-tools/logger"
	"erm-tools/model"
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
|old名称|old字段|old索引类型||new名称|new字段|new索引类型|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
`
	mdOutName = "diff.md"
)

type MdOut struct {
	OutPath  string
	diffFile *os.File
}

func (out *MdOut) Writer(diffTables []*model.DiffTable) {
	if fp, err := os.Create(out.OutPath + string(os.PathSeparator) + mdOutName); err != nil {
		logger.Error("创建DIFF文件失败", out.OutPath, err)
		return
	} else {
		out.diffFile = fp
	}
	defer out.diffFile.Close()
	for _, diffTab := range diffTables {
		out.colDiff(diffTab)
		out.idxDiff(diffTab)
	}
}

func (out *MdOut) colDiff(diffTable *model.DiffTable) {
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

func (out *MdOut) idxDiff(diffTab *model.DiffTable) {
	idxFlag := len(diffTab.DiffIndexes) > 0 || len(diffTab.DiffPks) > 0
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
		out.diffFile.WriteString("|PRIMARY|" + strings.Join(oldPks, ", ") + "|主键||PRIMARY|" + strings.Join(newPks, ", |主键|\n"))
	}
	out.indexDiff(diffTab.DiffIndexes)

}
func (out *MdOut) indexDiff(diffIdx []model.DiffIndex) {

	for _, idx := range diffIdx {

		var oldIdxType string
		var newIdxType string
		var oldName string
		var newName string
		var oldColName string
		var newColName string

		if idx.OldIndex != nil {
			oldName = idx.OldIndex.Name
			oldColName = model.ColumnsName(idx.OldIndex.Columns)
			if !idx.OldIndex.NonUnique {
				oldIdxType = "唯一索引"
			} else {
				oldIdxType = "索引"
			}
		}
		if idx.NewIndex != nil {
			if !idx.NewIndex.NonUnique {
				newIdxType = "唯一索引"
			} else {
				newIdxType = "索引"
			}
			newName = idx.NewIndex.Name
			newColName = model.ColumnsName(idx.NewIndex.Columns)
		}

		out.diffFile.WriteString("|" + oldName + "|" + oldColName + "|" + oldIdxType + "||" + newName + "|" + newColName + "|" + newIdxType + "|\n")
	}
}

func appendColName(col *model.Column, nameList []string) []string {
	if col != nil {
		return append(nameList, col.PhysicalName)
	}
	return nameList
}
