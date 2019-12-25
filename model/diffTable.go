package model

import "strings"

type DiffTable struct {
	Name        string
	Comment     string
	IsNew       bool
	DiffColumns []DiffColumn
	DiffIndexes []DiffIndex
	DiffPks     []DiffColumn
}

type DiffTableSlice []*DiffTable

func (p DiffTableSlice) Len() int           { return len(p) }
func (p DiffTableSlice) Less(i, j int) bool { return strings.Compare(p[i].Name, p[j].Name) < 0 }
func (p DiffTableSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type DiffColumn struct {
	Name      string
	OldColumn *Column
	NewColumn *Column
}

type DiffIndex struct {
	Name     string
	OldIndex *Index
	NewIndex *Index
}

func NewDiffTable(physicalName, comment string) *DiffTable {
	return &DiffTable{Name: physicalName, Comment: comment, DiffColumns: []DiffColumn{}, DiffIndexes: []DiffIndex{}, DiffPks: []DiffColumn{}}
}
