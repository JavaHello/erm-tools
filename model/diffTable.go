package model

type DiffTable struct {
	Name        string
	Comment     string
	IsNew       bool
	DiffColumns []DiffColumn
	DiffIndexes []DiffIndex
	DiffPks     []DiffColumn
}

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

func NewDiffTable(physicalName, comment string) DiffTable {
	return DiffTable{Name: physicalName, Comment: comment, DiffColumns: []DiffColumn{}, DiffIndexes: []DiffIndex{}, DiffPks: []DiffColumn{}}
}
