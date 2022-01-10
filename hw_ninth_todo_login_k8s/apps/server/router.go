package server

import (
	"net/http"
	"strings"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Set_router() {
	r := gin.New()
	r.Use(gin.Recovery())

	r.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		ps := strings.Split(p, "/")
		fs := gin.Dir("./views", true)
		_, err := fs.Open(ps[1])
		if err != nil {
			c.File("./views/index.html")
		} else {
			h := http.FileServer(http.Dir("./views"))
			h.ServeHTTP(c.Writer, c.Request)
		}
	})

	// for k8s livenessprobe
	r.GET("/alive", alive)

	// 登入
	r.POST("/login", login_api)
	r.POST("/register", register_api)
	r.GET("/logout", logout_api)
	r.GET("/auth", auth_api)

	// todolist
	r.GET("/todolist", todolist_api)
	r.POST("/todolist", todolist_api)
	r.PATCH("/todolist", todolist_api)
	r.DELETE("/todolist", todolist_api)

	err := endless.ListenAndServe(":30004", r)
	if err != nil {
		log.Error().Caller().Str("func", "endless.ListenAndServe(\":30004\", r)").Err(err).Msg("Web")
		return
	}
}
