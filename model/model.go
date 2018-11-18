package model

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type Database struct {
	Self   *gorm.DB
}

var DB *Database

func openDB(username, password, url, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=%t&loc=%s",
		username, password, url, name, true, "Asia/Shanghai")
	db, err := gorm.Open("mysql", config)
	if err != nil {
		fmt.Sprint("数据库连接失败")
	}

	setupDB(db)

	return db
}

func setupDB(db *gorm.DB) {
	db.LogMode(viper.GetBool("mysqllog"))
	db.DB().SetMaxOpenConns(10)
	db.DB().SetMaxIdleConns(0)
}

func linkSelfDB() *gorm.DB {
	return openDB(viper.GetString("database.username"),
		viper.GetString("database.password"),
		viper.GetString("database.url"),
		viper.GetString("database.name"))
}

func GetSelfDB() *gorm.DB {
	return linkSelfDB()
}

func (db *Database) Init() {
	DB = &Database {
		Self: GetSelfDB(),
	}
}

func (db *Database) Close() {
	DB.Self.Close()
}
