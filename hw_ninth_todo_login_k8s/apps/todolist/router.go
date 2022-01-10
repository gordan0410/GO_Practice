package todolist

import (
	"hw_ninth/apps/middlewares"
	"hw_ninth/tools"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Set_router(configs *tools.Config_data) {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middlewares.Load_configs(configs))

	// for k8s livenessprobe
	r.GET("/alive", alive)

	// funcitons
	r.GET("/api", get)       // 首頁
	r.POST("/api", create)   // 建立
	r.PATCH("/api", update)  // 更新
	r.DELETE("/api", delete) // 刪除

	err := endless.ListenAndServe(":30003", r)
	if err != nil {
		log.Error().Caller().Str("func", "endless.ListenAndServe(\":30003\", r)").Err(err).Msg("Web")
		return
	}
}
