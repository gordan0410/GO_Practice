package api

import (
	"encoding/base32"
	"encoding/json"
	"errors"
	"strings"

	"github.com/gorilla/securecookie"
)

type SessionService interface {
	NewSession(name string, maxage int, prefix, key string, data map[string]string) *SessionData
	Save(sd *SessionData) (sessionID string, err error)
	load(sessionID string) (*SessionData, error)
}

type sessionService struct {
	Store  SessionRepository
	Name   string
	key    string
	Prefix string
	MaxAge int
}

type SessionRepository interface {
	Get(key string) (interface{}, error)
	SaveOrCreate(key string, maxage int, value []byte) error
	Delete(key string) error
}

func NewSessionService(store SessionRepository, name, key, prefix string, maxage int) SessionService {
	return &sessionService{
		Store:  store,
		Name:   name,
		key:    key,
		Prefix: prefix,
		MaxAge: maxage,
	}
}

func (ss *sessionService) NewSession(name string, maxage int, prefix, key string, data map[string]string) *SessionData {
	return &SessionData{
		Name:   name,
		MaxAge: maxage,
		Prefix: prefix,
		Codecs: securecookie.CodecsFromPairs([]byte(key)),
		ID:     strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "="),
		Datas:  data,
	}
}

func (ss *sessionService) Save(sd *SessionData) (sessionID string, err error) {
	// 設定資料過期時間
	age := sd.MaxAge
	if age < 0 {
		err := ss.Store.Delete(sd.Prefix + sd.ID)
		if err != nil {
			return "", err
		}
		err = errors.New("session expired and delete")
		return "", err
	} else if age == 0 {
		age = 240
	}

	// json序列化, map扁平
	b, err := json.Marshal(sd.Datas)
	if err != nil {
		return "", err
	}

	// 設定key
	key := sd.Prefix + sd.ID

	// 存入redis
	err = ss.Store.SaveOrCreate(key, age, b)
	if err != nil {
		return "", err
	}

	// encode sessionID
	sessionID, err = securecookie.EncodeMulti(sd.Name, sd.ID, sd.Codecs...)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (ss *sessionService) load(sessionID string) (*SessionData, error) {
	// implement session
	sessionData := SessionData{
		Name:   ss.Name,
		Codecs: securecookie.CodecsFromPairs([]byte(ss.key)),
		Prefix: ss.Prefix,
		MaxAge: ss.MaxAge,
	}

	// decode sessionID
	err := securecookie.DecodeMulti(sessionData.Name, sessionID, sessionData.ID, sessionData.Codecs...)
	if err != nil {
		return nil, err
	}

	// find the session
	v, err := ss.Store.Get(sessionData.Prefix + sessionData.ID)
	if err != nil {
		return nil, err
	}
	b, ok := v.([]byte)
	if !ok {
		errMsg := ("redis value error")
		return nil, errors.New(errMsg)
	}
	err = json.Unmarshal(b, sessionData.Datas)
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}
