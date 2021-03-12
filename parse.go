package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
			Attr:      desc.fieldAttrName(),
			Type:      desc.fieldType(),
			Tag:       desc.columnTag(),
			IsPrimary: desc.isPrimaryKey(),
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
	return "`json:\"" + t.Field + "\" gorm:\"column:" + t.Field + "\"`"
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
	Attr      string
	Type      fieldType
	Tag       string
	IsPrimary bool
}

type modelParse struct {
	PackageName         string
	FileName            string
	ModelName           string
	Fields              []parseField
	TableName           string
	Directory           string
	RepositoryDirectory string
	DaoDirectory        string
}

func (m modelParse) primaryKeyType() fieldType {
	for _, value := range m.Fields {
		if value.IsPrimary {
			return value.Type
		}
	}
	return TypeUnknown
}

func (m modelParse) RepositoryInterfaceName() string {
	return m.ModelName + "Repository"
}

func (m modelParse) DaoStructName() string {
	return m.ModelName + "Dao"
}

func writeFile(mp *modelParse) error {
	bf := new(bytes.Buffer)
	bf.WriteString("package " + mp.PackageName + "\n\n")
	bf.WriteString(fmt.Sprintf("type %s struct { \n", mp.ModelName))
	for _, field := range mp.Fields {
		bf.WriteString(fmt.Sprintf("	%s %s %s\n", field.Attr, field.Type, field.Tag))
	}
	bf.WriteString("}\n\n")
	if len(mp.TableName) > 0 {
		bf.WriteString(fmt.Sprintf(`func(%s) TableName() string {
	return "%s"
}`, mp.ModelName, mp.TableName))
		bf.WriteString("\n")
	}
	if len(mp.Directory) > 0 {
		createDirectoryIfNotExist(mp.Directory)
		return ioutil.WriteFile(mp.Directory+"/"+mp.FileName+".go", bf.Bytes(), 0755)
	}
	return ioutil.WriteFile(mp.FileName+".go", bf.Bytes(), 0755)
}

func writeDaoFile(mp *modelParse) error {
	if len(mp.DaoDirectory) > 0 {
		daoPath := strings.TrimRight(mp.DaoDirectory, "/") + "/dao"
		createDirectoryIfNotExist(daoPath)
		bf := new(bytes.Buffer)
		bf.WriteString("package dao\n\n")
		bf.WriteString(fmt.Sprintf("type %s struct { }\n\n", mp.DaoStructName()))
		functions := []string{
			fmt.Sprintf("func(%s) List() []*models.%s {\n\n}\n\n", mp.DaoStructName(), mp.ModelName),
			fmt.Sprintf("func(%s) GetById(id %s) (*models.%s, error) {\n\n}\n\n", mp.DaoStructName(), mp.primaryKeyType(), mp.ModelName),
		}
		for _, fs := range functions {
			bf.WriteString(fs)
		}
		return ioutil.WriteFile(daoPath+"/"+mp.FileName+"_dao.go", bf.Bytes(), 0755)
	}
	return nil
}

func writeRepoFile(mp *modelParse) error {
	if len(mp.RepositoryDirectory) > 0 {
		createDirectoryIfNotExist(strings.TrimRight(mp.RepositoryDirectory, "/") + "/repo")
		//bf := new(bytes.Buffer)
	}
	return nil
}

func createDirectoryIfNotExist(path string) {
	_, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		_ = os.Mkdir(path, 0755)
	}
}
