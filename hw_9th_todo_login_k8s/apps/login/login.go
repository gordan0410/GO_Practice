package login

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
	"hw_ninth/models"
	"hw_ninth/tools"
)

func login(c *gin.Context) {
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

	// 找帳號
	var account models.Account
	err = db_conn.Where("Username = ?", req.Username).Take(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := tools.Msg_send(c, "error", "user not exist", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
		} else {
			log.Error().Caller().Str("func", "db_conn.Where").Err(err).Msg("Web")
			if err := tools.Msg_send(c, "error", "db query failed", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
			}
		}
		return
	}
	if account.Password != password {
		if err := tools.Msg_send(c, "error", "user password incorrect", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// 開Redis
	redis_store, err := models.Redis_conn_web(c)
	defer models.Redis_close(redis_store)
	if err != nil {
		log.Error().Caller().Str("func", "tools.Redis_conn(c)").Err(err).Msg("Redis")
		if err := tools.Msg_send(c, "error", "redis error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// new Session
	session := sessions.NewSession(redis_store, models.Session_name)
	session.Values["user_id"] = account.ID
	session.Values["user_name"] = account.Username
	// 儲存, 加密, 返回
	code_session_id, err := models.Redis_save(redis_store, session)
	if err != nil {
		log.Error().Caller().Str("func", "Redis_save").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "session create error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	// new JWT
	//準備聲明內容
	claims := tools.Claims{
		S_id: code_session_id,
	}

	// 生成token
	token, err := tools.Jwt_token_init(c, claims)
	if err != nil {
		log.Error().Caller().Str("func", "Jwt_token_init").Err(err).Msg("Web")
		if err := tools.Msg_send(c, "error", "JWT create error", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}
		return
	}

	data := map[string]interface{}{"jwt": token}
	// 設定respond cookie
	// c.SetCookie("Authorization", token, 0, "", "", false, false)
	if err := tools.Msg_send(c, "success", "token given", data); err != nil {
		log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		return
	}
}
