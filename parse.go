package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type tableDcs []tableDesc

func (c tableDcs) ToString() string {
	b, _ := json.MarshalIndent(c, "", "	")
	return string(b)
}

type tableDesc struct {
	Field   string `gorm:"column:Field"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Key     string `gorm:"column:Key"`
	Default string `gorm:"column:Default"`
	Extra   string `gorm:"column:Extra"`
}

func (t tableDesc) defaultValue() string {
	return t.Default
}

func (t tableDesc) nullable() bool {
	return t.Null != "NO"
}

func (t tableDesc) isPrimaryKey() bool {
	return t.Key == "PRI"
}

func (t tableDesc) fieldType() fieldType {
	if strings.HasPrefix(t.Type, "int") || strings.HasPrefix(t.Type, "bigint") {
		if strings.HasSuffix(t.Type, "unsigned") {
			return TypeUInt32
		}
		return TypeInt32
	}
	if strings.HasPrefix(t.Type, "tinyint") || strings.HasPrefix(t.Type, "smallint") || strings.HasPrefix(t.Type, "mediumint") {
		if strings.HasSuffix(t.Type, "unsigned") {
			return TypeUInt8
		}
		return TypeInt8
	}
	if strings.HasPrefix(t.Type, "varchar") || strings.HasPrefix(t.Type, "text") {
		return TypeString
	}
	if strings.HasPrefix(t.Type, "float") || strings.HasPrefix(t.Type, "double") || strings.HasPrefix(t.Type, "decimal") {
		return TypeFloat
	}
	return TypeUnknown
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
	//for _, field := range mp.Fields {
	//
	//}
	return nil
}
