package todolist

import (
	"hw_ninth/models"
	"hw_ninth/tools"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Binding from JSON
type Update_subject struct {
	Id      string `json:"id" binding:"required,min=1"`
	Subject string `json:"subject" binding:"required,min=1"`
	Status  string `json:"status" binding:"required,min=1"`
}

func update(c *gin.Context) {
	// 接收前端訊息
	var subject Update_subject
	if err := c.ShouldBindJSON(&subject); err != nil {
		log.Error().Caller().Str("func", "ShouldBindJSON(&subject)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "subject type error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// id string to int
	subject_id, err := strconv.Atoi(subject.Id)
	if err != nil {
		log.Error().Caller().Str("func", "strconv.Atoi(subject.Id)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "id type error", nil); err != nil {
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

	var todolist models.Todolist
	// 先找出欲異動資料
	err = db_conn.Where("ID = ? AND Status <> ?", uint(subject_id), 0).Take(&todolist).Error
	if err != nil {
		log.Error().Caller().Str("func", "db_conn.Where(\"ID = ? AND Status <> ?\", uint(id_int), 0)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "db query error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// update
	// complete則更新指定ID的status為2（完成的）
	if subject.Status == "complete" {
		// 更新
		err = db_conn.Model(&todolist).Update("Status", "2").Error
		if err != nil {
			log.Error().Caller().Str("func", "db_conn.Model(&todolist).Update(\"Status\", \"2\")").Err(err).Msg("Web")
			if err := tools.Msg_send(c, "error", "db updated error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			return
		}

		if err := tools.Msg_send(c, "success", "updated", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		// active則更新指定ID的status為1（未完成的）
	} else if subject.Status == "active" {
		// 更新
		err = db_conn.Model(&todolist).Update("Status", "1").Error
		if err != nil {
			log.Error().Caller().Str("func", "db_conn.Model(&todolist).Update(\"Status\", \"1\")").Err(err).Msg("Web")
			if err := tools.Msg_send(c, "error", "db updated error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			return
		}

		if err := tools.Msg_send(c, "success", "updated", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		// 判斷是否為更改標題
	} else if subject.Status == "subject_change" {
		err = db_conn.Model(&todolist).Update("Subject", subject.Subject).Error
		if err != nil {
			log.Error().Caller().Str("func", "db_conn.Model(&todolist).Update(\"Subject\", subject.Subject)").Err(err).Msg("Web")
			if err := tools.Msg_send(c, "error", "db updated error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			return
		}

		if err := tools.Msg_send(c, "success", "updated", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}

		// 非以上則錯誤
	} else {
		log.Error().Caller().Str("func", "none").Str("msg", "status unknown").Msg("Web")
		if err := tools.Msg_send(c, "error", "status error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
	}
}
