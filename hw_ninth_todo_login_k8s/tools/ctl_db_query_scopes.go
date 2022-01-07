package tools

import (
	"gorm.io/gorm"
)

// DB搜尋方法，參數為前端回傳group資料＋DB資料限制讀取數＋資料忽略數（之後可再把limit and offset給前端做調整）
func Db_query_group(group string, limit_record, offset_record int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		switch {
		case group == "all":
			return db.Where("Status <> ?", 0).Offset(offset_record).Limit(limit_record)
		case group == "active":
			return db.Where("Status = ?", 1).Offset(offset_record).Limit(limit_record)
		case group == "complete":
			return db.Where("Status = ?", 2).Offset(offset_record).Limit(limit_record)
		default:
			return db.Where("Status <> ?", 0).Offset(offset_record).Limit(limit_record)
		}
	}
}

// 設定db頁面呈現方式，參數為前端回傳當前page資料＋欲顯示分頁數（之後可再把分頁數給前端做調整）
func Db_query_page(page, slice_target int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * slice_target
		return db.Offset(offset).Limit(slice_target)
	}
}
