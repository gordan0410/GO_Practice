package app

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *Server) Routers() *gin.Engine {
	router := s.router
	router.NoRoute(func(c *gin.Context) {
		p := c.Request.URL.Path
		ps := strings.Split(p, "/")
		fs := gin.Dir("./web", true)
		_, err := fs.Open(ps[1])
		if err != nil {
			c.File("./web/index.html")
		} else {
			h := http.FileServer(http.Dir("./web"))
			h.ServeHTTP(c.Writer, c.Request)
		}
	})
	v1 := router.Group("/v1/api")
	{
		v1.GET("/todolist", s.Get())
		v1.Use(s.Temp())
		v1.POST("/todolist", s.Create())
		v1.PATCH("/todolist", s.Update())
		v1.DELETE("/todolist", s.Delete())
	}
	{
		v1.POST("/login", s.Login())
		v1.POST("/login/register", s.Register())
	}
	return router
}
