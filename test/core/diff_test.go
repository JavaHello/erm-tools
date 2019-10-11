package core

import (
	"erm-tools/core"
	"erm-tools/model"
	"github.com/stretchr/testify/assert"
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
	assert.Len(t, diffTab.DiffColumns, 1, "差异对比错误: 没有比对出差异")
	assert.NotNil(t, diffTab.DiffColumns[0].OldColumn, "差异对比错误: 差异比对错误")
}
