package model

// Table 数据库表结构
type Table struct {
	PhysicalName string
	LogicalName  string
	Description  string
	Columns      []Column
	Indexs       []Index
	Uniques      []Index
}

// Column 字段信息
type Column struct {
	PhysicalName  string
	LogicalName   string
	Type          string
	AutoIncrement bool
	DefaultValue  string
	Length        int8
	Decimal       int8
	PrimaryKey    bool
	UniqueKey     bool
	NotNull       bool
	Description   string
	Desc          bool
}

// Index 索引信息
type Index struct {
	Name      string
	NonUnique bool
	Columns   []Column
}

func NewTable(name string) Table {
	var tb = Table{
		LogicalName: name,
	}
	tb.Columns = make([]Column, 8)
	tb.Indexs = make([]Index, 1)
	tb.Uniques = make([]Index, 1)
	return tb
}
