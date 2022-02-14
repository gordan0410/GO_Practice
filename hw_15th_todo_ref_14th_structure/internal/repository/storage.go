package repository

import (
	"errors"
	"todolist/internal/api"

	"github.com/jinzhu/gorm"
)

type Storage interface {
	RunMigration() error
	GetUser(username, password string) error
	GetSubject(userID, slice_target, page int, group string) ([]api.SubjectRespond, error)
	CreateSubject(userID int, newSubject string) error
	UpdateStatus(subjectID int) error
	UpdateSubject(subjectID int, subject string) error
	DeleteSubject(userID int) error
}

type storage struct {
	db *gorm.DB
}

func NewStorage(db *gorm.DB) Storage {
	return &storage{
		db: db,
	}
}

func (s *storage) RunMigration() error {
	all_tables := []interface{}{&Account{}, &Todolist{}}
	// 更新資料庫資料
	err := s.db.AutoMigrate(all_tables...).Error
	if err != nil {
		return err
	}

	// 確認table存在
	for _, v := range all_tables {
		has := s.db.HasTable(v)
		if !has {
			errMsg := "table not exist"
			return errors.New(errMsg)
		}

	}
	return nil
}

// func (s *storage) CreateFirstAccount() error {
// 	var account Account
// 	result := s.db.First(&account)
// 	if result.RecordNotFound() {
// 		account.ID = 1
// 		account.Username = "root"
// 		p, err := encodePassword("root")
// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 		account.Password = p
// 		s.db.Create(&account)
// 		return nil
// 	}
// 	return nil
// }

func (s *storage) GetUser(username, password string) error {
	var account Account
	err := s.db.Where("username = ? AND password = ?", username, password).First(&account).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) GetSubject(userID, slice_target, page int, group string) ([]api.SubjectRespond, error) {
	var todolists []Todolist
	err := s.db.Order("id desc").Scopes(dbQueryGroup(group, 100, 0), dbQueryPage(page, slice_target)).Find(&todolists).Error
	if err != nil {
		return []api.SubjectRespond{}, err
	}
	var res []api.SubjectRespond
	for _, v := range todolists {
		var r api.SubjectRespond
		r.ID = int(v.ID)
		r.Status = v.Status
		r.Subject = v.Subject
		res = append(res, r)
	}

	return res, nil
}

func (s *storage) CreateSubject(userID int, newSubject string) error {
	// 找user
	var user Account
	err := s.db.Where("Id = ? ", userID).First(&user).Error
	if err != nil {
		return err
	}
	// 建立資料並寫入
	err = s.db.Model(&user).Association("Todolists").Append([]Todolist{{Subject: newSubject, Status: 1}}).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) UpdateStatus(subjectID int) error {
	// 先找出欲異動資料
	var todolist Todolist
	err := s.db.Where("ID = ? AND Status <> ?", uint(subjectID), 0).Take(&todolist).Error
	if err != nil {
		return err
	}
	switch todolist.Status {
	case 1:
		err = s.db.Model(&todolist).Update("Status", 2).Error
		if err != nil {
			return err
		}
	case 2:
		err = s.db.Model(&todolist).Update("Status", 1).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *storage) UpdateSubject(subjectID int, subject string) error {
	var todolist Todolist
	err := s.db.Where("ID = ? AND Status <> ?", uint(subjectID), 0).Take(&todolist).Update("Subject", subject).Error
	if err != nil {
		return err
	}
	return nil
}

func (s *storage) DeleteSubject(userID int) error {
	var todolist []Todolist
	count := 0
	// 若資料量>100, 則分批刪除
	for {
		err := s.db.Where("Status = ? AND account_id = ? ", 2, userID).Limit(100).Find(&todolist).Error
		if err != nil {
			return err
		}
		// 判斷是否有資料，無資料則回傳"no object selected"，以提示無資料經選取
		if len(todolist) == 0 && count == 0 {
			errMsg := "no object selected"
			return errors.New(errMsg)

			// 無剩餘可刪除，刪除結束
		} else if len(todolist) == 0 && count > 0 {
			return nil

			// 有資料則更改status為0（軟刪除）並count++
		} else {
			for _, v := range todolist {
				result := s.db.Model(v).Update("Status", "0")
				if result.Error != nil {
					errMsg := "db deleted error"
					return errors.New(errMsg)
				}
			}
			count++
		}
	}
}

// DB搜尋方法，參數為前端回傳group資料＋DB資料限制讀取數＋資料忽略數（之後可再把limit and offset給前端做調整）
func dbQueryGroup(group string, limit_record, offset_record int) func(db *gorm.DB) *gorm.DB {
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
func dbQueryPage(page, slice_target int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * slice_target
		return db.Offset(offset).Limit(slice_target)
	}
}
