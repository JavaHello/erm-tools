package core

import (
	"erm-tools/helper"
	"erm-tools/logger"
	"erm-tools/model"
	"os"
	"strings"
)

const (
	CREATE            = "CREATE "
	TABLE             = "TABLE "
	LF                = "\n"
	ddlOutName string = "gen_dbname.sql"
)

type OptType string

const (
	ADD    OptType = " ADD "
	DROP   OptType = " DROP "
	MODIFY OptType = " MODIFY "
)

type DdlOut struct {
	OutPath string
	outFile *os.File
}

func (out *DdlOut) Writer(diffTables []*model.DiffTable) {
	ddlFileName := strings.Replace(ddlOutName, "dbname", helper.Env.DbName, 1)
	if fp, err := os.Create(out.OutPath + string(os.PathSeparator) + ddlFileName); err != nil {
		logger.Error.Println("创建DDL文件失败", out.OutPath, err)
		return
	} else {
		out.outFile = fp
	}
	defer out.outFile.Close()
	for _, diffTab := range diffTables {
		out.genDdl(diffTab)
	}
}

func (out *DdlOut) genDdl(diffTab *model.DiffTable) {
	if diffTab.IsNew {
		var createDdl string
		createDdl += CREATE
		createDdl += TABLE
		createDdl += diffTab.Name
		createDdl += "("

		createDdl += LF
		for _, col := range diffTab.DiffColumns {
			createDdl += "\t"
			createDdl += col.NewColumn.ToDDL() + ",\n"
		}

		var idxType string
		for _, idx := range diffTab.DiffIndexes {
			if idx.NewIndex.NonUnique {
				idxType = "\tUNIQUE KEY "

			} else {
				idxType = "\tKEY "
			}
			if idx.NewIndex != nil {
				createDdl += idxType + "(" + model.IndexsColName(idx.NewIndex.Columns) + "),\n"
			}
		}

		if len(diffTab.DiffPks) > 0 {
			createDdl += "\tPRIMARY KEY ("
			var pks []string
			for _, pk := range diffTab.DiffPks {
				pks = appendColName(pk.NewColumn, pks)
			}
			createDdl += strings.Join(pks, ", ") + ",\n"
		}
		createDdl = createDdl[:len(createDdl)-2]
		out.outFile.WriteString(createDdl)
		out.outFile.WriteString(")")
		out.outFile.WriteString(diffTab.Name)
	} else {
		var colDdl string
		for _, diffCol := range diffTab.DiffColumns {
			var opt = out.optColType(diffCol.OldColumn, diffCol.NewColumn)
			if opt == MODIFY {
				colDdl += "-- " + diffCol.OldColumn.ColumnType + " 修改为 " + diffCol.NewColumn.ColumnType + "\n"
			}

			colDdl += "ALTER TABLE " + diffTab.Name + string(opt) + "COLUMN "
			if opt == ADD || opt == MODIFY {
				colDdl += diffCol.NewColumn.ToDDL() + ";\n"
			} else if opt == DROP {
				colDdl += diffCol.OldColumn.PhysicalName + ";\n"
			}
		}
		out.outFile.WriteString(colDdl)

		var idxDdl string
		var idxType string
		for _, diffIdx := range diffTab.DiffIndexes {
			var opt = out.optIdxType(diffIdx.OldIndex, diffIdx.NewIndex)
			if opt == MODIFY {
				idxDdl += diffTab.Name + " 表有索引类型修改，请自行确认\n"
			} else {

				idxDdl += "ALTER TABLE " + diffTab.Name + string(opt)
				if opt == ADD {
					if diffIdx.NewIndex.NonUnique {
						idxType = "INDEX "
					} else {
						idxType = "UNIQUE INDEX "
					}
					idxDdl += idxType + diffIdx.NewIndex.Name + " (" + model.IndexsColName(diffIdx.NewIndex.Columns) + ");\n"
				} else if opt == DROP {
					if diffIdx.OldIndex.NonUnique {
						idxType = " INDEX "
					} else {
						idxType = " UNIQUE INDEX "
					}
					idxDdl += diffIdx.Name + ";\n"
				}
			}
		}
		out.outFile.WriteString(idxDdl)
		if len(diffTab.DiffPks) > 0 {
			out.outFile.WriteString(diffTab.Name + " 表主键有修改，请自行确认~\n")
		}
	}
	out.outFile.WriteString("\n")
}

func (out *DdlOut) optColType(oldCol, newCol *model.Column) OptType {
	if oldCol == nil {
		return ADD
	} else if newCol == nil {
		return DROP
	} else {
		return MODIFY
	}
}
func (out *DdlOut) optIdxType(oldIdx, newIdx *model.Index) OptType {
	if oldIdx == nil {
		return ADD
	} else if newIdx == nil {
		return DROP
	} else {
		return MODIFY
	}
}
