package views

import (
	"net/http"
	"github.com/rs/zerolog/log"
	"github.com/gin-gonic/gin"
)


func Msg_send(c *gin.Context, status string, msg string, data map[string]interface{} ){
	if status == "error" || status == "success"{
		c.JSON(http.StatusOK, gin.H{status:msg, "data": data})
	} else {
		log.Warn().Caller().Str("func", "Msg_send").Str("msg", "msg not send").Msg("Web")
	}
}