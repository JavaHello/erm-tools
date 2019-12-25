package core

import (
	"erm-tools/helper"
	"erm-tools/model"
	"sort"
)

// IRead 读取表结构
type IRead interface {
	ReadAll(path string)
	Read(tableName string) *model.Table
}

// AbstractRead 获取单表结构的默认实现
type AbstractRead struct {
	AllTable map[string]*model.Table
}

func (red *AbstractRead) Read(name string) *model.Table {
	return red.AllTable[name]
}

// IDiff 比对两张表结构的差异
type IDiff interface {
	Diff(oldTable *model.Table, newTable *model.Table) model.DiffTable
}

// IDiffOut 差异写出到文件
type IDiffOut interface {
	Writer(diffTables []*model.DiffTable)
}

type DiffExec interface {
	Exec()
}

type ErmErmDiffExec struct{}
type ErmDbDiffExec struct{}
type DbDbDiffExec struct{}

func (e *ErmErmDiffExec) Exec() {
	newErmRead := NewErmRead()
	for _, file := range helper.Env.NewErmFiles() {
		newErmRead.ReadAll(file)
	}

	oldRead := NewErmRead()
	for _, file := range helper.Env.OldErmFiles() {
		oldRead.ReadAll(file)
	}

	diffTables := diff(&newErmRead.AbstractRead, &oldRead.AbstractRead)
	out(diffTables)
}

func (e *ErmDbDiffExec) Exec() {
	newErmRead := NewErmRead()
	for _, file := range helper.Env.NewErmFiles() {
		newErmRead.ReadAll(file)
	}
	var oldRead = NewDbRead(helper.Env.DbUser,
		helper.Env.DbPassword,
		helper.Env.DbHost,
		helper.Env.DbPort,
		helper.Env.DbName)
	oldRead.ReadAll(helper.Env.DbName)
	diffTables := diff(&newErmRead.AbstractRead, &oldRead.AbstractRead)
	out(diffTables)
}
func (e *DbDbDiffExec) Exec() {
	var newRead = NewDbRead(helper.Env.DbUser,
		helper.Env.DbPassword,
		helper.Env.DbHost,
		helper.Env.DbPort,
		helper.Env.DbName)
	dbname := helper.Env.DbName
	newRead.ReadAll(dbname)
	out := MdOut{OutPath: helper.Env.OutPath}
	ddlOut := DdlOut{OutPath: helper.Env.OutPath}
	for _, dc := range helper.Env.TargetDbList {
		var oldRead = NewDbRead(dc.DbUser,
			dc.DbPassword,
			dc.DbHost,
			dc.DbPort,
			dc.DbName)
		oldRead.ReadAll(dbname)
		diffTables := diff(&newRead.AbstractRead, &oldRead.AbstractRead)
		out.Writer(diffTables)
		if helper.Env.GenDdl {
			ddlOut.Writer(diffTables)
		}
	}

}

func diff(newRead, oldRead *AbstractRead) []*model.DiffTable {
	diff := TableDiff{}
	var diffTables []*model.DiffTable
	for _, newTab := range newRead.AllTable {
		oldTab := oldRead.Read(newTab.PhysicalName)
		if oldTab == nil {
			oldTab = &model.Table{PhysicalName: newTab.PhysicalName}
		}
		diffTab := diff.Diff(oldTab, newTab)
		diffTables = append(diffTables, &diffTab)
	}
	sort.Sort(model.DiffTableSlice(diffTables))
	return diffTables
}

func out(diffTables []*model.DiffTable) {
	sort.Sort(model.DiffTableSlice(diffTables))
	out := MdOut{OutPath: helper.Env.OutPath}
	out.Writer(diffTables)
	if helper.Env.GenDdl {
		ddlOut := DdlOut{OutPath: helper.Env.OutPath}
		ddlOut.Writer(diffTables)
	}
}

func GetExec(dt string) DiffExec {
	var exec DiffExec
	switch dt {
	case helper.ERM_ERM:
		exec = &ErmErmDiffExec{}
	case helper.ERM_DB:
		exec = &ErmDbDiffExec{}
	case helper.DB_DB:
		exec = &DbDbDiffExec{}
	}
	return exec
}
