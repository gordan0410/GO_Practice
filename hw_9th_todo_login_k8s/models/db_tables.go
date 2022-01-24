package models

import "gorm.io/gorm"

/// Accounts table
type Account struct {
	gorm.Model
	Username  string     `gorm:"type:varchar(30) NOT NULL;"`
	Password  string     `gorm:"type:varchar(64) NOT NULL;"`
	Todolists []Todolist `gorm:"foreignKey:AccountID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

// Todolist table
type Todolist struct {
	gorm.Model
	Subject   string `gorm:"type:varchar(30) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL;"`
	Status    int    `gorm:"type:int NOT NULL;"`
	AccountID uint
}
