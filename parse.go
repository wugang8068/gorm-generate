package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

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
	ModelDirectory      string
	RepositoryDirectory string
	DaoDirectory        string
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
	if len(mp.ModelDirectory) > 0 {
		createDirectoryIfNotExist(mp.ModelDirectory)
		return ioutil.WriteFile(mp.ModelDirectory+"/"+mp.FileName+".go", bf.Bytes(), 0755)
	}
	return ioutil.WriteFile(mp.FileName+".go", bf.Bytes(), 0755)
}

func writeDaoFile(mp *modelParse) error {
	if len(mp.DaoDirectory) > 0 {
		daoPath := strings.TrimRight(mp.DaoDirectory, "/") + "/dao"
		createDirectoryIfNotExist(daoPath)
		bf := new(bytes.Buffer)
		bf.WriteString("package dao\n\n")
		modelAbsPath := mp.modelDirectoryAbsPath()
		if len(modelAbsPath) > 0 {
			bf.WriteString(fmt.Sprintf("import models \"%s\"\n\n", modelAbsPath))
		}
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

func createDirectoryIfNotExist(p string) {
	_, err := os.Stat(p)
	if err != nil || os.IsNotExist(err) {
		ps := strings.Split(p, "/")
		if len(ps) > 0 {
			var pt string
			for _, name := range ps {
				pt += name
				_ = os.Mkdir(pt, 0755)
				pt += "/"
			}
		} else {
			_ = os.Mkdir(p, 0755)
		}
	}
}
