package main

import (
	"erm-tools/core"
	"erm-tools/helper"
)

func main() {
	helper.Env.Init()
	core.GetExec(helper.Env.Type).Exec()
}
