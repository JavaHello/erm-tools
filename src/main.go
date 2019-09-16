package main

import (
	"fmt"

	"erm-tools/src/core"
)

func main() {
	read := core.NewErmRead()
	read.ReadAll(`D:\workspace\JavaProjects\demo\src\main\resources\db.erm`)
	tb := read.Read("tm_test")
	fmt.Println(tb)

	fmt.Println("-------------------------------------------------")
	dbRead := core.NewDbRead()
	dbRead.ReadAll("demodb")
	tb = dbRead.Read("tm_test2")
	fmt.Println(tb)
}
