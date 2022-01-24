package middlewares

import (
	"fmt"
	"hw_ninth/models"
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
)

func Validate_session(c *gin.Context) {
	// 拿加密session_id
	s_id := c.GetString("code_session_id")
	if s_id == "" {
		log.Error().Caller().Str("func", "c.GetString(\"code_session_id\")").Str("msg", "get token error").Msg("Middleware")
		if err := tools.Msg_send(c, "error", "get token error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		c.Abort()
		return
	}

	// 開Redis
	redis_store, err := models.Redis_conn_web(c)
	defer models.Redis_close(redis_store)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Redis_conn(c)").Err(err).Msg("Redis")
		if err := tools.Msg_send(c, "error", "redis error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		c.Abort()
		return
	}

	// create empty session and assign value from db
	session := sessions.NewSession(redis_store, models.Session_name)
	b, err := models.Redis_load(redis_store, s_id, session)

	// 成功
	if b && err == nil {
		// interface 轉 string
		for key, value := range session.Values {
			str_key := fmt.Sprintf("%v", key)
			str_value := fmt.Sprintf("%v", value)
			c.Set(str_key, str_value)
		}
		c.Next()
		// 錯誤
	} else if err != nil {
		log.Error().Caller().Str("func", "Redis_load").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "token error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		c.Abort()
		// 內部資料為空
	} else if !b {
		if err := tools.Msg_send(c, "error", "token expired", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		c.Abort()
	}
}
