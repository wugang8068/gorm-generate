package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"unicode"
)

type fileConfig struct {
	DB string
}

type config struct {
	fileConfig     `json:"fileConfig"`
	TableName      string
	FileName       string
	DB             string
	ConfigFilePath string
	ModelName      string
	Directory      string
	DaoDirectory   string
	RepDirectory   string
}

func (c config) ToString() string {
	b, _ := json.MarshalIndent(c, "", "	")
	return string(b)
}

func (c config) GetTableName() string {
	return c.TableName
}

func (c config) GetDirectory() string {
	return c.Directory
}

func (c config) GetFileName() string {
	if len(c.FileName) > 0 {
		return c.FileName
	}
	mn := c.ModelName
	if len(mn) == 0 {
		mn = c.GetTableName()
	}
	mns := strings.SplitAfter(mn, ".")
	if len(mns) >= 2 {
		mn = mns[1]
	}
	if len(mn) == 1 {
		return strings.ToLower(mn)
	}
	buffer := new(bytes.Buffer)
	for i, r := range mn {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteByte('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

func (c config) GetModelName() string {
	mn := c.ModelName
	if len(mn) == 0 {
		mn = c.GetTableName()
	}
	mns := strings.SplitAfter(mn, ".")
	if len(mns) >= 2 {
		mn = mns[1]
	}
	if len(mn) == 1 {
		return strings.ToUpper(mn)
	}
	//下划线转驼峰
	mn = strings.Replace(mn, "_", " ", -1)
	mn = strings.Title(mn)
	return strings.Replace(mn, " ", "", -1)
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
