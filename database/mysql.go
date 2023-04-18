package database

import (
	"fmt"
	"github.com/liuyp5181/base/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var mysqlList = make(map[string]*gorm.DB)

func NewMysql(cfg config.Database) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", cfg.User, cfg.Pass, cfg.Host, cfg.Port, cfg.DB)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	mysqlList[cfg.Name] = db
	return nil
}

func GetMysql(name string) *gorm.DB {
	// todo ping
	return mysqlList[name]
}
