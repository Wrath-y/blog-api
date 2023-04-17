package db

import (
	"blog-api/pkg/def"
	"blog-api/pkg/logging"
	"context"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	glog "gorm.io/gorm/logger"
	"log"
	"time"
)

var Orm *gorm.DB

func Setup() {
	Orm = NewMysqlDB("default")
}

func NewMysqlDB(store string) *gorm.DB {
	dbViper := viper.Sub("mysql." + store)
	if dbViper == nil {
		log.Fatal("mysql配置缺失", store)
	}

	address := dbViper.GetString("address")
	username := dbViper.GetString("username")
	password := dbViper.GetString("password")
	database := dbViper.GetString("database")
	maxIdleConns := dbViper.GetInt("max_idle_conns")
	maxOpenConns := dbViper.GetInt("max_open_conns")
	timeout := dbViper.GetString("timeout")
	if timeout == "" {
		timeout = "20"
	}

	dsn := username + ":" + password + "@tcp(" + address + ")/" + database +
		"?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=true&loc=Local&timeout=" + timeout + "s"

	logLevel := glog.Silent
	if viper.GetString("app.env") != def.EnvProduction {
		logLevel = glog.Info
	}

	orm, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: glog.Default.LogMode(logLevel),
	})
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, _ := orm.DB()
	sqlDB.SetMaxIdleConns(maxIdleConns)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	return orm
}

type gormLog struct {
	glog.Interface
}

const msg = "gorm"

func (*gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, rows := fc()
	if err != nil && err != gorm.ErrRecordNotFound {
		logging.FromContext(ctx).ErrorL(msg, sql, err)
	} else if viper.GetString("app.env") == def.EnvDevelopment {
		logging.FromContext(ctx).Info(msg, sql, rows, begin)
	}
}
