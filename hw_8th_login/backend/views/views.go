package views

import (
	"encoding/hex"
	"errors"
	"fmt"
	"hw_eighth_login/backend/models"
	"net/http"
	"time"

	"crypto/sha256"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog/log"
	"gopkg.in/boj/redistore.v1"
	"gorm.io/gorm"
)

// Binding from JSON
type Login struct {
	Username string `json:"username" binding:"required,alphanum,max=10,min=4"`
	Password string `json:"password" binding:"required,alphanum,max=10,min=4"`
}

// db連線
var db_conn *gorm.DB

// redis session store
var s_store *redistore.RediStore

// validate expire
var expire_sec = 30

func Set_router(db *gorm.DB) {
	// mysql
	db_conn = db

	// redis session store
	var err error
	s_store, err = redistore.NewRediStore(10, "tcp", ":6379", "", session_key)
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "redistore.NewRediStore").Msg("Redis")
		return
	}
	defer s_store.Close()

	// gin
	r := gin.Default()
	r.NoRoute(gin.WrapH(http.FileServer(http.Dir("./templates"))))
	r.POST("/register", register) 	// 註冊
	r.POST("/login", login)       	// 登入
	r_login := r.Group("/api")
	r_login.Use(validate)			// 驗證
	r_login.GET("/", api)           // 首頁
	r_login.GET("/logout", log_out) // 登出

	err = r.Run(":3000")
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "r.Run:3000").Msg("Web")
		return
	}
}

// sha256加密
func encode_password(password string) string {
	h := sha256.New224()
	h.Write([]byte(password))
	result := hex.EncodeToString(h.Sum(nil))
	return result
}

func register(c *gin.Context) {
	// 接收前端訊息並驗證
	var req Login
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Caller().Err(err).Str("func", "c.ShouldBindJSON(&json)").Msg("Web")
		Msg_send(c, "error", "帳號密碼格式錯誤", nil)
		return
	}

	// 加密
	password := encode_password(req.Password)

	// 驗證資料是否存在並儲存
	var account models.Account
	q_result := db_conn.Debug().Where("Username = ?", req.Username).Limit(1).Find(&account)
	if q_result.RowsAffected == 0 {
		// 建立資料並寫入
		create := models.Account{Username: req.Username, Password: password, CreatedAt: time.Now().Local()}
		result := db_conn.Debug().Create(&create)
		if result.Error != nil {
			log.Warn().Caller().Err(result.Error).Str("func", "db_conn.Debug().Create(&create)").Msg("Web")
			Msg_send(c, "error", "db create failed", nil)
			return
		}
		if result.RowsAffected != 1 {
			log.Warn().Caller().Str("func", "db_conn.Debug().Create(&create) and result.RowsAffected != 1").Str("msg", "wrong amount of data been effected").Msg("Web")
			Msg_send(c, "error", "wrong amount of data been effected", nil)
			return
		}
		Msg_send(c, "success", "account created", nil)
	} else if q_result.Error != nil{
		log.Warn().Caller().Err(q_result.Error).Str("func", "db_conn.Debug().Where").Msg("Web")
		Msg_send(c, "error", "db query failed", nil)
	} else if q_result.RowsAffected != 0{
		Msg_send(c, "error", "user has already created", nil)
	}
}

func login(c *gin.Context) {
	// 接收前端訊息並驗證
	var req Login
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Warn().Caller().Err(err).Str("func", "c.ShouldBindJSON(&json)").Msg("Web")
		Msg_send(c, "error", "帳號密碼格式錯誤", nil)
		return
	}

	//加密
	password := encode_password(req.Password)

	// find user in data
	var account models.Account
	err := db_conn.Debug().Where("Username = ? AND Password = ?", req.Username, password).Take(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			Msg_send(c, "error", "user not exist or wrong password", nil)
			return
		} else {
			log.Warn().Caller().Err(err).Str("func", "db_conn.Debug().Where").Msg("Web")
			Msg_send(c, "error", "db query failed", nil)
			return
		}
	}

	// new Session
	session := sessions.NewSession(s_store, session_name)
	session.Values["user_id"] = account.ID
	session.Values["user_name"] = account.Username
	session.Options.MaxAge = expire_sec
	// 儲存, 加密, 返回
	code_session_id, err := Redis_save(s_store, session)
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "Redis_save").Msg("Web")
		Msg_send(c, "error", "session create error", nil)
		return
	}

	// new JWT
	//準備聲明內容
	now := time.Now()
	claims := Claims{
		User_id: account.ID,
		S_id:    code_session_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: now.Add( time.Duration(expire_sec) * time.Second).Unix(),
		},
	}

	// 生成token
	token, err := Jwt_token_get(claims)
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "Jwt_token_get").Msg("Web")
		Msg_send(c, "error", "JWT create error", nil)
		return
	}

	// 設定respond cookie
	c.SetCookie("Authorization", token, 0, "", "", false, false)
	Msg_send(c, "success", "token given", nil)
}

func validate(c *gin.Context) {
	// 解JWT
	claim, err := AuthRequired(c)
	if err != nil {
		Msg_send(c, "error", err.Error(), nil)
		c.Abort()
		return
	}

	// create empty session and assign value from db
	session := sessions.NewSession(s_store, session_name)
	b, err := Redis_load(s_store, claim.S_id, session)

	// 成功
	if b && err == nil {
		// interface 轉 string
		for key, value := range session.Values{
			str_key := fmt.Sprintf("%v", key)
			str_value := fmt.Sprintf("%v", value)
			c.Set(str_key, str_value)
		}
		c.Set("code_session_id", claim.S_id)
		c.Next()
		// 錯誤
	} else if err != nil {
		log.Warn().Caller().Err(err).Str("func", "Redis_load").Msg("Web")
		Msg_send(c, "error", "redis load session failed", nil)
		c.Abort()
		// 內部資料為空
	} else if !b {
		log.Warn().Caller().Str("func", "Redis_load").Msg("Web")
		Msg_send(c, "error", "session data is empty", nil)
		c.Abort()
	}
}

func api(c *gin.Context) {
	data := map[string]interface{}{"user_name": c.GetString("user_name"), "user_id": c.GetString("user_id")}
	Msg_send(c, "success", "login", data)
}

func log_out(c *gin.Context) {
	s_id := c.GetString("code_session_id")
	err := Redis_delete(s_store, s_id)
	if err != nil {
		log.Warn().Caller().Err(err).Str("func", "Redis_delete").Msg("Web")
		Msg_send(c, "error", "session delete failed", nil)
		return
	}
	Msg_send(c, "success", "logout", nil)
}
