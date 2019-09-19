package model

import (
	"sort"
	"strings"
)

// Table 数据库表结构
type Table struct {
	PhysicalName string
	LogicalName  string
	Description  string
	Columns      []*Column
	PrimaryKeys  []*Column
	Indexs       []*Index
}

// Column 字段信息
type Column struct {
	PhysicalName  string
	LogicalName   string
	Type          string
	AutoIncrement bool
	DefaultValue  string
	Length        int
	Decimal       int
	PrimaryKey    bool
	UniqueKey     bool
	NotNull       bool
	Description   string
	Desc          bool
	ColumnType    string
}

func (col *Column) ToDDL() string {
	var notNull string
	if col.NotNull {
		notNull = " NOT NULL"
	}
	var inc string
	if col.AutoIncrement {
		inc = " AUTO_INCREMENT"
	}
	var comment string
	if len(col.Description) > 0 {
		comment = " COMMENT '" + col.Description + "'"
	}
	return col.PhysicalName + " " +
		col.ColumnType +
		notNull + inc + comment
}

// Index 索引信息
type Index struct {
	Name      string
	NonUnique bool
	Columns   []*Column
}

func (idx *Index) ToCreateDDL() string {
	var idxType string
	if idx.NonUnique {
		idxType = "KEY "
	} else {
		idxType = "UNIQUE KEY "
	}
	return idxType + idx.Name + IndexsColName(idx.Columns)
}

func NewTable(name string) *Table {
	var tb = Table{
		PhysicalName: name,
	}
	tb.Columns = []*Column{}
	tb.Indexs = []*Index{}
	return &tb
}

func IndexsColName(cols []*Column) string {
	if cols == nil {
		return ""
	}
	var names []string
	for _, col := range cols {
		var idxColName = col.PhysicalName
		if col.Desc {
			idxColName += " DESC"
		}
		names = append(names, idxColName)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

func ColumnsName(cols []*Column) string {
	if cols == nil {
		return ""
	}
	var names []string
	for _, col := range cols {
		names = append(names, col.PhysicalName)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}
