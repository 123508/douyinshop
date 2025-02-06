package db

import (
	"fmt"

	"github.com/123508/douyinshop/pkg/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Conf.MySQLConfig.User,
		config.Conf.MySQLConfig.Password,
		config.Conf.MySQLConfig.Host,
		config.Conf.MySQLConfig.Port,
		config.Conf.MySQLConfig.DBName,
	)
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := DB.DB()
	sqlDB.SetMaxOpenConns(config.Conf.MySQLConfig.MaxOpenConns) // 设置最大打开连接数
	sqlDB.SetMaxIdleConns(config.Conf.MySQLConfig.MaxIdleConns) // 设置最大空闲连接数

	if err != nil {
		return nil, err
	}
	return DB, nil
}
