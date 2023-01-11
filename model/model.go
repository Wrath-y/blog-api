package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

type Database struct {
	Self *gorm.DB
}

type Base struct {
	Id        int    `json:"id"`
	UpdatedAt string `json:"updated_at"`
	CreatedAt string `json:"created_at"`
}

var DB *Database

func openDB(username, password, url, port, name string) *gorm.DB {
	config := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		username, password, url, port, name)
	db, err := gorm.Open("mysql", config)
	if err != nil {
		panic(err)
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
		viper.GetString("database.port"),
		viper.GetString("database.name"))
}

func GetSelfDB() *gorm.DB {
	return linkSelfDB()
}

func (db *Database) Init() {
	DB = &Database{
		Self: GetSelfDB(),
	}
}

func (db *Database) Close() {
	err := DB.Self.Close()
	if err != nil {
		panic(err)
	}
}
