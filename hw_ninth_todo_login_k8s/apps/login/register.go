package login

import (
	"hw_ninth/models"
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func register(c *gin.Context) {
	// 接收前端訊息並驗證
	var req Login
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().Caller().Str("func", "c.ShouldBindJSON(&json)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "帳號密碼格式錯誤", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 加密
	password, err := tools.Encode_password(req.Password)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Encode_password(req.Password)").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "密碼加密錯誤", nil); err != nil {
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

	// 驗證資料是否存在並儲存
	var account models.Account
	q_result := db_conn.Where("Username = ?", req.Username).Limit(1).Find(&account)
	if q_result.RowsAffected == 0 {
		// 建立資料並寫入
		create := models.Account{Username: req.Username, Password: password}
		result := db_conn.Create(&create)
		if result.Error != nil {
			log.Error().Caller().Str("func", "db_conn.Create(&create)").Err(result.Error).Msg("Db")
			if err := tools.Msg_send(c, "error", "db create failed", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			return
		}
		if result.RowsAffected != 1 {
			log.Error().Caller().Str("func", "db_conn.Create(&create) and result.RowsAffected != 1").Str("msg", "wrong amount of data been effected").Msg("Web")
			if err := tools.Msg_send(c, "error", "wrong amount of data been effected", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
			return
		}
		if err := tools.Msg_send(c, "success", "account created", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
	} else if q_result.Error != nil {
		log.Error().Caller().Str("func", "db_conn.Where(\"Username = ?\", req.Username)").Err(q_result.Error).Msg("Db")
		if err := tools.Msg_send(c, "error", "db query failed", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
	} else if q_result.RowsAffected != 0 {
		if err := tools.Msg_send(c, "error", "user has already created", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			return
		}
	}
}
