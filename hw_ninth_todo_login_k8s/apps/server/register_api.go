package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"hw_ninth/tools"
)

func register_api(c *gin.Context) {
	//set payload
	payload := c.Request.Body

	// request
	r_msg, err := tools.Request_api("http://localhost:30002/register", "POST", payload, nil)
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
