package login

import (
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func auth(c *gin.Context) {

	user_id := c.GetString("user_id")
	if user_id == "" {
		log.Error().Caller().Str("func", "c.GetString(\"user_id\")").Str("msg", "get user id failed").Msg("Web")
		if err := tools.Msg_send(c, "error", "get user id failed", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}
	data := map[string]interface{}{"user_id": user_id}
	if err := tools.Msg_send(c, "success", "login", data); err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}
}
