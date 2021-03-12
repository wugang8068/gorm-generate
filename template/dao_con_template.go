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
