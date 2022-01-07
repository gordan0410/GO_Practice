package server

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"hw_ninth/tools"
)

func login_api(c *gin.Context) {
	//set payload
	payload := c.Request.Body

	// request
	r_msg, err := tools.Request_api("http://localhost:30002/login", "POST", payload, nil)
	if err != nil {
		log.Error().Caller().Str("func", "request_api(...)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "request_api error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		return
	}

	// set return msg
	token_raw := r_msg.Data["jwt"]
	var token string
	if token_raw == nil {
		token = ""
	} else {
		token = token_raw.(string)
	}
	status := r_msg.Status
	msg := r_msg.Msg

	// 設定respond cookie
	if status == "success" {
		c.SetCookie("Authorization", token, 0, "", "", false, false)
	}
	// return

	if err := tools.Msg_send(c, status, msg, nil); err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}

}
