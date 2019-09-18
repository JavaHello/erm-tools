package core

import "erm-tools/src/model"

// IRead 读取表结构
type IRead interface {
	ReadAll(path string)
	Read(tableName string) model.Table
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
	Writer(diffTables []model.DiffTable)
}
