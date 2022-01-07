package models

import (
	"encoding/base32"
	"errors"
	"hw_ninth/tools"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	new_redis "github.com/gomodule/redigo/redis"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gopkg.in/boj/redistore.v1"
)

// redis store prefix
var Session_prefix string

// session_name
var Session_name string

// session_maxage
var Session_maxage int

// init redis
func Redis_init(configs *tools.Config_data) (redis_store *redistore.RediStore, err error) {
	// session info load
	Session_prefix = configs.Session.SessionPrefix
	Session_name = configs.Session.SessionName
	Session_maxage = configs.Session.SessionMaxage

	session_key := []byte(configs.Session.SessionKey)

	// 測試連線是否正常，pool內dial不會返回錯誤
	c, err := new_redis.Dial(configs.Redis.Network, configs.Redis.Address)
	if err != nil {
		return nil, err
	}
	c.Close()

	redis_store, err = redistore.NewRediStore(configs.Redis.Size, configs.Redis.Network, configs.Redis.Address, configs.Redis.Password, session_key)
	if err != nil {
		return nil, err
	}

	return redis_store, nil
}

// save session前檢查
func Redis_save(redis_store *redistore.RediStore, session *sessions.Session) (code_session_id string, err error) {
	// 是否過期
	if session.Options.MaxAge < 0 {
		conn := redis_store.Pool.Get()
		defer conn.Close()
		if _, err := conn.Do("DEL", Session_prefix+session.ID); err != nil {
			return "", err
		}
		err = errors.New("session expired and delete")
		return "", err
	} else {
		// Build an alphanumeric key for the redis store.
		if session.ID == "" {
			session.ID = strings.TrimRight(base32.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32)), "=")
		}

		// save
		if err := Redis_upadte(redis_store, session); err != nil {
			return "", err
		}

		// 返回加密後sessionID
		encoded, err := securecookie.EncodeMulti(session.Name(), session.ID, redis_store.Codecs...)
		if err != nil {
			return "", err
		}
		return encoded, nil
	}
}

// 真正的save or update
func Redis_upadte(redis_store *redistore.RediStore, session *sessions.Session) error {
	// 連接redis
	conn := redis_store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}
	// 設定資料過期時間
	age := Session_maxage
	if age == 0 {
		age = redis_store.DefaultMaxAge
	}

	// json序列化, map扁平
	var js redistore.JSONSerializer
	b, err := js.Serialize(session)
	if err != nil {
		return err
	}
	// 存進redis
	_, err = conn.Do("SETEX", Session_prefix+session.ID, age, b)
	return err
}

// 讀session
func Redis_load(redis_store *redistore.RediStore, code_session_id string, session *sessions.Session) (bool, error) {
	// connect to redis
	conn := redis_store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return false, err
	}

	// decode code_session_id
	var session_id string
	err := securecookie.DecodeMulti(Session_name, code_session_id, &session_id, redis_store.Codecs...)
	if err != nil {
		return false, err
	}
	session.ID = session_id

	// 找資料
	data, err := conn.Do("GET", Session_prefix+session_id)
	if data == nil {
		return false, nil // no data was associated with this key
	}
	// byte化
	b, err := redis.Bytes(data, err)
	if err != nil {
		return false, err
	}

	// json 反序列化，並寫入session
	var js redistore.JSONSerializer
	err = js.Deserialize(b, session)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 刪除session
func Redis_delete(redis_store *redistore.RediStore, code_session_id string) error {
	conn := redis_store.Pool.Get()
	defer conn.Close()
	if err := conn.Err(); err != nil {
		return err
	}

	// decode code_session_id
	var session_id string
	err := securecookie.DecodeMulti(Session_name, code_session_id, &session_id, redis_store.Codecs...)
	if err != nil {
		return err
	}

	// 找資料
	data, err := conn.Do("DEL", Session_prefix+session_id)
	if err != nil {
		return err
	}
	if data.(int64) != 1 {
		err = errors.New("affected amount error")
		return err
	}

	return nil
}

func Redis_conn_web(c *gin.Context) (*redistore.RediStore, error) {
	configs_raw, b := c.Get("configs")
	if !b {
		err := errors.New("can't get configs")
		return nil, err
	}
	configs := configs_raw.(*tools.Config_data)
	redis_store, err := Redis_init(configs)
	if err != nil {
		return nil, err
	}
	return redis_store, nil
}

// 關閉session
func Redis_close(redis_store *redistore.RediStore) (s_store_err error) {
	err := redis_store.Close()
	if err != nil {
		return err
	}
	return nil
}
