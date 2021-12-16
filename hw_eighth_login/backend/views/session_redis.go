package views

import (
	"encoding/base32"
	"errors"
	"os"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gopkg.in/boj/redistore.v1"
)

// redis store prefix
var seesion_prefix = "session_"

// session_name
var session_name = "my-session"

// session key
var session_key = []byte(os.Getenv("COOKIE_KEY"))

// save session前檢查
func Redis_save(store *redistore.RediStore, session *sessions.Session) (code_session_id string, err error) {
	// 是否過期
	if session.Options.MaxAge < 0 {
		conn := store.Pool.Get()
		defer conn.Close()
		if _, err := conn.Do("DEL", seesion_prefix+session.ID); err != nil {
			return "", err
		}
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}
		
		// save
		if err := Redis_upadte(store, session); err != nil {
			return "", err
		}
		// 返回加密後sessionID
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, store.Codecs...)
		if err != nil {
			return "", err
		}
		return encoded, err
	}
	return "", nil
}

// 真正的save or update
func Redis_upadte(store *redistore.RediStore ,session *sessions.Session) error {
	// 連接redis
	conn := store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	// 設定資料過期時間
	age := session.Options.MaxAge
	if age == 0 {
		age = store.DefaultMaxAge
	}

	// gob序列化, map扁平
	var js redistore.JSONSerializer
	b, err := js.Serialize(session)
	if err != nil{
		return err
	}
	// 存進redis
	_, err = conn.Do("SETEX", seesion_prefix+session.ID, age, b)
	return err
}

// 讀session
func Redis_load(store *redistore.RediStore ,code_session_id string, session *sessions.Session) (bool, error) {
	// connect to redis
	conn := store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, err
	}
	
	// decode code_session_id
	var session_id string
	err := securecookie.DecodeMulti(session_name, code_session_id, &session_id, store.Codecs...)
	if err != nil {
		return false, err
	}
	session.ID = session_id

	// 找資料
	data, err := conn.Do("GET", seesion_prefix+session_id)
	if err != nil {
		return false, err
	}

	if data == nil {
		return false, nil // no data was associated with this key
	}

	// byte化
	b, err := redis.Bytes(data, err)
	if err != nil {
		return false, err
	}

	// gob 反序列化，並寫入session
	var js redistore.JSONSerializer
	return true, js.Deserialize(b, session)
}

// 刪除session
func Redis_delete(store *redistore.RediStore ,code_session_id string) error{
	conn := store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}

	// decode code_session_id
	var session_id string
	err := securecookie.DecodeMulti(session_name, code_session_id, &session_id, store.Codecs...)
	if err != nil {
		return err
	}
	
	// 找資料
	data, err := conn.Do("DEL", seesion_prefix+session_id)
	if err != nil {
		return err
	}
	if data.(int64) != 1 {
		err = errors.New("affected amount error")
		return err
	}

	return nil
}