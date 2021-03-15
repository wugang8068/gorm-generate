package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var conTemplate = `
package mysql

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var con *gorm.DB

func DefaultConnection() *gorm.DB {
	if con == nil {
		con = connect("{{ dns }}")
	}
	return con
}

func connect(dns string) *gorm.DB {
	var err error
	connection, err := gorm.Open("mysql", dns)
	if err != nil {
		panic(errors.New("db connection error"))
	}
	connection.DB().SetConnMaxLifetime(time.Duration(300) * time.Second)
	connection.DB().SetMaxOpenConns(200)
	connection.DB().SetMaxIdleConns(50)
	return connection.Unscoped()
}
`

func writeConnectionUtilTool(mp *modelParse) error {
	if len(mp.DaoDirectory) > 0 {
		path := mysqlDirectory
		createDirectoryIfNotExist(path)
		conPath := path + "/con.go"
		_, e := os.Stat(conPath)
		if e != nil && os.IsNotExist(e) {
			return ioutil.WriteFile(path+"/con.go", []byte(conTemplate), 0755)
		}
		return nil
	}
	return nil
}

func writeModelFile(mp *modelParse) error {
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
		daoPath := strings.TrimRight(mp.DaoDirectory, "/")
		createDirectoryIfNotExist(daoPath)
		bf := new(bytes.Buffer)
		bf.WriteString("package dao\n\n")
		modelAbsPath := mp.modelDirectoryAbsPath()
		if len(modelAbsPath) > 0 || modelAbsPath != mp.daoDirectoryAbsPath() {
			bf.WriteString("import (\n")
			bf.WriteString("	\"github.com/jinzhu/gorm\"\n")
			bf.WriteString(fmt.Sprintf("	models \"%s\"\n", modelAbsPath))
			bf.WriteString(fmt.Sprintf("	\"%s\"\n", mp.mysqlDirectoryAbsPath()))
			bf.WriteString(")\n\n")
		}
		bf.WriteString(fmt.Sprintf("type %s struct { }\n\n", mp.DaoStructName()))
		functions := []string{
			fmt.Sprintf("func(%s) List() (l []*models.%s) {\n"+
				"	mysql.DefaultConnection().Order(\"%s desc\").Find(&l) \n"+
				"	return\n"+
				"}\n\n", mp.DaoStructName(), mp.ModelName, mp.primaryKey()),
			fmt.Sprintf("func(%s) GetById(id %s) (*models.%s, error) {\n"+
				"	var m models.%s\n"+
				"	e := mysql.DefaultConnection().Where(\"%s = ?\", id).First(&m).Error\n"+
				"	if e != nil {\n"+
				"		if gorm.IsRecordNotFoundError(e) {\n"+
				"			return nil, nil\n"+
				"		}\n"+
				"		return nil, e \n"+
				"	}\n"+
				"	return &m, nil\n"+
				"}\n\n", mp.DaoStructName(), mp.primaryKeyType(), mp.ModelName, mp.ModelName, mp.primaryKey()),
			fmt.Sprintf("func (%s) Create(m models.%s) (*models.%s, error)  {\n"+
				"	e := mysql.DefaultConnection().Create(&m).Error\n"+
				"	if e != nil {\n"+
				"		return nil, e\n"+
				"	}\n"+
				"	return &m, nil\n"+
				"}\n\n", mp.DaoStructName(), mp.ModelName, mp.ModelName),
			fmt.Sprintf("func (%s) Update(m models.%s, updates map[string]interface{}) (*models.%s, error) {\n"+
				"	if len(updates) == 0 {\n"+
				"		return &m, nil\n"+
				"	}\n"+
				"	e := mysql.DefaultConnection().Model(&m).UpdateColumns(updates).Error\n"+
				"	if e != nil {\n"+
				"		return nil, e\n"+
				"	}\n"+
				"	return &m, nil\n"+
				"}\n\n", mp.DaoStructName(), mp.ModelName, mp.ModelName),
			fmt.Sprintf("func (%s) Delete(m models.%s) error {\n"+
				"	return mysql.DefaultConnection().Delete(m).Error\n"+
				"}\n\n", mp.DaoStructName(), mp.ModelName),
		}
		for _, fs := range functions {
			bf.WriteString(fs)
		}
		if e := writeConnectionUtilTool(mp); e != nil {
			return e
		}
		return ioutil.WriteFile(daoPath+"/"+mp.FileName+"_dao.go", bf.Bytes(), 0755)
	}
	return nil
}

func writeRepoFile(mp *modelParse) error {
	if len(mp.RepositoryDirectory) > 0 {
		repoPath := strings.TrimRight(mp.RepositoryDirectory, "/")
		createDirectoryIfNotExist(repoPath)
		bf := new(bytes.Buffer)
		bf.WriteString("package repo\n\n")
		daoAbsPath := mp.daoDirectoryAbsPath()
		modelAbsPath := mp.modelDirectoryAbsPath()
		bf.WriteString("import (\n")
		bf.WriteString(fmt.Sprintf("	models \"%s\"\n", modelAbsPath))
		if len(daoAbsPath) > 0 {
			bf.WriteString(fmt.Sprintf("	dao \"%s\"\n", daoAbsPath))
		}
		bf.WriteString(")\n\n")
		bf.WriteString(fmt.Sprintf("type %sRepository interface {\n", mp.ModelName))
		bf.WriteString(fmt.Sprintf("	List() (l []*models.%s)\n", mp.ModelName))
		bf.WriteString(fmt.Sprintf("	GetById(id %s) (*models.%s, error)\n", mp.primaryKeyType(), mp.ModelName))
		bf.WriteString(fmt.Sprintf("	Create(m models.%s) (*models.%s, error)\n", mp.ModelName, mp.ModelName))
		bf.WriteString(fmt.Sprintf("	Update(m models.%s, updates map[string]interface{}) (*models.%s, error)\n", mp.ModelName, mp.ModelName))
		bf.WriteString(fmt.Sprintf("	Delete(m models.%s) error\n", mp.ModelName))
		bf.WriteString("}\n\n")
		if len(daoAbsPath) > 0 {
			bf.WriteString(fmt.Sprintf("func New%sRepository() %sRepository {\n"+
				"	return dao.%s{}\n"+
				"}\n\n", mp.ModelName, mp.ModelName, mp.DaoStructName()))
		}
		return ioutil.WriteFile(repoPath+"/"+mp.FileName+"_repo.go", bf.Bytes(), 0755)
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
