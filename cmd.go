package main

import (
	"errors"
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

var cf config

func init() {
	flag.StringVar(&cf.Directory, "d", "", "Generated directory")
	flag.StringVar(&cf.ModelName, "name", "", "Model name")
	flag.StringVar(&cf.DB, "db", "", "DB connect dns")
	flag.StringVar(&cf.TableName, "t", "", "Table name of generated model")
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
	fmt.Printf("\n%+v\n", cf)
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
