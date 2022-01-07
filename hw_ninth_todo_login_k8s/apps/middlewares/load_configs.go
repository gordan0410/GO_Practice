package middlewares

import (
	"hw_ninth/tools"

	"github.com/gin-gonic/gin"
)

func Load_configs(config *tools.Config_data) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("configs", config)
	}
}
