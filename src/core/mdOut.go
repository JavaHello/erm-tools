package core

import (
	"erm-tools/src/model"
	"log"
	"os"
	"strconv"
)

const (
	colTitle = `
|old名称|old类型|old长度|old精度||new名称|new类型|new长度|new精度|
|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
`
	idxTitle = `
|old名称|old字段||new名称|new字段|索引类型|
|:-:|:-:||:-:|:-:|:-:|
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
