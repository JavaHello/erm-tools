package core

import (
	"database/sql"
	"strconv"

	"erm-tools/src/model"

	_ "github.com/go-sql-driver/mysql"
)

// DbRead 读取数据库表结构
type DbRead struct {
	AbstractRead
	database *sql.DB
}

// NewDbRead 创建 DbRead
func NewDbRead() DbRead {
	return DbRead{AbstractRead: AbstractRead{AllTable: make(map[string]*model.Table, 16)}}
}

func (read *DbRead) db() *sql.DB {
	if read.database == nil {
		db, _ := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/information_schema?charset=utf8")
		db.SetMaxIdleConns(100)
		db.SetMaxOpenConns(10)
		read.database = db
	}

	return read.database
}

// ReadAll 读取所有数据库表结构
func (red *DbRead) ReadAll(dbname string) {
	red.readTable("tm_test2")
}

func (read *DbRead) readTable(name string) {
	db := read.db()
	stmt, _ := db.Prepare(
		`select TABLE_NAME,
		COLUMN_NAME,
		IS_NULLABLE,
		DATA_TYPE,
		CHARACTER_MAXIMUM_LENGTH,
		NUMERIC_PRECISION,
		NUMERIC_SCALE,
		COLUMN_COMMENT,
		COLUMN_TYPE,
		EXTRA,
		COLUMN_DEFAULT
	from COLUMNS
	where TABLE_SCHEMA = ?
	and TABLE_NAME = ?`)
	defer stmt.Close()
	rows, _ := stmt.Query("demodb", name)
	defer rows.Close()
	tb := read.AllTable[name]
	if tb == nil {
		tb = model.NewTable(name)
	}
	var tableName string
	var colName string
	var isNull string
	var dataType string
	var charLen string
	var numLen string
	var numScale string
	var colComment string
	var colType string
	var extra string
	var defval string
	var col model.Column
	var colMap = map[string]model.Column{}
	for rows.Next() {
		rows.Scan(&tableName, &colName, &isNull, &dataType, &charLen, &numLen, &numScale, &colComment, &colType, &extra, &defval)
		col.PhysicalName = colName
		col.LogicalName = colComment
		col.Type = dataType
		col.DefaultValue = defval
		if charLen != "" && charLen != "null" {
			l, _ := strconv.Atoi(charLen)
			col.Length = int8(l)
		} else if numLen != "" && numLen != "null" {
			l, _ := strconv.Atoi(numLen)
			col.Length = int8(l)
		}
		if numScale != "" && numScale != "null" {
			l, _ := strconv.Atoi(numScale)
			col.Decimal = int8(l)
		}
		if extra == "auto_increment" {
			col.AutoIncrement = true
		}
		col.NotNull = isNull == "NO"
		col.ColumnType = colType
		tb.Columns = append(tb.Columns, col)
		colMap[colName] = col
	}
	read.readIndex(tb, colMap)
	read.AllTable[tb.PhysicalName] = tb
}

func (red *DbRead) readIndex(table *model.Table, colMap map[string]model.Column) {
	db := red.db()
	stmt, _ := db.Prepare(
		`select TABLE_NAME, NON_UNIQUE, INDEX_NAME, COLUMN_NAME
		from STATISTICS
		where TABLE_SCHEMA = ?
		  and TABLE_NAME = ?`)
	defer stmt.Close()
	rows, _ := stmt.Query("demodb", table.PhysicalName)
	defer rows.Close()
	var tableName string
	var nonUnique string
	var indexName string
	var colName string
	var oldIndexName string
	var index model.Index
	flag := false
	for rows.Next() {
		rows.Scan(&tableName, &nonUnique, &indexName, &colName)
		if indexName == "PRIMARY" {
			col := colMap[colName]
			col.PrimaryKey = true
			table.PrimaryKeys = append(table.PrimaryKeys, col)
			continue
		}
		if oldIndexName != indexName {
			if flag {
				if index.NonUnique {
					table.Indexs = append(table.Indexs, index)
				} else {
					table.Uniques = append(table.Uniques, index)
				}
			}
			flag = true
			index = model.Index{Name: indexName}
			index.NonUnique, _ = strconv.ParseBool(nonUnique)
			index.Columns = []model.Column{}

		}
		index.Columns = append(index.Columns, colMap[colName])
		oldIndexName = indexName
	}
	if index.NonUnique {
		table.Indexs = append(table.Indexs, index)
	} else {
		table.Uniques = append(table.Uniques, index)
	}
}
