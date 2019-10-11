package main

import (
	"erm-tools/core"
	"erm-tools/helper"
	"erm-tools/model"
)

func main() {

	helper.Env.Init()
	newErmRead := core.NewErmRead()
	for _, file := range helper.Env.NewErmFiles() {
		newErmRead.ReadAll(file)
	}

	var oldRead core.IRead
	if helper.Env.Type == helper.ERM_ERM {
		oldRead = core.NewErmRead()
	} else if helper.Env.Type == helper.ERM_DB {
		oldRead = core.NewDbRead()
	}
	for _, file := range helper.Env.OldErmFiles() {
		oldRead.ReadAll(file)
	}

	diff := core.TableDiff{}
	var diffTables []*model.DiffTable
	for _, newTab := range newErmRead.AllTable {
		diffTab := diff.Diff(oldRead.Read(newTab.PhysicalName), newTab)
		diffTables = append(diffTables, &diffTab)
	}
	out := core.MdOut{OutPath: helper.Env.OutPath}
	out.Writer(diffTables)
	if helper.Env.GenDdl {
		ddlOut := core.DdlOut{OutPath: helper.Env.OutPath}
		ddlOut.Writer(diffTables)
	}
}
