package main

import (
	"errors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

func connect(dns string) (*gorm.DB, error) {
	var err error
	connection, err := gorm.Open("mysql", dns)
	if err != nil {
		return nil, errors.New("db connection error:" + err.Error())
	}
	connection.LogMode(true)
	connection.DB().SetConnMaxLifetime(time.Duration(300) * time.Second)
	connection.DB().SetMaxOpenConns(200)
	connection.DB().SetMaxIdleConns(50)
	return connection.Unscoped(), nil
}
