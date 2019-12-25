package core

import (
	"database/sql"
	"erm-tools/logger"
	"fmt"
	"strconv"

	"erm-tools/model"

	_ "github.com/go-sql-driver/mysql"
)

// DbRead 读取数据库表结构
type DbRead struct {
	AbstractRead
	database *sql.DB
	User     string
	Pass     string
	Host     string
	Port     string
	DbName   string
}

// NewDbRead 创建 DbRead
func NewDbRead(user, pass, host, port, dbName string) *DbRead {
	return &DbRead{AbstractRead: AbstractRead{AllTable: make(map[string]*model.Table, 16)},
		User:   user,
		Pass:   pass,
		Host:   host,
		Port:   port,
		DbName: dbName,
	}
}

func (dr *DbRead) url() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/information_schema?charset=utf8", dr.User, dr.Pass, dr.Host, dr.Port)
}

func (read *DbRead) db() *sql.DB {
	if read.database == nil {
		url := read.url()
		db, err := sql.Open("mysql", url)
		if err != nil {
			logger.Error("数据库连接失败,URL:", url)
			panic(err)
		}
		_, err = db.Query("select 1")
		if err != nil {
			logger.Error("数据库连接失败,URL:", url)
			panic(err)
		}
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(10)
		read.database = db
	}

	return read.database
}

// ReadAll 读取所有数据库表结构
func (red *DbRead) ReadAll(dbname string) {
	db := red.db()
	stmt, _ := db.Prepare("SELECT TABLE_NAME FROM TABLES WHERE TABLE_SCHEMA = ?")
	defer stmt.Close()
	rows, _ := stmt.Query(dbname)
	defer rows.Close()
	defer rows.Close()
	var tableName string
	for rows.Next() {
		rows.Scan(&tableName)
		red.readTable(tableName)
	}
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
	rows, _ := stmt.Query(read.DbName, name)
	defer rows.Close()
	tb := read.AllTable[name]
	if tb == nil {
		tb = model.NewTable(name)
	}
	var colMap = map[string]*model.Column{}
	for rows.Next() {
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
		rows.Scan(&tableName, &colName, &isNull, &dataType, &charLen, &numLen, &numScale, &colComment, &colType, &extra, &defval)
		var col model.Column
		col.PhysicalName = colName
		col.LogicalName = colComment
		col.Type = dataType
		col.DefaultValue = defval
		if charLen != "" && charLen != "null" {
			l, _ := strconv.Atoi(charLen)
			col.Length = int(l)
		} else if numLen != "" && numLen != "null" {
			l, _ := strconv.Atoi(numLen)
			col.Length = int(l)
		}
		if numScale != "" && numScale != "null" {
			l, _ := strconv.Atoi(numScale)
			col.Decimal = int(l)
		}
		if extra == "auto_increment" {
			col.AutoIncrement = true
		}
		col.NotNull = isNull == "NO"
		if colType == "" {
			col.ColumnType = dataType
		} else {
			col.ColumnType = colType
		}
		tb.Columns = append(tb.Columns, &col)
		colMap[colName] = &col
	}
	read.readIndex(tb, colMap)
	read.AllTable[tb.PhysicalName] = tb
}

func (red *DbRead) readIndex(table *model.Table, colMap map[string]*model.Column) {
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
	var index *model.Index
	for rows.Next() {
		rows.Scan(&tableName, &nonUnique, &indexName, &colName)
		if indexName == "PRIMARY" {
			col := colMap[colName]
			col.PrimaryKey = true
			table.PrimaryKeys = append(table.PrimaryKeys, col)
		} else if oldIndexName != indexName {
			index = &model.Index{Name: indexName}
			index.Columns = []*model.Column{}
			index.NonUnique, _ = strconv.ParseBool(nonUnique)
			table.Indices = append(table.Indices, index)
			index.Columns = append(index.Columns, colMap[colName])
		} else {
			index.Columns = append(index.Columns, colMap[colName])
		}
		oldIndexName = indexName
	}
}
