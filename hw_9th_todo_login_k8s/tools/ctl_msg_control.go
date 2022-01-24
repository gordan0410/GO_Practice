package tools

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 前端訊息回覆
func Msg_send(c *gin.Context, status string, msg string, data map[string]interface{}) error {
	if status == "error" || status == "success" {
		c.JSON(http.StatusOK, gin.H{"status": status, "msg": msg, "data": data})
		return nil
	} else {
		err := errors.New("status is wrong")
		return err
	}
}
