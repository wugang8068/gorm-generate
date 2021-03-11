package main

import (
	"bytes"
	"fmt"
)

type tableDcs []tableDesc

type tableDesc struct {
	Field   string `gorm:"column:Field"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Key     string `gorm:"column:Key"`
	Default string `gorm:"column:Default"`
	Extra   string `gorm:"column:Extra"`
}

func (t tableDesc) fieldType() fieldType {

}

type parseField struct {
	Attr string
	Type fieldType
	Tag  string
}

type modelParse struct {
	PackageName string
	FileName    string
	ModelName   string
	Fields      []parseField
	TableName   string
}

func writeFile(mp modelParse) error {
	bf := new(bytes.Buffer)
	bf.WriteString("package " + mp.PackageName + "\n\n\n")
	bf.WriteString(fmt.Sprintf("type %s struct { \n", mp.ModelName))
	for _, field := range mp.Fields {

	}
}
