package login

import (
	"hw_ninth/apps/middlewares"
	"hw_ninth/tools"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// Binding from JSON
type Login struct {
	Username string `json:"username" binding:"required,alphanum,max=10,min=4"`
	Password string `json:"password" binding:"required,alphanum,max=10,min=4"`
}

func Set_router(configs *tools.Config_data) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.Load_configs(configs))
	r.Use()
	r.GET("/alive", alive)
	r.POST("/register", register) // 註冊
	r.POST("/login", login)       // 登入
	r_login := r.Group("/auth")
	r_login.Use(middlewares.Validate_jwt)
	r_login.Use(middlewares.Validate_session)
	r_login.POST("/", auth)
	r_login.GET("/logout", logout) // 登出

	err := endless.ListenAndServe(":30002", r)
	// err := r.Run(":30002")
	if err != nil {
		log.Error().Caller().Str("func", "endless.ListenAndServe(\":30002\", r)").Err(err).Msg("Web")
		return
	}
}
