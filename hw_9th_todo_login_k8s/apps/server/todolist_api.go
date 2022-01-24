package server

import (
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func todolist_api(c *gin.Context) {
	// get method
	method := c.Request.Method
	// get param
	param := c.Request.URL.RawQuery
	// get payload
	payload := c.Request.Body

	// get token
	token, err := c.Cookie("Authorization")
	if err != nil {
		if err := tools.Msg_send(c, "error", "no auth cookies found", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		return
	}

	// 先驗證
	req_data := map[string]string{"token": token}
	r_msg, err := tools.Request_api("http://localhost:30002/auth/", "POST", nil, req_data)
	if err != nil {
		log.Error().Caller().Str("func", "request_api(...)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "request_api error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		return
	}

	// get return account data
	user_id_raw := r_msg.Data["user_id"]
	user_id, b := user_id_raw.(string)
	if !b {
		log.Error().Caller().Str("func", "user_data_raw.(string)").Str("msg", "interface convert to string error").Msg("Web")
		if err := tools.Msg_send(c, "error", "interface convert to string error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
	}

	// 驗證完成取資料
	if method == "GET" {
		param = param + "&user_id=" + user_id
		r_msg, err = tools.Request_api("http://localhost:30003/api?"+param, method, nil, nil)
	} else {
		user_data := map[string]string{"user_id": user_id}
		r_msg, err = tools.Request_api("http://localhost:30003/api", method, payload, user_data)
	}
	if err != nil {
		log.Error().Caller().Str("func", "request_api(...)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "request_api error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
		return
	}

	status := r_msg.Status
	msg := r_msg.Msg
	data := r_msg.Data

	// return
	if err := tools.Msg_send(c, status, msg, data); err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}
}
