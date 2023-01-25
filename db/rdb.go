package db

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	dsn = "root:root@tcp(127.0.0.1:3306)/information_schema?charset=utf8mb4&parseTime=True&loc=Local"
)

var rdb *gorm.DB

func ConnectDB() error {
	var err error

	rdb, err = gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("connot connect to db : ", err.Error())
		return err
	}
	rdb.LogMode(true)

	rdb.DB().SetMaxOpenConns(10)
	rdb.DB().SetMaxIdleConns(10)

	return nil
}

func DB() *gorm.DB {
	return rdb
}

func CloseDB() error {
	err := rdb.DB().Close()
	if err != nil {
		return err
	}

	return nil
}
