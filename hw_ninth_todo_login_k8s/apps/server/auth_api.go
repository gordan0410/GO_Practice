package server

import (
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func auth_api(c *gin.Context) {
	// get token
	token, err := c.Cookie("Authorization")
	if err != nil {
		if err := tools.Msg_send(c, "error", "no auth cookies found", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		return
	}
	data := map[string]string{"token": token}
	// request
	r_msg, err := tools.Request_api("http://localhost:30002/auth/", "POST", nil, data)
	if err != nil {
		log.Error().Caller().Str("func", "request_api(...)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "request_api error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		return
	}

	// set return msg
	status := r_msg.Status
	msg := r_msg.Msg

	// return
	if err := tools.Msg_send(c, status, msg, nil); err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}
}
