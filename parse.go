package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type tableDcs []tableDesc

func (c tableDcs) ToString() string {
	b, _ := json.MarshalIndent(c, "", "	")
	return string(b)
}

func (c tableDcs) parseFields() (fields []parseField) {
	for _, desc := range c {
		fields = append(fields, parseField{
			Attr: desc.fieldAttrName(),
			Type: desc.fieldType(),
			Tag:  desc.columnTag(),
		})
	}
	return
}

type tableDesc struct {
	Field   string `gorm:"column:Field"`
	Type    string `gorm:"column:Type"`
	Null    string `gorm:"column:Null"`
	Key     string `gorm:"column:Key"`
	Default string `gorm:"column:Default"`
	Extra   string `gorm:"column:Extra"`
}

func (t tableDesc) fieldAttrName() string {
	name := strings.Replace(t.Field, "_", " ", -1)
	name = strings.Title(name)
	return strings.Replace(name, " ", "", -1)
}

func (t tableDesc) columnTag() string {
	return "`gorm:\"column:" + t.Field + "\"`"
}

func (t tableDesc) defaultValue() string {
	return t.Default
}

func (t tableDesc) nullable() bool {
	return t.Null == "YES"
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
	for _, field := range mp.Fields {
		bf.WriteString(fmt.Sprintf("	%s %s %s\n", field.Attr, field.Type, field.Tag))
	}
	bf.WriteString("}\n\n\n")
	if len(mp.TableName) > 0 {
		bf.WriteString(fmt.Sprintf(`func(%s) TableName() string {
	return "%s"
}`, mp.ModelName, mp.TableName))
		bf.WriteString("\n")
	}
	return ioutil.WriteFile(mp.FileName, bf.Bytes(), 0755)
}
