package models

import (
	"errors"
	"fmt"
	"hw_ninth/tools"

	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 資料庫連接資料
type Db_message struct {
	Username      string
	Password      string
	Addr          string
	Database      string
	Max_lifetime  int
	Max_openconns int
	Max_idleconns int
}

func Db_init(config *tools.Config_data, fk_constraint bool) (db_conn *gorm.DB, err error) {
	// gorm資料庫連接
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8mb4&parseTime=True&loc=Local", config.Mysql.Username, config.Mysql.Password, config.Mysql.Addr)
	db_conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		return nil, err
	}

	// 若無則新建資料庫(schema)
	err = db_conn.Exec("CREATE DATABASE IF NOT EXISTS " + config.Mysql.Database).Error
	if err != nil {
		return nil, err
	}

	// 重新連接指定schema
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.Mysql.Username, config.Mysql.Password, config.Mysql.Addr, config.Mysql.Database)
	db_conn, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DisableForeignKeyConstraintWhenMigrating: !fk_constraint,
	})
	if err != nil {
		return nil, err
	}

	// mysql資料庫本體
	db, err := db_conn.DB()
	if err != nil {
		return nil, err
	}

	// 確保資料庫安全關閉
	db.SetConnMaxLifetime(time.Duration(config.Mysql.MaxLifetime) * time.Second)
	// 閒置連接數（官方建議跟SetMaxOpenConns一致）
	db.SetMaxIdleConns(config.Mysql.MaxIdleconns)
	// 限制資料庫連接數
	db.SetMaxOpenConns(config.Mysql.MaxOpenconns)

	return db_conn, nil
}

func Db_migrate(db_conn *gorm.DB) error {
	all_tables := []interface{}{&Account{}, &Todolist{}}
	// 更新資料庫資料
	err := db_conn.AutoMigrate(all_tables...)
	if err != nil {
		return err
	}

	// 確認table存在
	migrator := db_conn.Migrator()
	for _, v := range all_tables {
		has := migrator.HasTable(v)
		if !has {
			return err
		}

	}

	// 關閉DB
	err = Db_close(db_conn)
	if err != nil {
		return err
	}
	return nil
}

func Db_conn_web(c *gin.Context) (*gorm.DB, error) {
	configs_raw, b := c.Get("configs")
	if !b {
		err := errors.New("can't get configs")
		return nil, err
	}
	configs := configs_raw.(*tools.Config_data)
	db_conn, err := Db_init(configs, false)
	if err != nil {
		return nil, err
	}
	return db_conn, nil
}

func Db_close(db_conn *gorm.DB) error {
	db, err := db_conn.DB()
	if err != nil {
		return err
	}
	err = db.Close()
	if err != nil {
		return err
	}
	return nil
}
