package core

import (
	"erm-tools/core"
	"erm-tools/model"
	"testing"
)

func TestDiff(t *testing.T) {
	diff := core.TableDiff{}
	newTab := model.Table{
		PhysicalName: "test",
		LogicalName:  "TEST",
	}
	oldTab := model.Table{
		PhysicalName: "test",
		LogicalName:  "TEST",
	}
	oldTab.Columns = append(oldTab.Columns, &model.Column{
		PhysicalName: "c1",
		LogicalName:  "C1",
		Type:         "int",
		Length:       11,
		Decimal:      0,
	})
	diffTab := diff.Diff(&oldTab, &newTab)
	if len(diffTab.DiffColumns) == 0 {
		t.Error("差异对比错误: 没有比对出差异")
	}
	if diffTab.DiffColumns[0].OldColumn == nil {
		t.Error("差异对比错误: 差异比对错误")
	}
}
