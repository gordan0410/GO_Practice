package todolist

import (
	"hw_ninth/models"
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Binding from JSON
type New_subject struct {
	Subject string `json:"subject" binding:"required,min=1"`
	User_id string `json:"user_id"`
}

func create(c *gin.Context) {
	var subject New_subject
	// 接收前端訊息
	if err := c.ShouldBindJSON(&subject); err != nil {
		log.Error().Caller().Str("func", "ShouldBindJSON(&subject)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "subject type error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	user_id, err := tools.Parse_int(subject.User_id)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Parse_int(user_id_raw").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "user_id convert error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 開DB
	db_conn, err := models.Db_conn_web(c)
	defer models.Db_close(db_conn)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Db_conn(c)").Err(err).Msg("Db")
		if err := tools.Msg_send(c, "error", "db error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 找user
	var user models.Account
	err = db_conn.Where("Id = ? ", user_id).First(&user).Error
	if err != nil {
		log.Error().Caller().Str("func", "db_conn.Where(\"id = ?\", user_id))").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "db query error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}
	// 建立資料並寫入
	err = db_conn.Model(&user).Association("Todolists").Append([]models.Todolist{{Subject: subject.Subject, Status: 1}})
	if err != nil {
		log.Error().Caller().Str("func", "db_conn.Model(&user).Association(\"Todolist\")").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "db create error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 成功
	if err := tools.Msg_send(c, "success", "subject created", nil); err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}
}
