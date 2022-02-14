package app

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"todolist/internal/api"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// todolist

func (s *Server) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := getQueryToInt(c, "user_id")
		if err != nil {
			log.Error().Caller().Str("func", "getQueryToInt(c, \"user_id\"").Err(err).Msg("Server")
			if err := msgSend(c, "error", "user_id convert error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}

		sliceTarget, err := getQueryToInt(c, "slice_target")
		if err != nil {
			log.Error().Caller().Str("func", "getQueryToInt(c, \"slice_target\"").Err(err).Msg("Server")
			if err := msgSend(c, "error", "slice_target convert error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}

		page, err := getQueryToInt(c, "page")
		if err != nil {
			log.Error().Caller().Str("func", "getQueryToInt(c, \"page\"").Err(err).Msg("Server")
			if err := msgSend(c, "error", "page convert error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}

		group, ok := c.GetQuery("group")
		if !ok {
			log.Error().Caller().Str("func", "c.GetQuery(\"group\")").Str("msg", "get group failed").Msg("Server")
			if err := msgSend(c, "error", "get group failed", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}
		data, err := s.todolistService.Get(userID, sliceTarget, page, group)
		if err != nil {
			switch err.Error() {
			case "no other page":
				if err := msgSend(c, "error", "no other page", nil); err != nil {
					log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
				}
				return
			case "no data send":
				if err := msgSend(c, "success", "no data send", nil); err != nil {
					log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
				}
				return
			default:
				if err := msgSend(c, "error", "no other page", nil); err != nil {
					log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
				}
				return
			}
		}

		if err := msgSend(c, "success", "list send", data); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			return
		}
	}
}

func (s *Server) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newSubject api.NewSubjectRequest
		err := c.ShouldBindJSON(&newSubject)
		if err != nil {
			log.Error().Caller().Str("func", "ShouldBindJSON(&subject)").Err(err).Msg("Server")
			if err := msgSend(c, "error", "subject type error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}
		err = s.todolistService.Create(&newSubject)
		if err != nil {
			log.Error().Caller().Str("func", "s.todolistService.Create").Err(err).Msg("Server")
			if err := msgSend(c, "error", "create failed", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}
		// 成功
		if err := msgSend(c, "success", "subject created", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			return
		}

	}
}

func (s *Server) Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var updateSubject api.UpdateSubjectRequest
		err := c.ShouldBindJSON(&updateSubject)
		if err != nil {
			log.Error().Caller().Str("func", "ShouldBindJSON(&updateSubject)").Err(err).Msg("Server")
			if err := msgSend(c, "error", "subject type error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}
		err = s.todolistService.Update(&updateSubject)
		if err != nil {
			log.Error().Caller().Str("func", "s.todolistService.Update").Err(err).Msg("Server")
			if err := msgSend(c, "error", "update failed", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}
		// 成功
		if err := msgSend(c, "success", "updated", nil); err != nil {
			log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Web")
		}

	}
}

func (s *Server) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var deleteSubject api.DeleteSubjectRequest
		err := c.ShouldBindJSON(&deleteSubject)
		if err != nil {
			log.Error().Caller().Str("func", "ShouldBindJSON(&deleteSubject)").Err(err).Msg("Server")
			if err := msgSend(c, "error", "subject type error", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}
		err = s.todolistService.Delete(&deleteSubject)
	}
}

// 前端訊息回覆
func msgSend(c *gin.Context, status string, msg string, data map[string]interface{}) error {
	if status == "error" || status == "success" {
		c.JSON(http.StatusOK, gin.H{"status": status, "msg": msg, "data": data})
		return nil
	} else {
		err := errors.New("status is wrong")
		return err
	}
}

// get query and covert to int
func getQueryToInt(c *gin.Context, key string) (int, error) {
	rv, ok := c.GetQuery(key)
	if !ok {
		errMsg := fmt.Sprintf("get %s failed.", key)
		return 0, errors.New(errMsg)
	}
	v, err := strconv.Atoi(rv)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// login

func (s *Server) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest api.LoginRequest
		err := c.ShouldBindJSON(&loginRequest)
		if err != nil {
			log.Error().Caller().Str("func", "ShouldBindJSON(&loginRequest)").Err(err).Msg("Server")
			if err := msgSend(c, "error", "帳號密碼格式錯誤", nil); err != nil {
				log.Error().Caller().Str("func", "tools.Msg_send").Err(err).Msg("Server")
			}
			return
		}
		err = 
	}
}
