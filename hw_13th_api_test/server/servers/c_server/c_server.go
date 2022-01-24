package c_server

import (
	"log"
	"net/http"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func C_server(port string) error {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(auth)
	r.GET("/v2/teacher", v2_teacher)
	r.POST("/v2/student", v2_student)
	err := endless.ListenAndServe(":"+port, r)
	if err != nil {
		return err
	}
	return nil
}

type req_msg struct {
	S string `json:"s"`
}

type res_msg[T result_all] struct {
	Err_no  int    `json:"err_no"`
	Err_txt string `json:"err_txt"`
	Result  T      `json:"result"`
}

type result_all interface {
	teacher | student | map[string]string
}

type teacher struct {
	Class string `json:"class"`
}

type student struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

// id錯誤
var id_fail = res_msg[map[string]string]{
	Err_no:  2,
	Err_txt: "id error",
	Result:  nil,
}

// 驗證失敗資訊
var token_fail = res_msg[map[string]string]{
	Err_no:  1,
	Err_txt: "token error",
	Result:  nil,
}

// 驗證
func auth(c *gin.Context) {
	token := c.Request.Header.Get("token")
	if token == "gordan" {
		c.Next()
	} else {
		c.JSON(http.StatusOK, token_fail)
		c.Abort()
		return
	}
}

func v2_teacher(c *gin.Context) {
	// 成功資訊
	var success = res_msg[teacher]{
		Err_no:  0,
		Err_txt: "ok",
		Result:  teacher{Class: "Larry"},
	}

	// 辨識身份
	id := c.Query("t")
	if id == "1" {
		c.JSON(http.StatusOK, success)
	} else {
		c.JSON(http.StatusOK, id_fail)
	}
}

func v2_student(c *gin.Context) {
	// 成功資訊
	success := res_msg[student]{
		Err_no:  0,
		Err_txt: "ok",
		Result:  student{Id: "1", Name: "Larry"},
	}

	// 辨識身份
	var r req_msg
	err := c.ShouldBindJSON(&r)
	if err != nil {
		c.JSON(http.StatusOK, id_fail)
		log.Println(err)
	}
	if r.S == "1" {
		c.JSON(http.StatusOK, success)
	} else {
		c.JSON(http.StatusOK, id_fail)
	}
}
