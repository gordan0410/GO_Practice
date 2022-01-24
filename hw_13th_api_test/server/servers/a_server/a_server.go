package a_server

import (
	"log"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func A_server(port string) error {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(auth)
	r.GET("/api/1", api_1)
	r.POST("/api/2", api_2)
	err := endless.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}
	return nil
}

type req_msg struct {
	Id string `json:"id"`
}

type res_msg[T data_all] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

type data_all interface {
	data1 | data2 | map[string]string
}

type data1 struct {
	Name string `json:"name"`
}

type data2 struct {
	Book string `json:"book"`
	Num  int    `json:"num"`
}

// 驗證失敗資訊
var token_fail = res_msg[map[string]string]{
	Code: 0,
	Msg:  "token error",
	Data: nil,
}

// id錯誤
var id_fail = res_msg[map[string]string]{
	Code: 0,
	Msg:  "id error",
	Data: nil,
}

// 驗證
func auth(c *gin.Context) {
	token := c.Request.Header.Get("key")
	if token == "123" {
		c.Next()
	} else {
		c.JSON(http.StatusOK, token_fail)
		c.Abort()
		return
	}
}

func api_1(c *gin.Context) {
	// 成功資訊
	success := res_msg[data1]{
		Code: 1,
		Msg:  "ok",
		Data: data1{Name: "Larry"},
	}

	// 辨識身份
	id := c.Query("id")
	if id == "1" {
		c.JSON(http.StatusOK, success)
	} else {
		c.JSON(http.StatusOK, id_fail)
	}
}

func api_2(c *gin.Context) {
	// 成功資訊
	success := res_msg[data2]{
		Code: 1,
		Msg:  "ok",
		Data: data2{Book: "Larry", Num: 1},
	}

	// 辨識身份
	var r req_msg
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, id_fail)
		log.Println(err)
	}
	if r.Id == "1" {
		c.JSON(http.StatusOK, success)
	} else {
		c.JSON(http.StatusOK, id_fail)
	}
}
