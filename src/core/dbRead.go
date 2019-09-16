package core

import (
	"database/sql"
	"fmt"

	"erm-tools/src/model"
	_ "github.com/go-sql-driver/mysql"
)

type DbRead struct {
	AbstractRead
	database *sql.DB
}

func NewDbRead() DbRead {
	return DbRead{AbstractRead: AbstractRead{AllTable: make(map[string]model.Table, 16)}}
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

func (red *DbRead) ReadAll(dbname string) {

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
		COLUMN_TYPE
	from COLUMNS
	where TABLE_SCHEMA = ?
	and TABLE_NAME = 'tm_test'`)
	defer stmt.Close()
	rows, _ := stmt.Query("demodb")
	defer rows.Close()
	var tableName string
	var colName string
	var idNull string
	var dataType string
	var charLen string
	var numLen string
	var numScale string
	var colComment string
	var colType string
	for rows.Next() {
		rows.Scan(&tableName, &colName, &idNull, &dataType, &charLen, &numLen, &numScale, &colComment, &colType)
		fmt.Print("tableName = ", tableName)
		fmt.Println(", colName = ", colName)
	}
	read.readIndex()
}

func (red *DbRead) readIndex() {
	db := red.db()
	stmt, _ := db.Prepare(
		`select TABLE_NAME, NON_UNIQUE, INDEX_NAME, COLUMN_NAME
		from STATISTICS
		where TABLE_SCHEMA = ?
		  and TABLE_NAME = 'tm_test'`)
	defer stmt.Close()
	rows, _ := stmt.Query("demodb")
	defer rows.Close()
	var tableName string
	var nonUnique string
	var indexName string
	var colName string
	for rows.Next() {
		rows.Scan(&tableName, &nonUnique, &indexName, &colName)
		fmt.Print("tableName = ", tableName)
		fmt.Print(", indexName = ", indexName)
		fmt.Println(", colName = ", colName)
	}
}
