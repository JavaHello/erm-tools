package core

import "erm-tools/src/model"

type IRead interface {
	ReadAll(path string)
	Read(tableName string) model.Table
}

type AbstractRead struct {
	AllTable map[string]*model.Table
}

func (red *AbstractRead) Read(name string) model.Table {
	return *red.AllTable[name]
}
