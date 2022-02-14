package api

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodolistService interface {
	Get(userID, slice_target, page int, group string) (map[string]interface{}, error)
	Create(n *NewSubjectRequest) error
	Update(u *UpdateSubjectRequest) error
	Delete(d *DeleteSubjectRequest) error
}

type todolistService struct {
	storage TodolistRepository
}

type TodolistRepository interface {
	GetSubject(userID, slice_target, page int, group string) ([]SubjectRespond, error)
	CreateSubject(userID int, newSubject string) error
	UpdateStatus(subjectID int) error
	UpdateSubject(subjectID int, subject string) error
	DeleteSubject(userID int) error
}

func NewTodolistService(tr TodolistRepository) TodolistService {
	return &todolistService{
		storage: tr,
	}
}

func (t *todolistService) Get(userID, slice_target, page int, group string) (map[string]interface{}, error) {
	// 無<1的頁面
	if page <= 0 {
		errMsg := fmt.Sprintf("no other page")
		return nil, errors.New(errMsg)
	}
	todolists, err := t.storage.GetSubject(userID, slice_target, page, group)
	if err != nil {
		return nil, err
	}

	if len(todolists) > 0 {
		data := make(map[string]interface{})
		for i, v := range todolists {
			si := fmt.Sprint(i)
			data[si] = gin.H{"id": v.ID,
				"status":  v.Status,
				"subject": v.Subject}
		}
		return data, nil
		// 第一次啟動
	} else if len(todolists) == 0 && page == 1 {
		errMsg := "no data send"
		return nil, errors.New(errMsg)
		// 無分頁
	} else {
		errMsg := "no other page"
		return nil, errors.New(errMsg)
	}
}

func (t *todolistService) Create(n *NewSubjectRequest) error {
	userID, err := strconv.Atoi(n.User_id)
	if err != nil {
		return err
	}
	err = t.storage.CreateSubject(userID, n.Subject)
	if err != nil {
		return err
	}
	return nil
}

func (t *todolistService) Update(u *UpdateSubjectRequest) error {
	subjectID, err := strconv.Atoi(u.Id)
	if err != nil {
		return err
	}
	if u.Status == "subject_change" {
		err := t.storage.UpdateSubject(subjectID, u.Subject)
		if err != nil {
			return err
		}
		return nil
	} else {
		err := t.storage.UpdateStatus(subjectID)
		if err != nil {
			return err
		}
		return nil
	}
}

func (t *todolistService) Delete(d *DeleteSubjectRequest) error {
	userID, err := strconv.Atoi(d.User_id)
	if err != nil {
		return err
	}
	err = t.storage.DeleteSubject(userID)
	if err != nil {
		return err
	}
	return nil
}
