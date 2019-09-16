package core

import (
	"encoding/xml"
	"erm-tools/src/helper"
	"erm-tools/src/model"
	"log"
)

type ErmRead struct {
	AbstractRead
}

func NewErmRead() ErmRead {
	return ErmRead{AbstractRead: AbstractRead{AllTable: make(map[string]*model.Table, 16)}}
}

func (red *ErmRead) ReadAll(path string) {
	var ermInfo model.Diagram
	err := xml.Unmarshal(helper.ReadFile(path), &ermInfo)
	if err != nil {
		log.Println("DbRead ReadAll Error", err)
	}
	helper.ErmToTable(&ermInfo, red.AllTable)
}
