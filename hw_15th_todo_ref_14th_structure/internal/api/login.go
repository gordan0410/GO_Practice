package api

import (
	"crypto/sha256"
	"encoding/hex"
)

type LoginService interface {
	Login(lr *LoginRequest) error
}

type loginService struct {
	storage LoginRepository
	session LoginSession
}

type LoginRepository interface {
	GetUser(username, password string) error
}

type LoginSession interface {
}

func NewLoginService(lr LoginRepository, ss LoginSession) LoginService {
	return &loginService{
		storage: lr,
		session: ss,
	}
}

// func (ls *loginService) CreateFirstAccount() error {
// 	err := ls.CreateFirstAccount()
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (ls *loginService) Login(lr *LoginRequest) error {
	password, err := encodePassword(lr.Password)
	if err != nil {
		return err
	}
	err = ls.storage.GetUser(lr.Username, password)
	if err != nil {
		return err
	}

}

func encodePassword(password string) (string, error) {
	h := sha256.New224()
	_, err := h.Write([]byte(password))
	if err != nil {
		return "", err
	}
	result := hex.EncodeToString(h.Sum(nil))
	return result, nil
}
