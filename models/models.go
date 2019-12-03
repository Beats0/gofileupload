package models

import (
	"fmt"
	"github.com/beats0/gofileupload/conf"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var db *gorm.DB

func Setup() {
	var err error
	log.Print("[info] start mysql db")
	dbPath := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.DatabaseSetting.DbUser,
		conf.DatabaseSetting.DbPass,
		conf.DatabaseSetting.DbHost,
		conf.DatabaseSetting.DbPort,
		conf.DatabaseSetting.DbName,
	)
	db, err = gorm.Open(conf.DatabaseSetting.DbType, dbPath)
	if err != nil {
		log.Fatalf("models.Setup err: %v", err)
	}

	db.SingularTable(true)
	CheckTable()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
}

func CheckTable() {
	if !db.HasTable("upload") {
		db.CreateTable(Upload{})
	} else {
		db.AutoMigrate(Upload{})
	}
	if !db.HasTable("user") {
		db.CreateTable(User{})
	} else {
		db.AutoMigrate(User{})
	}
}
