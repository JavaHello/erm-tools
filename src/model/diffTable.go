package model

type DiffTable struct {
	Name        string
	DiffColumns []DiffColumn
	DiffIndexes []DiffIndex
	DiffUniques []DiffIndex
	DiffPks     []DiffColumn
}

type DiffColumn struct {
	Name      string
	OldColumn *Column
	NewColumn *Column
}

type DiffIndex struct {
	Name     string
	OldIndex *Column
	NewIndex *Column
}

func NewDiffTable(name string) DiffTable {
	return DiffTable{Name: name, DiffColumns: []DiffColumn{}, DiffIndexes: []DiffIndex{}, DiffUniques: []DiffIndex{}, DiffPks: []DiffColumn{}}
}
