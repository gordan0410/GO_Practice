package login

import (
	"hw_ninth/models"
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func logout(c *gin.Context) {
	// 拿加密session_id
	s_id := c.GetString("code_session_id")
	if s_id == "" {
		log.Error().Caller().Str("func", "c.GetString(\"code_session_id\")").Str("msg", "get token error").Msg("Middleware")
		if err := tools.Msg_send(c, "error", "get token error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
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
		return
	}

	// 刪除session
	err = models.Redis_delete(redis_store, s_id)
	if err != nil {
		log.Error().Caller().Str("func", "Redis_delete").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "session delete failed", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	if err := tools.Msg_send(c, "success", "logout", nil); err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}
}
