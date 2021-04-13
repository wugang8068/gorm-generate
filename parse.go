package main

import (
	"os"
	"path"
)

type parseField struct {
	Attr       string
	Type       fieldType
	Tag        string
	ColumnName string
	IsPrimary  bool
}

type modelParse struct {
	ModelPackageName    string
	DaoPackageName      string
	RepoPackageName     string
	FileName            string
	ModelName           string
	Fields              []parseField
	TableName           string
	ModelDirectory      string
	RepositoryDirectory string
	DaoDirectory        string
}

const mysqlDirectory = "mysql"

func (m modelParse) mysqlDirectoryAbsPath() string {
	dir, e := os.Getwd()
	var rootDirectory string
	if e == nil {
		_, rootDirectory = path.Split(dir)
	}
	if len(rootDirectory) > 0 {
		return rootDirectory + "/" + mysqlDirectory
	}
	return mysqlDirectory
}

func (m modelParse) daoDirectoryAbsPath() string {
	dir, e := os.Getwd()
	var rootDirectory string
	if e == nil {
		_, rootDirectory = path.Split(dir)
	}
	if len(rootDirectory) > 0 {
		return rootDirectory + "/" + m.DaoDirectory
	}
	return m.DaoDirectory
}

func (m modelParse) modelDirectoryAbsPath() string {
	dir, e := os.Getwd()
	var rootDirectory string
	if e == nil {
		_, rootDirectory = path.Split(dir)
	}
	if len(rootDirectory) > 0 {
		return rootDirectory + "/" + m.ModelDirectory
	}
	return m.ModelDirectory
}

func (m modelParse) primaryKey() string {
	for _, value := range m.Fields {
		if value.IsPrimary {
			return value.ColumnName
		}
	}
	return "id"
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
