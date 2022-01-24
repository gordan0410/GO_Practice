package models

import (
	"os"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
	"time"
	"encoding/json"
	"github.com/rs/zerolog/log"
	"fmt"
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

// 建立table
type Account struct {
	ID         uint       `gorm:"type:bigint(20) NOT NULL auto_increment;primary_key;"`
	Username   string     `gorm:"type:varchar(30) NOT NULL;"`
	Password   string     `gorm:"type:varchar(64) NOT NULL;"`
	CreatedAt  time.Time  `gorm:"type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP"`
}

func Load_database(addr string)(db_conn *gorm.DB, db_conn_err error) {
	// 讀取config.json資料
	file, err := os.Open(addr)
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "os.open").Msg("DB")
		return
	}
	var msg Db_message
	d_data := json.NewDecoder(file)
	err = d_data.Decode(&msg)
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "json.Decode").Msg("DB")
		return
	}

	// gorm資料庫連接
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/?charset=utf8&parseTime=True&loc=Local", msg.Username, msg.Password, msg.Addr)
	db_conn, db_conn_err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if db_conn_err != nil {
		log.Warn().Caller().Err(err).Str("func", "gorm.Open").Msg("DB")
		return
	}

	// 若無則新建資料庫(schema)
	db_conn.Exec("CREATE DATABASE IF NOT EXISTS " + msg.Database)

	// 重新連接指定schema
	dsn = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local", msg.Username, msg.Password, msg.Addr, msg.Database)
	db_conn, db_conn_err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if db_conn_err != nil {
		log.Warn().Caller().Err(err).Str("func", "gorm.Open").Msg("DB")
		return
	}

	// mysql資料庫本體
	db, err := db_conn.DB()
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "db_conn.DB()").Msg("DB")
		return
	}

	// 確保資料庫安全關閉
	db.SetConnMaxLifetime(time.Duration(msg.Max_lifetime) * time.Second)
	// 閒置連接數（官方建議跟SetMaxOpenConns一致）
	db.SetMaxIdleConns(msg.Max_idleconns)
	// 限制資料庫連接數
	db.SetMaxOpenConns(msg.Max_openconns)

	// 更新資料庫資料
	err = db_conn.AutoMigrate(&Account{})
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "db_conn.AutoMigrate").Msg("DB")
		return
	}

	// 確認table存在
	migrator := db_conn.Migrator()
	has := migrator.HasTable(&Account{})
	if !has {
		log.Warn().Caller().Str("func", "migrator.HasTable(&Account{})").Str("msg", "table is not exist").Msg("DB")
		return
	}
	
	return db_conn, db_conn_err
}

func Leave_database(db_conn *gorm.DB) {
	db, err := db_conn.DB()
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "db_conn.DB()").Msg("DB")
		return
	}
	db.Close()
}