package main

import "errors"

type fileConfig struct {
	DB string
}

type config struct {
	fileConfig
	TableName      string
	DB             string
	ConfigFilePath string
	ModelName      string
	Directory      string
}

func (c config) GetTableName() string {
	return c.TableName
}

func (c config) GetDNS() string {
	if len(c.DB) > 0 {
		return c.DB
	}
	return c.fileConfig.DB
}

func (c config) Validate() error {
	if len(c.GetDNS()) == 0 {
		return errors.New("db connect name must not be blank")
	}
	if len(cf.GetTableName()) == 0 {
		return errors.New("table name must not be blank")
	}
	return nil
}
