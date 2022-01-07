package middlewares

import (
	"hw_ninth/tools"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Received_token struct {
	Token string `json:"token"`
}

func Validate_jwt(c *gin.Context) {
	// 接收前端訊息並驗證
	var rt Received_token
	if err := c.ShouldBindJSON(&rt); err != nil {
		log.Error().Caller().Str("func", "c.ShouldBind(&rt)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "token資訊有誤", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		c.Abort()
		return
	}

	// 解JWT
	claim, err := tools.AuthRequired(c, rt.Token)
	if err != nil {
		var log_message string
		var message string
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				log_message = "token is malformed"
			} else if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				log_message = "token could not be verified because of signing problems"
			} else if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				log_message = "signature validation failed"
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				message = "token expired"
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				log_message = "token is not yet valid before sometime"
			} else {
				log_message = "can not handle this token"
			}
			if log_message != "" {
				log.Error().Caller().Str("func", "tools.AuthRequired").Str("msg", log_message).Msg("Web")
				if err := tools.Msg_send(c, "error", "token error", nil); err != nil {
					log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
				}
			} else {
				if err := tools.Msg_send(c, "error", message, nil); err != nil {
					log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
				}
			}
		} else {
			log.Error().Caller().Str("func", "tools.AuthRequired").Err(err).Msg("Web")
			if err := tools.Msg_send(c, "error", "token error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
		}
		c.Abort()
		return
	}
	c.Set("code_session_id", claim.S_id)
	c.Next()
}
