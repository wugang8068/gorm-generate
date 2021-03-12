package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var cf config
var con *gorm.DB

func init() {
	flag.StringVar(&cf.FileName, "f", "", "Generate file name")
	flag.StringVar(&cf.Directory, "d", "", "Generated directory")
	flag.StringVar(&cf.ModelName, "name", "", "Model name")
	flag.StringVar(&cf.DB, "db", "", "DB connect dns")
	flag.StringVar(&cf.TableName, "t", "", "Table name of generated model")
	flag.StringVar(&cf.DaoDirectory, "dao", "", "The directory of dao generate.")
	flag.StringVar(&cf.RepDirectory, "repo", "", "The directory of repository generate.")
	flag.StringVar(&cf.ConfigFilePath, "c", "", "Special config file, format: .yml")
}

func main() {
	flag.Parse()
	// load config
	if e := readConfigFromFile(&cf); e != nil {
		fmt.Println(e.Error())
		return
	}
	if e := cf.Validate(); e != nil {
		fmt.Println(e.Error())
		return
	}
	c, e := connect(cf.GetDNS())
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	con = c
	mp, e := getTableDescription()
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	if e := writeFile(mp); e != nil {
		fmt.Println(e.Error())
		return
	}
	if e := writeDaoFile(mp); e != nil {
		fmt.Println(e.Error())
		return
	}
}

func getTableDescription() (*modelParse, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("%+v\n", err)
		}
	}()
	tableName := cf.GetTableName()
	if con.HasTable(tableName) == false {
		return nil, errors.New("table \"" + tableName + "\" not exist")
	}
	var result tableDcs
	con.Raw("DESCRIBE " + tableName).Scan(&result)
	parse := modelParse{
		PackageName:         "models",
		Directory:           cf.GetDirectory(),
		FileName:            cf.GetFileName(),
		ModelName:           cf.GetModelName(),
		Fields:              result.parseFields(),
		TableName:           cf.GetTableName(),
		DaoDirectory:        cf.DaoDirectory,
		RepositoryDirectory: cf.RepDirectory,
	}
	return &parse, nil
}

func readConfigFromFile(cfg *config) error {
	configFile := cfg.ConfigFilePath
	// load default config if input config file is empty
	if len(configFile) == 0 {
		defaultFile, e := os.Stat(".yml")
		if e == nil && defaultFile.IsDir() == false {
			configFile = ".yml"
		}
	}
	if len(configFile) > 0 {
		bts, e := ioutil.ReadFile(configFile)
		if e != nil {
			return errors.New("config file open error:" + e.Error())
		}
		e = yaml.Unmarshal(bts, &cfg.fileConfig)
		if e != nil {
			return errors.New("config file format error:" + e.Error())
		}
	}
	return nil
}
