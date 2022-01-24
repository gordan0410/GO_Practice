package todolist

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"hw_ninth/tools"
)

func alive(c *gin.Context) {
	err := tools.Msg_send(c, "success", "service alive", nil)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}
}
