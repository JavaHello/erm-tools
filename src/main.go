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

	dbRead := core.NewDbRead()
	dbRead.ReadAll("demodb")
}
