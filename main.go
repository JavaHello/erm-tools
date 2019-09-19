package main

import (
	"fmt"

	"erm-tools/core"
	"erm-tools/model"
)

func main() {
	read := core.NewErmRead()
	read.ReadAll(`D:\workspace\JavaProjects\demo\src\main\resources\db.erm`)

	ermRead := core.NewErmRead()
	ermRead.ReadAll(`D:\workspace\JavaProjects\demo\src\main\resources\db2.erm`)

	diff := core.TableDiff{}
	var diffTables []model.DiffTable
	for _, newTab := range read.AllTable {
		diffTab := diff.Diff(ermRead.Read(newTab.PhysicalName), newTab)
		fmt.Println(diffTab)
		diffTables = append(diffTables, diffTab)
	}
	out := core.MdOut{OutPath: "D:\\workspace\\GoProjects\\erm-tools\\diff.md"}
	out.Writer(diffTables)
	ddlOut := core.DdlOut{OutPath: "D:\\workspace\\GoProjects\\erm-tools\\gen.sql"}
	ddlOut.Writer(diffTables)
}
