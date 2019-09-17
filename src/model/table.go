package model

// Table 数据库表结构
type Table struct {
	PhysicalName string
	LogicalName  string
	Description  string
	Columns      []Column
	PrimaryKeys  []Column
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
	ColumnType    string
}

// Index 索引信息
type Index struct {
	Name      string
	NonUnique bool
	Columns   []Column
}

func NewTable(name string) *Table {
	var tb = Table{
		PhysicalName: name,
	}
	tb.Columns = []Column{}
	tb.Indexs = []Index{}
	tb.Uniques = []Index{}
	return &tb
}
