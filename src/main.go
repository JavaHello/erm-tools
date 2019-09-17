package main

import (
	"fmt"

	"erm-tools/src/core"
)

func main() {
	read := core.NewErmRead()
	read.ReadAll(`D:\workspace\JavaProjects\demo\src\main\resources\db.erm`)

	ermRead := core.NewErmRead()
	ermRead.ReadAll(`D:\workspace\JavaProjects\demo\src\main\resources\db2.erm`)

	diff := core.TableDiff{}
	for _, newTab := range read.AllTable {
		diffTab := diff.Diff(ermRead.Read(newTab.PhysicalName), newTab)
		fmt.Println(diffTab)
	}
}
