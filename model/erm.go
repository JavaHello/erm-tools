package model

import "encoding/xml"

type Diagram struct {
	XMLName    xml.Name   `xml:"diagram"`
	Dictionary Dictionary `xml:"dictionary"`
	Contents   Contents   `xml:"contents"`
}

type Dictionary struct {
	Words []Word `xml:"word"`
}

type Contents struct {
	Table []ErmTable `xml:"table"`
}

type ErmTable struct {
	Id           string  `xml:"id"`
	PhysicalName string  `xml:"physical_name"`
	LogicalName  string  `xml:"logical_name"`
	Description  string  `xml:"description"`
	Columns      Columns `xml:"columns"`
	Indexes      Indexes `xml:"indexes"`
}
type Columns struct {
	NormalColumn []NormalColumn `xml:"normal_column"`
}

type Indexes struct {
	Inidex []ErmInidex `xml:"inidex"`
}
type ErmInidex struct {
	FullText    string       `xml:"full_text"`
	NonUnique   string       `xml:"non_unique"`
	Name        string       `xml:"name"`
	Type        string       `xml:"type"`
	Description string       `xml:"description"`
	Columns     IndexColumns `xml:"columns"`
}

type IndexColumns struct {
	Column []IndexColumn `xml:"column"`
}

type IndexColumn struct {
	Id   string `xml:"id"`
	Desc string `xml:"desc"`
}

type NormalColumn struct {
	WordId        string `xml:"word_id"`
	Id            string `xml:"id"`
	Description   string `xml:"description"`
	UniqueKeyName string `xml:"unique_key_name"`
	LogicalName   string `xml:"logical_name"`
	PhysicalName  string `xml:"physical_name"`
	Type          string `xml:"type"`
	Constraint    string `xml:"constraint"`
	DefaultValue  string `xml:"default_value"`
	AutoIncrement string `xml:"auto_increment"`
	ForeignKey    string `xml:"foreign_key"`
	NotNull       string `xml:"not_null"`
	PrimaryKey    string `xml:"primary_key"`
	UniqueKey     string `xml:"unique_key"`
	CharacterSet  string `xml:"character_set"`
	Collation     string `xml:"collation"`
}
type Word struct {
	Id           string `xml:"id"`
	Length       string `xml:"length"`
	Decimal      string `xml:"decimal"`
	Description  string `xml:"description"`
	PhysicalName string `xml:"physical_name"`
	LogicalName  string `xml:"logical_name"`
	Type         string `xml:"type"`
}
